package controllers

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"talknity/helpers"
	"talknity/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GenerateHashPassword(c echo.Context) error {
	password := c.Param("user_password")
	hash, _ := helpers.HashPassword(password)

	return c.JSON(http.StatusOK, hash)
}

func CheckLogin(c echo.Context) error {
	usernameemail := c.FormValue("user_usernameemail")
	password := c.FormValue("user_password")

	res, obj, err := models.CheckLogin(usernameemail, password)

	if err != nil {
		errMsg := err.Error()

		if errMsg == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			errMsg = "Password does not match!"
		} else if err.Error() == "sql: no rows in result set" {
			errMsg = "Username/Email does not exist!"
		}

		return c.JSON(http.StatusOK,
			map[string]string{
				"message": errMsg,
			})
	}

	if !res {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["usernameemail"] = usernameemail
	claims["level"] = "application"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	mytoken, err := token.SignedString([]byte("secret"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK,
		map[string]string{
			"message":          "Login successful",
			"token":            mytoken,
			"user_id":          strconv.Itoa(obj.Id),
			"user_username":    obj.Username,
			"user_displayname": obj.Displayname,
			"user_image":       obj.Image,
			"user_email":       obj.Email,
		})
}

func RegisterUser(c echo.Context) error {

	username := c.FormValue("user_username")
	email := c.FormValue("user_email")
	password, _ := helpers.HashPassword(c.FormValue("user_password"))

	result, err := models.RegisterUser(username, email, password)

	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry '" + username + "' for key 'users_user_username_unique'" {
			return c.JSON(http.StatusBadGateway,
				result)
		} else if err.Error() == "Error 1062: Duplicate entry '" + email + "' for key 'users_user_email_unique'" {
			return c.JSON(http.StatusBadRequest,
				result)
		} else {
			return c.JSON(http.StatusInternalServerError,
				result)
		}
	}

	return c.JSON(http.StatusOK, result)
}

func UserProfile(c echo.Context) error {
	uid := c.FormValue("user_id")

	result, err := models.UserProfile(uid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func UpdateProfile(c echo.Context) error {
	uid, err := strconv.ParseUint(c.FormValue("user_id"), 10, 64)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"message": err.Error()})
	}

	username := c.FormValue("user_username")
	displayname := c.FormValue("user_displayname")
	email := c.FormValue("user_email")
	password := c.FormValue("user_password")
	image, err := c.FormFile("user_image")

	if err != nil {
		if err != http.ErrMissingFile {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		} else {
			result, err := models.UpdateProfile(uid, displayname, username, email, password, "")

			if err != nil {
				errMsg := err.Error()
	
				if errMsg == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
					errMsg = "Password does not match!"
	
					return c.JSON(http.StatusBadRequest,
						map[string]string{"message": errMsg})
				} else {
					return c.JSON(http.StatusInternalServerError,
						map[string]string{"message": errMsg})
				}
			}

			return c.JSON(http.StatusOK, result)
		}
	} else {
		imgTemp, err := models.GetUserImage(uid)

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
		fileName := "images/user/" + strings.Replace(uuid.New().String(), "-", "", -1) + "." + strings.Split(image.Filename, ".")[1]

		err = os.WriteFile(fileName, fileByte, 0777)

		if err != nil {
			return c.JSON(http.StatusInternalServerError,
				map[string]string{"message": err.Error()})
		}

		result, err := models.UpdateProfile(uid, displayname, username, email, password, fileName)

		if err != nil {
			errMsg := err.Error()

			if errMsg == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
				errMsg = "Password does not match!"

				return c.JSON(http.StatusBadRequest,
					map[string]string{"message": errMsg})
			} else {
				return c.JSON(http.StatusInternalServerError,
					map[string]string{"message": errMsg})
			}
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
