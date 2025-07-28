package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Mattcazz/Peer-Presure.git/utils"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userId int, userName string) (string, error) {
	expiration := time.Second * time.Duration(3600*24*7) // 7 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userId),
		"expiredAt": time.Now().Add(expiration).Unix(),
		"username":  userName,
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTAuth(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: no auth token"))
			return
		}

		token_string := cookie.Value

		token, err := JWTAuthWeb(token_string)

		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		id, err := strconv.Atoi(claims["userID"].(string))

		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: Id needs to be an int"))
			return
		}
		username := claims["username"].(string)

		// we add user_id to context so that we can access it later using r.Context().Value("user_id").(int)
		ctx := context.WithValue(r.Context(), ctxKeyUserID, id)
		// we add username to context so that we can access it later using r.Context().Value("username").(string)
		ctx = context.WithValue(ctx, ctxUsername, username)

		next(w, r.WithContext(ctx))

	}
}

func JWTAuthWeb(token_string string) (*jwt.Token, error) {

	token, err := ValidateJWTtoken(token_string)

	if err != nil {
		return nil, fmt.Errorf("unauthorized: invalid token 1")
	}

	claims := token.Claims.(jwt.MapClaims)

	if !token.Valid || claims["expiredAt"] == nil || claims["userID"] == nil {

		return nil, fmt.Errorf("unauthorized: invalid token 2")
	}

	if time.Now().Unix() > int64(claims["expiredAt"].(float64)) {
		return nil, fmt.Errorf("unauthorized: expired token")
	}

	return token, nil
}

func ValidateJWTtoken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, err
}

type ctxKey string

const (
	ctxKeyUserID ctxKey = "user_id"
	ctxUsername  ctxKey = "username"
)
