package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"shifo-website-backend/internal/entity"
	"shifo-website-backend/internal/pkg/config"
	"shifo-website-backend/internal/pkg/repository/postgres"
	"strings"
	"time"
)

type Auth struct {
	postgresDB *postgres.Database
}

func New(postgresDB *postgres.Database) *Auth {
	return &Auth{postgresDB: postgresDB}
}

func (a Auth) GenerateToken(ctx context.Context, data GenerateToken) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(12 * time.Hour).Unix()
	claims["username"] = data.Username
	claims["role"] = data.Role

	tokenString, err := token.SignedString([]byte(config.GetConf().JWTKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a Auth) IsValidToken(ctx context.Context, tokenStr string) (entity.User, error) {
	claims := new(Claims)
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConf().JWTKey), nil
	})
	if err != nil {
		return entity.User{}, err
	}
	if !tkn.Valid {
		return entity.User{}, errors.New("invalid token")
	}

	query := fmt.Sprintf(`SELECT id, username,role FROM users WHERE username = '%s' AND deleted_at IS NULL`, claims.Username)

	rows, err := a.postgresDB.QueryContext(ctx, query)
	if err != nil {
		return entity.User{}, err
	}

	var detail entity.User

	for rows.Next() {
		if err = rows.Scan(&detail.Id, &detail.Username, &detail.Role); err != nil {
			return entity.User{}, err
		}
	}

	return detail, nil
}

func (a Auth) GetTokenData(ctx context.Context, token string) (TokenData, error) {
	detail, err := a.IsValidToken(ctx, token)
	if err != nil {
		return TokenData{}, err
	}

	if detail.Username == "" {
		return TokenData{}, errors.New("no such user")
	}

	var tokenData TokenData

	tokenData.Username = detail.Username
	tokenData.UserId = detail.Id

	return tokenData, nil
}

func (a Auth) HasPermission(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lang string
		if len(c.Request.Header["Accept-Language"]) > 0 {
			lang = c.Request.Header["Accept-Language"][0]
		} else {
			defaultLang := config.GetConf().DefaultLang
			lang = defaultLang
		}
		tokenStr := c.Request.Header["Authorization"]

		if len(tokenStr) > 0 {

			splitToken := strings.Split(tokenStr[0], " ")
			if len(splitToken) != 2 || splitToken[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Invalid token!",
					"status":  false,
				})
			} else {
				ctx := context.Background()
				userDetail, err := a.IsValidToken(ctx, splitToken[1])
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": err.Error(),
						"status":  false,
					})
					return
				}

				hasPermission := false
				for _, r := range roles {
					if userDetail.Role == r {
						hasPermission = true
						break
					}
				}
				if !hasPermission {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"message": "User has not permission!",
						"status":  false,
					})
				}

				c.Set("role", userDetail.Role)
				c.Set("userId", userDetail.Id)
				c.Set("lang", lang)
				c.Next()
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Send token!",
				"status":  false,
			})
		}
	}
}
