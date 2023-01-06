package controllers

import (
	"net/http"
	"strconv"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func FetchOwnedComments(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.FetchOwnedComments(uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

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

func UpdateComment(c echo.Context) error {
	content := c.FormValue("comment_content")
	id, err := strconv.ParseUint(c.FormValue("comment_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.UpdateComment(content, id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func DeleteComment(c echo.Context) error {
	cid, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.DeleteComment(cid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}
