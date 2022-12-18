package controllers

import (
	"net/http"
	"strconv"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func FetchAllPosts(c echo.Context) error {

	result, err := models.FetchAllPosts()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StorePost(c echo.Context) error {

	title := c.FormValue("post_title")
	content := c.FormValue("post_content")
	image := c.FormValue("post_image")
	anonymous, err := strconv.ParseBool(c.FormValue("anonymous"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message1": err.Error()})
	}

	uid, err := strconv.ParseUint(c.FormValue("uid"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message2": err.Error()})
	}

	result, err := models.StorePost(title, content, image, anonymous, uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message3": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdatePost(c echo.Context) error {

	title := c.FormValue("post_title")
	content := c.FormValue("post_content")
	image := c.FormValue("post_image")
	anonymous, err := strconv.ParseBool(c.FormValue("anonymous"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	pid, err := strconv.ParseUint(c.FormValue("post_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}


	result, err := models.UpdatePost(title, content, image, anonymous, pid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func DeletePost(c echo.Context) error {

	pid := c.FormValue("post_id")

	result, err := models.DeletePost(pid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}