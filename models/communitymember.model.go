package models

import (
	"net/http"
	"talknity/db"
)

type CommunityMember struct {
	Id  uint64 `json:"id" validate:"required,numeric"`
	Uid uint64 `json:"user_id" validate:"required,numeric"`
	Cid uint64 `json:"community_id" validate:"required,numeric"`
}

// Read All
func FetchAllCommunityMember(cid uint64) (Response, error) {
	var obj CommunityMember
	var arrObj []CommunityMember
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT id, user_id, community_id FROM community_members WHERE"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Uid, &obj.Cid)

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
