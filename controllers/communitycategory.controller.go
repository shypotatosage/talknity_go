package controllers

import (
	"io"
	"net/http"
	"os"
	"strings"
	"talknity/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func FetchAllCommunityCategory(c echo.Context) error {

	result, err := models.FetchAllCommunityCategory()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func SearchCommunityCategory(c echo.Context) error {
	key := c.Param("search_key")

	result, err := models.SearchCommunityCategory(key)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StoreCommunityCategory(c echo.Context) error {
	name := c.FormValue("category_name")
	logo, err := c.FormFile("category_logo")
	color1 := c.FormValue("category_color1")
	color2 := c.FormValue("category_color2")
	color3 := c.FormValue("category_color3")

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
	fileName := "images/category/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(logo.Filename, ".")[1]

	err = os.WriteFile(fileName, fileByte, 0777)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.StoreCommunityCategory(name, fileName, color1, color2, color3)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}