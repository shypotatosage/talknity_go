package controllers

import (
	"net/http"
	"strconv"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCommunityMember(c echo.Context) error {

	cid, err := strconv.ParseUint(c.FormValue("community_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.FetchAllCommunityMember(cid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}