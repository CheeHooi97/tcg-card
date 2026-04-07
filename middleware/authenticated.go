package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

const currentActorKey = "CurrentActor"

func GetActor(c echo.Context) (*Actor, error) {
	actor, ok := c.Get(currentActorKey).(*Actor)
	if !ok {
		return nil, errors.New("missing actor from requests")
	}
	return actor, nil
}

type Actor struct {
	Id    string `json:"-"`
	Token string `json:"-"`
	Role  string `json:"-"`
}

type tokenIntrospectionFunc func(echo.Context, string) (*Actor, error)
type tokenExtractFunc func(echo.Context) string

func Authenticated(tokenExtract tokenExtractFunc, tokenVerify tokenIntrospectionFunc) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := tokenExtract(c)
			data, err := tokenVerify(c, token)

			if data == nil || err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
			}

			c.Set(currentActorKey, data)
			return next(c)
		}
	}
}

func OptionalAuthenticated(tokenExtract tokenExtractFunc, tokenVerify tokenIntrospectionFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := tokenExtract(c)
			if token == "" {
				return next(c)
			}

			data, err := tokenVerify(c, token)
			if data == nil || err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Token is invalid")
			}

			c.Set(currentActorKey, data)
			return next(c)
		}
	}
}
