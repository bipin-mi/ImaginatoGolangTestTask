package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"os"
	"strings"
	"time"
)

type IMiddleware interface {
	AuthHandler() gin.HandlerFunc
}

// Middleware is
type Middleware struct{}

func NewMiddlewareService() IMiddleware {
	return &Middleware{}
}

// AuthHandler is
func (m *Middleware) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")

		b := "Bearer "
		if !strings.Contains(token, b) {
			c.JSON(403, gin.H{"message": "Your request is not authorized", "status": 403})
			c.Abort()
			return
		}
		t := strings.Split(token, b)
		if len(t) < 2 {
			c.JSON(403, gin.H{"message": "An authorization token was not supplied", "status": 403})
			c.Abort()
			return
		}

		// Validate token
		valid, err := ValidateToken(t[1], os.Getenv("AccessSecret"))
		if err != nil {
			c.JSON(403, gin.H{"message": "Invalid authorization token", "status": 403})
			c.Abort()
			return
		}

		// set userId Variable
		c.Set("userData", valid.Claims.(jwt.MapClaims)["userData"])
		c.Next()
	}
}

func GenerateToken(k []byte, userData interface{}) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Hour * 8760).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(k)
	return tokenString, err
}

func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	})

	return token, err
}
