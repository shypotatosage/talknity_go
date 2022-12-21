package models

import (
	"net/http"
	"talknity/db"
)

type Comment struct {
	Id        uint64 `json:"id" validate:"required,numeric"`
	Pid       uint64 `json:"post_id" validate:"required,numeric"`
	Content   string `json:"user_image" validate:"required"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user" validate:"required"`
}

// Read All
func FetchAllComment(pid uint64) (Response, error) {
	var obj Comment
	var usr User
	var arrObj []Comment
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT comments.id, comments.post_id, comments.comment_content, comments.created_at, users.id, users.user_username, users.user_displayname, users.user_email, users.user_image FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ?"

	rows, err := conn.Query(sqlStatement, pid)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Pid, &obj.Content, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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
