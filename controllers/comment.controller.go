package controllers

import (
	"net/http"
	"strconv"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func StoreComment(c echo.Context) error {
	pid, err := strconv.ParseUint(c.FormValue("post_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	content := c.FormValue("comment_content")
	
	uid, err := strconv.ParseUint(c.FormValue("user_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.StoreComment(pid, content, uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}