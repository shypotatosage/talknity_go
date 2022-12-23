package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"talknity/models"
	"time"

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

func FetchCommunity(c echo.Context) error {

	result, err := models.FetchCommunity()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreCommunity(c echo.Context) error {
	name := c.FormValue("community_name")
	description := c.FormValue("community_description")
	contact := c.FormValue("community_contact")
	logo, err := c.FormFile("community_logo")

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	src, err := logo.Open()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}
	
	defer src.Close()

	fileByte, _ := io.ReadAll(src)
	fileName := "community/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"

	err = os.WriteFile(fileName, fileByte, 0777)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	cid, err := strconv.ParseUint(c.FormValue("category_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	lid, err := strconv.ParseUint(c.FormValue("leader_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.StoreCommunity(name, description, contact, fileName, cid, lid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}