package server

import (
	"encoding/json"
	"fmt"
	"memos/api"
	"memos/common"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) registerAuthRoutes(g *echo.Group) {
	g.POST("/auth/login", func(c echo.Context) error {
		login := &api.Login{}
		if err := json.NewDecoder(c.Request().Body).Decode(login); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted login request").SetInternal(err)
		}

		userFind := &api.UserFind{
			Name: &login.Name,
		}
		user, err := s.UserService.FindUser(userFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate user").SetInternal(err)
		}
		if user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("User not found: %s", login.Name))
		}

		// Compare the stored hashed password, with the hashed version of the password that was received.
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
			// If the two passwords don't match, return a 401 status.
			return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect password").SetInternal(err)
		}

		err = setUserSession(c, user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set login session").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode user response").SetInternal(err)
		}

		return nil
	})

	g.POST("/auth/logout", func(c echo.Context) error {
		err := removeUserSession(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set logout session").SetInternal(err)
		}

		c.Response().WriteHeader(http.StatusOK)
		return nil
	})

	g.POST("/auth/signup", func(c echo.Context) error {
		signup := &api.Signup{}
		if err := json.NewDecoder(c.Request().Body).Decode(signup); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted signup request").SetInternal(err)
		}

		if len(signup.Name) <= 5 {
			return echo.NewHTTPError(http.StatusBadRequest, "Username is too short, minimum length is 6.")
		}
		if len(signup.Password) <= 5 {
			return echo.NewHTTPError(http.StatusBadRequest, "Password is too short, minimum length is 6.")
		}

		userFind := &api.UserFind{
			Name: &signup.Name,
		}
		user, err := s.UserService.FindUser(userFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to authenticate user").SetInternal(err)
		}
		if user != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Existed user found: %s", signup.Name))
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate password hash").SetInternal(err)
		}

		userCreate := &api.UserCreate{
			Name:         signup.Name,
			PasswordHash: string(passwordHash),
			OpenId:       common.GenUUID(),
		}
		user, err = s.UserService.CreateUser(userCreate)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user").SetInternal(err)
		}

		err = setUserSession(c, user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to set signup session").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(user)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode created user response").SetInternal(err)
		}

		return nil
	})
}
