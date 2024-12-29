package api

import (
	"eatsome/internal/terrors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *API) Health(c echo.Context) error {
	stats, err := a.storage.Health()
	if err != nil {
		return terrors.InternalServerError(err, "cannot get health stats")
	}

	return c.JSON(http.StatusOK, stats)
}
