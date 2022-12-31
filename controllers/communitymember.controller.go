package controllers

import (
	"net/http"
	"strconv"
	"talknity/models"

	"github.com/labstack/echo/v4"
)

func FetchAllCommunityMember(c echo.Context) error {

	cid, err := strconv.ParseUint(c.Param("community_id"), 10, 64)

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

func DeleteMember(c echo.Context) error {

	id, err := strconv.ParseUint(c.Param("community_member_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.DeleteMember(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func SignoutCommunity(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("user_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	cid, err := strconv.ParseUint(c.Param("community_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.SignoutCommunity(uid, cid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func JoinCommunity(c echo.Context) error {
	cid, err := strconv.ParseUint(c.FormValue("community_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	uid, err := strconv.ParseUint(c.FormValue("user_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.JoinCommunity(cid, uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}