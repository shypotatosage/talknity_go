package controllers

import (
	"net/http"
	"strconv"
	"talknity/helpers"
	"talknity/models"
	"time"

	"github.com/dgrijalva/jwt-go"
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

		if (err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password") {
			errMsg = "Password does not Match!"
		} else if (err.Error() == "sql: no rows in result set") {
			errMsg = "Username/Email does not exist!"
		}
		
		return c.JSON(http.StatusOK,
			map[string]string {
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
			"message": "Login successful",
			"token": mytoken,
			"user_id": strconv.Itoa(obj.Id),
			"user_username": obj.Username,
			"user_displayname": obj.Displayname,
			"user_image": obj.Image,
			"user_email": obj.Email,
		})
}

func RegisterUser(c echo.Context) error {

	username := c.FormValue("user_username")
	email := c.FormValue("user_email")
	password, _ := helpers.HashPassword(c.FormValue("user_password"))

	result, err := models.RegisterUser(username, email, password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			result)
	}

	return c.JSON(http.StatusOK, result)
}