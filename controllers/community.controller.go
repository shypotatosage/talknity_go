package controllers

import (
	"net/http"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCommunities(c echo.Context) error {

	result, err := models.FetchAllCommunities()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}