package auth

import (
	"errors"

	"net/http"
	"strings"
	"time"

	"invoiceinaja/config"
	"invoiceinaja/domain/user"
	"invoiceinaja/helper"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Service interface {
	GenerateTokenJWT(id int, name string, role string) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateTokenJWT(id int, name string, role string) (string, error) {

	claim := jwt.MapClaims{}
	claim["id"] = id
	claim["fullname"] = name
	claim["role"] = role
	claim["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(config.InitConfiguration().JWT_KEY))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(config.InitConfiguration().JWT_KEY), nil

	})

	if err != nil {
		return token, err
	}

	return token, nil
}

func AuthMiddleware(authService Service, userService user.IService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// mendapatkan bearer token dari header
		authHeader := c.GetHeader("Authorization")

		// cek apakah terdapat bearer token
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.ApiResponse("Unauthorized1", http.StatusUnauthorized, "error", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// mengambil token jwt
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// validasi token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.ApiResponse("Unauthorized2", http.StatusUnauthorized, "error", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.ApiResponse("Unauthorized3", http.StatusUnauthorized, "error", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["id"].(float64))

		user, err := userService.GetUserById(userID)
		if err != nil {
			response := helper.ApiResponse("Unauthorized4", http.StatusUnauthorized, "error", nil, nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
