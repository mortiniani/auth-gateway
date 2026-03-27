package auth_gateway

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	SESSION_STORE   = "session-store"
	SALT_LENGTH    = 16
	COOKIE_EXPIRES = 30 * 24 * time.Hour
)

func generateRandomSalt() (salt []byte, err error) {
	salt = make([]byte, SALT_LENGTH)
	_, err = rand.Read(salt)
	if err != nil {
		return
	}
	return
}

func hashPassword(password string) (hashedPassword string, err error) {
	salt, err := generateRandomSalt()
	if err != nil {
		return
	}
	hashedPassword, err = bcrypt.GenerateFromPassword(salt, []byte(password))
	if err != nil {
		return
	}
	return
}

func comparePassword(password, hashedPassword string) (match bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return
	}
	match = true
	return
}

func generateToken(userID int) (token string, err error) {
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(COOKIE_EXPIRES).Unix(),
	}).SignedString([]byte("auth-secret"))
	return
}

func validateToken(tokenString string) (userID int, err error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{
		"userID": jwt.IdRequired(),
		"exp":    jwt.IdRequired(),
	}, func(token *jwt.Token) (interface{}, error {
		return []byte("auth-secret"), nil
	})
	if err != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userID"].(float64)
		if !ok {
			return
		}
		return int(userID), nil
	}
	return
}

func getStore() (store *sessions.CookieStore) {
	return sessions.NewCookieStore(SESSION_STORE)
}

func GetCurrentUser(r *http.Request) (userID int, err error) {
	store := getStore()
	session, err := store.Get(r, SESSION_STORE)
	if err != nil {
		return
	}
	tokenString, ok := session.Values["token"]
	if !ok || tokenString == nil {
		return
	}
	return validateToken(fmt.Sprintf("%s", tokenString))
}

func SetCurrentUser(r *http.Request, userID int) (err error) {
	store := getStore()
	session, err := store.Get(r, SESSION_STORE)
	if err != nil {
		return
	}
	token, err := generateToken(userID)
	if err != nil {
		return
	}
	session.Values["token"] = token
	return store.Save(r, session)
}

func ClearCurrentUser(r *http.Request) (err error) {
	store := getStore()
	session, err := store.Get(r, SESSION_STORE)
	if err != nil {
		return
	}
	session.Values["token"] = nil
	return store.Save(r, session)
}

func GetError(err error) (errMessage string) {
	if err == nil {
		return
	}
	switch t := err.(type) {
	case bcrypt.ErrMismatchedHashAndPassword:
		errMessage = "Invalid username or password"
	case jwt.ValidationError:
		errMessage = "Invalid JWT"
	default:
		errMessage = "Internal error"
	}
	return
}