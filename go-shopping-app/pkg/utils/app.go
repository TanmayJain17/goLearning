package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

//signature
var jwtKey = []byte("secretkey")

//payload
type JWTClaim struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

/* func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err == nil {
		return nil
	} else {
		return fmt.Errorf("enter a valid email")
	}
} */

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWTToken(email string, firstname string, lastname string) (tokenString string, err error) {
	//setting expiration time
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	if ok {
		return claims, nil
	}
	return
}
