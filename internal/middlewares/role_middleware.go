package middleware

import (
	"net/http"
	"strings"

	"marketly-app/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AllowRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "missing authorization header",
				})
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			claims, err := utils.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"message": "invalid or expired token",
				})
			}

			token := &jwt.Token{
				Claims: claims,
				Valid:  true,
			}

			c.Set("user", token)
			c.Set("userId", claims.UserID)
			c.Set("role", claims.Role)

			if len(roles) == 0 {
				return next(c)
			}

			for _, r := range roles {
				if r == claims.Role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"message": "forbidden: insufficient role",
			})
		}
	}
}
