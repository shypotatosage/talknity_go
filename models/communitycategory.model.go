package models

import (
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type CommunityCategory struct {
	Id     int    `json:"id"`
	Name   string `json:"category_name" validate:"required,max=100"`
	Logo   string `json:"category_logo" validate:"required"`
	Color1 string `json:"category_color1" validate:"required"`
	Color2 string `json:"category_color2" validate:"required"`
	Color3 string `json:"category_color3" validate:"required"`
}

// Read All
func FetchAllCommunityCategory() (Response, error) {
	var obj CommunityCategory
	var arrObj []CommunityCategory
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM community_categories ORDER BY community_categories.category_name ASC"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Logo, &obj.Color1, &obj.Color2, &obj.Color3)

		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

// Search
func SearchCommunityCategory(key string) (Response, error) {
	var obj CommunityCategory
	var arrObj []CommunityCategory
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM community_categories WHERE community_categories.category_name LIKE ? ORDER BY community_categories.category_name ASC"

	rows, err := conn.Query(sqlStatement, "%" + key + "%")

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Logo, &obj.Color1, &obj.Color2, &obj.Color3)

		if err != nil {
			return res, err
		}

		arrObj = append(arrObj, obj)
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

// Insert Data
func StoreCommunityCategory(name string, logo string, color1 string, color2 string, color3 string) (Response, error) {

	var res Response

	v := validator.New()

	category := CommunityCategory{
		Name:   name,
		Logo:   logo,
		Color1: color1,
		Color2: color2,
		Color3: color3,
	}

	err := v.Struct(category)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT INTO `community_categories`(`category_name`, `category_logo`, `category_color1`, `category_color2`, `category_color3`, `created_at`, `updated_at`) VALUES (?,?,?,?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	result, err := stmt.Exec(name, logo, color1, color2, color3)

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
