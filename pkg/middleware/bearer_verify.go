package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/isogram/clean-golang/pkg/entity"
	"github.com/isogram/clean-golang/pkg/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// BearerVerify function to verify token
func BearerVerify(IsRefreshToken bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			var myClaim entity.AuthTokenClaim

			req := c.Request()
			header := req.Header

			authorizationHeader := header.Get("Authorization")
			if authorizationHeader != "" {
				bearerToken := strings.Split(authorizationHeader, " ")
				if len(bearerToken) == 2 {
					token, err := jwt.ParseWithClaims(bearerToken[1], &myClaim, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("There was an error")
						}
						return []byte(os.Getenv("JWT_SECRET_KEY")), nil
					})

					if err != nil && token == nil {
						return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
					}

					if token.Valid {
						ID := utils.HashedToInt64(myClaim.UID)
						c.Set("JWT", token)
						c.Set("UID", ID)
						return next(c)
					} else if ve, ok := err.(*jwt.ValidationError); ok {
						var errorStr string
						if ve.Errors&jwt.ValidationErrorMalformed != 0 {
							errorStr = fmt.Sprintf("Invalid token format: %s", bearerToken[1])
						} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
							if IsRefreshToken == true {
								ID := utils.HashedToInt64(myClaim.UID)
								c.Set("JWT", token)
								c.Set("UID", ID)
								return next(c)
							}
							errorStr = "Token has been expired"
						} else {
							errorStr = fmt.Sprintf("Token Parsing Error: %s", err.Error())
						}
						return echo.NewHTTPError(http.StatusUnauthorized, errorStr)
					} else {
						return echo.NewHTTPError(http.StatusUnauthorized, "Unknown token error")
					}
				} else {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization token")
				}
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "An authorization header is required")
			}
		}
	}
}
