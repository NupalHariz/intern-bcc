package jwt

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// type IJwt interface {
// 	GenerateToken(userId uuid.UUID) (string, error)
// 	ValidateToken(tokenString string) (uuid.UUID, error)
// 	GetLoginUser(ctx *gin.Context) (domain.Users, error)
// }

// type jsonWebToken struct {
// 	SecretKey   string
// 	ExpiredTime time.Duration
// }

type Claims struct {
	UserId uuid.UUID
	jwt.RegisteredClaims
}

// func JwtInit() IJwt {
// 	secretKey := os.Getenv("SECRET_KEY")
// 	expiredTime, err := strconv.Atoi(os.Getenv("JWT_EXP_TIME"))
// 	if err != nil {
// 		log.Fatal("failed to set jwt expired time")
// 	}

//		return &jsonWebToken{
//			SecretKey:   secretKey,
//			ExpiredTime: time.Duration(expiredTime) * time.Hour,
//		}
//	}
func GenerateToken(userId uuid.UUID) (string, error) {
	expiredTime, err := strconv.Atoi(os.Getenv("JWT_EXP_TIME"))
	secretKey := os.Getenv("SECRET_KEY")
	if err != nil {
		log.Fatal("failed to set jwt expired time")
	}

	claim := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiredTime) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (uuid.UUID, error) {
	var userId uuid.UUID
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return userId, err
	}

	if !token.Valid {
		return userId, err
	}

	userId = claims.UserId

	return userId, nil
}

func GetLoginUserId(ctx *gin.Context) (uuid.UUID, error) {
	userId, ok := ctx.Get("userId")
	if !ok {
		return userId.(uuid.UUID), errors.New("failed to get user")
	}

	return userId.(uuid.UUID), nil
}
