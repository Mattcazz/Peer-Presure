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

func CreateJWT(secret []byte, userId int) (string, error) {
	expiration := time.Second * time.Duration(3600*24*7) // 7 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    strconv.Itoa(userId),
		"expires_at": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JwtAuth(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth_token")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: no auth token"))
			return
		}

		token_string := cookie.Value

		token, err := ValidateJWTtoken(token_string)

		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: invalid token"))
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		if !token.Valid || claims["expires_at"] == nil || claims["user_id"] == nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: invalid token"))
			return
		}

		if time.Now().Unix() > int64(claims["expires_at"].(float64)) {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: expired token"))
			return
		}

		id, err := strconv.Atoi(claims["user_id"].(string))

		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: Id needs to be an int"))
			return
		}

		// we add user_id to context so that we can access it later using r.Context().Value("user_id").(int)
		ctx := context.WithValue(r.Context(), "user_id", id)

		next(w, r.WithContext(ctx))

	}
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
