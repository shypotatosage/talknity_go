package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"talknity/models"

	"github.com/google/uuid"
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

func FetchPosts(c echo.Context) error {

	result, err := models.FetchPosts()

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func FetchOwnedPosts(c echo.Context) error {
	uid := c.Param("user_id")

	result, err := models.FetchOwnedPosts(uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func FetchPost(c echo.Context) error {
	pid, err := strconv.ParseUint(c.Param("post_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.FetchPost(pid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func SearchPosts(c echo.Context) error {
	key := c.Param("search_key")

	result, err := models.SearchPosts(key)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func StorePost(c echo.Context) error {
	title := c.FormValue("post_title")
	content := c.FormValue("post_content")
	image, err := c.FormFile("post_image")

	fileExist := true

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message1": err.Error()})
		} else {
			fileExist = false
		}
	}

	fileName := ""

	if (fileExist) {
		src, err := image.Open()
	
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}
		
		defer src.Close()
	
		fileByte, _ := io.ReadAll(src)
		fileName = "images/post/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(image.Filename, ".")[1]
	
		err = os.WriteFile(fileName, fileByte, 0777)
	
		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}
	}

	anonymous, err := strconv.ParseBool(c.FormValue("anonymous"))

	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	uid, err := strconv.ParseUint(c.FormValue("uid"), 10, 64)

	if err != nil {
		fmt.Println(content + title)
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.StorePost(title, content, fileName, anonymous, uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdatePost(c echo.Context) error {
	title := c.FormValue("post_title")
	content := c.FormValue("post_content")
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

	image, err := c.FormFile("post_image")

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		} else {
			result, err := models.UpdatePost(title, content, "", anonymous, pid)

			if err != nil {
				return c.JSON(http.StatusInternalServerError,
					map[string]string{"message": err.Error()})
			}

			return c.JSON(http.StatusOK, result)
		}
	} else {
		imgTemp, err := models.GetPostImage(pid)

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		src, err := image.Open()

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		defer src.Close()

		fileByte, _ := io.ReadAll(src)
		fileName := "images/post/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(image.Filename, ".")[1]

		err = os.WriteFile(fileName, fileByte, 0777)

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		result, err := models.UpdatePost(title, content, fileName, anonymous, pid)
	
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

func DeletePost(c echo.Context) error {
	pid, err := strconv.ParseUint(c.Param("post_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	image, err := models.GetPostImage(pid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	result, err := models.DeletePost(pid)

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