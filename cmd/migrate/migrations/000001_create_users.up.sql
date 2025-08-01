CREATE EXTENSION IF NOT EXISTS citext; 

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY, 
    email  citext UNIQUE NOT NULL, 
    username varchar (255) UNIQUE NOT NULL,
    password bytea NOT NULL, 
    profile_pic_url VARCHAR(10000) NOT NULL,
    created_at timestamp (0) with time zone NOT NULL DEFAULT NOW ()
);
