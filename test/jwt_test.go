package test

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"testing"
)

type UserClaims struct {
	Name     string `json:"name"`
	Identity string `json:"identity"`
	jwt.RegisteredClaims
}

var mykey = []byte("oj_key")

func TestGenerateToken(t *testing.T) {
	UserClaims := &UserClaims{
		Name:             "user_1",
		Identity:         "Get",
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims)
	tokenString, err := token.SignedString(mykey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tokenString)
}

func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidXNlcl8xIiwiaWRlbnRpdHkiOiJHZXQifQ.D4dh-J42_fNQgJShupWQuJ0GQO3TtN_zgDzdhSGLJEk"
	UserClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, UserClaims, func(token *jwt.Token) (interface{}, error) {
		return mykey, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if claims.Valid {
		fmt.Println(UserClaims)
	}

}
