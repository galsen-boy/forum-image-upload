package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Crytper(mot_pass string) []byte {
	passwd := []byte(mot_pass)

	hastedPasswd, err := bcrypt.GenerateFromPassword(passwd, bcrypt.DefaultCost)
	if err != err {
		fmt.Println(err.Error())
	}

	return hastedPasswd
}

func Decrytper(passSaisie string, passRecup []byte) bool {

	err := bcrypt.CompareHashAndPassword(passRecup, []byte(passSaisie))
	if err != nil {
		fmt.Println("Login et ou passewd incorrect")
		return false
	}
	return true
}
