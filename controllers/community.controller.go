package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"talknity/models"

	"github.com/google/uuid"
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

func FetchCommunities(c echo.Context) error {

	result, err := models.FetchCommunity()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func FetchOwnedCommunities(c echo.Context) error {
	uid := c.Param("user_id")

	result, err := models.FetchOwnedCommunity(uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func FetchCommunitiesCategory(c echo.Context) error {
	cid := c.Param("category_id")

	result, err := models.FetchCommunitiesCategory(cid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func SearchCommunitiesCategory(c echo.Context) error {
	cid := c.Param("category_id")
	key := c.Param("search_key")

	result, err := models.SearchCommunitiesCategory(cid, key)

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
	fileName := "images/community/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(logo.Filename, ".")[1]

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

func UpdateCommunity(c echo.Context) error {
	name := c.FormValue("community_name")
	description := c.FormValue("community_description")
	contact := c.FormValue("community_contact")
	cid, err := strconv.ParseUint(c.FormValue("category_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	id, err := strconv.ParseUint(c.FormValue("community_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	logo, err := c.FormFile("community_logo")

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		} else {
			result, err := models.UpdateCommunity(name, description, contact, "", cid, id)

			if err != nil {
				return c.JSON(http.StatusInternalServerError,
					map[string]string{"message": err.Error()})
			}

			return c.JSON(http.StatusOK, result)
		}
	} else {
		imgTemp, err := models.GetCommunityImage(id)

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
		fileName := "images/community/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(logo.Filename, ".")[1]

		err = os.WriteFile(fileName, fileByte, 0777)

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		result, err := models.UpdateCommunity(name, description, contact, fileName, cid, id)
	
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		if imgTemp != "" {
			err := os.Remove(imgTemp)

			if err != nil {
				return c.JSON(http.StatusInternalServerError,
					map[string]string{"message": err.Error()})
			}
		}
	
		return c.JSON(http.StatusOK, result)
	}
}

func DeleteCommunity(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("community_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	image, err := models.GetCommunityImage(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.DeleteCommunity(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	if image != "" {
		err := os.Remove(image)

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, result)
}