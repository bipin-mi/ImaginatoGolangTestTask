package utils

import (
	_const "ImaginatoGolangTestTask/shared/utils/const"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// RandomKeyGenerator is for generating random string
func RandomKeyGenerator(strSize int, randType string) string {
	var dictionary string

	if randType == _const.AlphaNum {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if randType == _const.Alpha {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == _const.Number {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	_, _ = rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// HashedPassword() is Password generator
func HashedPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("ERROR : ", err.Error())
		return ""
	}
	return string(hashedPassword)
}

func TokenGenerator() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func PageAttributes(pageStr, sizeStr string) (int, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = _const.PageNo
	}
	size, errSize := strconv.Atoi(sizeStr)
	if errSize != nil || size < 1 {
		size = _const.PerPageLimit
	}
	return page, size
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
