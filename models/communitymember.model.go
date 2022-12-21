package models

import (
	"net/http"
	"talknity/db"
)

type CommunityMember struct {
	Id   uint64 `json:"id" validate:"required,numeric"`
	Cid  uint64 `json:"community_id" validate:"required,numeric"`
	User User   `json:"user" validate:"required,numeric"`
}

// Read All
func FetchAllCommunityMember(cid uint64) (Response, error) {
	var obj CommunityMember
	var usr User
	var arrObj []CommunityMember
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT community_members.id, community_members.community_id, users.id, users.user_username, users.user_displayname, users.user_email, users.user_image FROM community_members INNER JOIN users ON community_members.user_id = users.id WHERE community_members.community_id = ?"

	rows, err := conn.Query(sqlStatement, cid)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Cid, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
		obj.User = usr

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
