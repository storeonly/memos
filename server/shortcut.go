package server

import (
	"encoding/json"
	"fmt"
	"memos/api"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) registerShortcutRoutes(g *echo.Group) {
	g.POST("/shortcut", func(c echo.Context) error {
		userId := c.Get(getUserIdContextKey()).(int)
		shortcutCreate := &api.ShortcutCreate{
			CreatorId: userId,
		}
		if err := json.NewDecoder(c.Request().Body).Decode(shortcutCreate); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted post shortcut request").SetInternal(err)
		}

		shortcut, err := s.ShortcutService.CreateShortcut(shortcutCreate)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create shortcut").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(shortcut)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode shortcut response").SetInternal(err)
		}

		return nil
	})

	g.PATCH("/shortcut/:shortcutId", func(c echo.Context) error {
		shortcutId, err := strconv.Atoi(c.Param("shortcutId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ID is not a number: %s", c.Param("shortcutId"))).SetInternal(err)
		}

		shortcutPatch := &api.ShortcutPatch{
			Id: shortcutId,
		}
		if err := json.NewDecoder(c.Request().Body).Decode(shortcutPatch); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Malformatted patch shortcut request").SetInternal(err)
		}

		shortcut, err := s.ShortcutService.PatchShortcut(shortcutPatch)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to patch shortcut").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(shortcut)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode shortcut response").SetInternal(err)
		}

		return nil
	})

	g.GET("/shortcut", func(c echo.Context) error {
		userId := c.Get(getUserIdContextKey()).(int)
		shortcutFind := &api.ShortcutFind{
			CreatorId: &userId,
		}
		list, err := s.ShortcutService.FindShortcutList(shortcutFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch shortcut list").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(list)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode shortcut list response").SetInternal(err)
		}

		return nil
	})

	g.GET("/shortcut/:shortcutId", func(c echo.Context) error {
		shortcutId, err := strconv.Atoi(c.Param("shortcutId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ID is not a number: %s", c.Param("shortcutId"))).SetInternal(err)
		}

		shortcutFind := &api.ShortcutFind{
			Id: &shortcutId,
		}
		shortcut, err := s.ShortcutService.FindShortcut(shortcutFind)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch shortcut").SetInternal(err)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		if err := json.NewEncoder(c.Response().Writer).Encode(composeResponse(shortcut)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode shortcut response").SetInternal(err)
		}

		return nil
	})

	g.DELETE("/shortcut/:shortcutId", func(c echo.Context) error {
		shortcutId, err := strconv.Atoi(c.Param("shortcutId"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("ID is not a number: %s", c.Param("shortcutId"))).SetInternal(err)
		}

		shortcutDelete := &api.ShortcutDelete{
			Id: shortcutId,
		}
		if err := s.ShortcutService.DeleteShortcut(shortcutDelete); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete shortcut").SetInternal(err)
		}

		c.JSON(http.StatusOK, true)

		return nil
	})
}
