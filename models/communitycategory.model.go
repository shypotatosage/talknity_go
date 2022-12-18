package models

import (
	"net/http"
	"talknity/db"
)

type CommunityCategory struct {
	Id      	 int 	`json:"id" validate:"required,numeric"`
	Name     	 string `json:"category_name" validate:"required,max=100"`
	Logo   	 	 string `json:"category_logo" validate:"required"`
	Bg 		 	 string `json:"category_bg" validate:"required"`
}

// Read All
func FetchAllCommunityCategory() (Response, error) {
	var obj CommunityCategory
	var arrObj []CommunityCategory
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT * FROM community_categories"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Logo, &obj.Bg)

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