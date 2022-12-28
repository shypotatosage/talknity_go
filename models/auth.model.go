package models

import (
	"net/http"
	"database/sql"
	"fmt"
	"talknity/db"
	"talknity/helpers"
	
	"github.com/go-playground/validator"
)

type User struct {
	Id       	int `json:"id"`
	Username 	string `json:"user_username" validate:"required,max=50"`
	Displayname string `json:"user_displayname" validate:"required,max=100"`
	Password 	string `json:"-" validate:"required,min=8"`
	Email    	string `json:"user_email" validate:"required"`
	Image	 	string `json:"user_image"`
}

func CheckLogin(username, password string) (bool, User, error) {
	var obj User
	var pwd string
	con := db.CreateCon()

	sqlStatement := "SELECT id, user_username, user_displayname, user_password, user_email FROM users WHERE user_username = ? OR user_email = ?"
	err := con.QueryRow(sqlStatement, username, username).Scan(
		&obj.Id, &obj.Username, &obj.Displayname, &pwd, &obj.Email,
	)

	if err == sql.ErrNoRows {
		fmt.Print("Username not found!")
		
		return false, obj, err
	}

	if err != nil {
		fmt.Print("Query error!")
		
		return false, obj, err
	}

	match, err := helpers.CheckPasswordHash(password, pwd)

	if !match {
		fmt.Print("Hash and password doesn't match!")
		
		return false, obj, err
	}
	
	return true, obj, nil
}

func RegisterUser(user_username string, user_email string, user_password string) (Response, error) {

	var res Response

	v := validator.New()

	usr := User{
		Username: user_username,
		Email: user_email,
		Password: user_password,
		Displayname: user_username,
	}

	err := v.Struct(usr)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT INTO `users`(`user_username`, `user_displayname`, `user_email`, `user_password`, `created_at`, `updated_at`) VALUES (?,?,?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}
		
		return res, err
	}

	result, err := stmt.Exec(user_username, user_username, user_email, user_password)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}
		
		return res, err
	}

	lastInsertedID, err := result.LastInsertId()

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}
		
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = map[string]int64{
		"last_inserted_id": lastInsertedID,
	}

	return res, nil
}

func UserProfile(uid string) (Response, error) {
	var obj User
	var res Response

	con := db.CreateCon()

	sqlStatement := "SELECT id, user_username, user_displayname, user_password, user_email, COALESCE(users.user_image, '') FROM users WHERE id = ?"
	err := con.QueryRow(sqlStatement, uid).Scan(
		&obj.Id, &obj.Username, &obj.Displayname, &obj.Password, &obj.Email, &obj.Image,
	)

	if err != nil {
		return res, err
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = obj
	
	return res, nil
}

func GetUserImage(uid uint64) (string, error) {
	img := ""

	con := db.CreateCon()

	sqlStatement := "SELECT user_image FROM users WHERE id = ?"
	err := con.QueryRow(sqlStatement, uid).Scan(
		&img,
	)

	if err != nil {
		return "", err
	}
	
	return img, nil
}

// Update Data
func UpdateProfile(uid uint64, displayname string, username string, email string, password string, image string) (Response, error) {
	var pwd string

	var res Response

	con := db.CreateCon()

	sqlStatement1 := "SELECT user_password FROM users WHERE id = ?"

	err := con.QueryRow(sqlStatement1, uid).Scan(
		&pwd,
	)

	if err != nil {
		return res, err
	}

	match, err := helpers.CheckPasswordHash(password, pwd)

	if !match {
		fmt.Print("Hash and password doesn't match!")
		
		return res, err
	}

	sqlStatement := ""

	if image != "" {
		sqlStatement = "UPDATE users SET user_displayname=?, user_username=?, user_email=?, user_image=? WHERE id=?"
	} else {
		sqlStatement = "UPDATE users SET user_displayname=?, user_username=?, user_email=? WHERE id=?"
	}

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	if image != "" {
		result, err := stmt.Exec(displayname, username, email, image, uid)
	
		if err != nil {
			return res, err
		}
	
		lastInsertedID, err := result.RowsAffected()
	
		if err != nil {
			return res, err
		}
	
		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = map[string]int64{
			"last_inserted_id2": lastInsertedID,
		}
	
		return res, nil
	} else {
		result, err := stmt.Exec(displayname, username, email, uid)
	
		if err != nil {
			return res, err
		}
	
		lastInsertedID, err := result.RowsAffected()
	
		if err != nil {
			return res, err
		}
	
		res.Status = http.StatusOK
		res.Message = "Success"
		res.Data = map[string]int64{
			"last_inserted_id1": lastInsertedID,
		}
	
		return res, nil
	}
}