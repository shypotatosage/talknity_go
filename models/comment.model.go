package models

import (
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type Comment struct {
	Id        uint64 `json:"id" validate:"required,numeric"`
	Content   string `json:"comment_content" validate:"required"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user" validate:"required"`
}

type CommentInsert struct {
	Id        uint64 `json:"id"`
	Pid       uint64 `json:"post_id" validate:"required,numeric"`
	Content   string `json:"comment_content" validate:"required"`
	CreatedAt string `json:"created_at"`
	Uid       uint64 `json:"user_id" validate:"required"`
}

// Read All
func FetchAllComment(pid uint64) (Response, []Comment, error) {
	var obj Comment
	var usr User
	var arrObj []Comment
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT comments.id, comments.comment_content, comments.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM comments INNER JOIN users ON comments.user_id = users.id WHERE comments.post_id = ?"

	rows, err := conn.Query(sqlStatement, pid)

	if err != nil {
		return res, arrObj, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Content, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
		obj.User = usr

		if err != nil {
			return res, arrObj, err
		}

		arrObj = append(arrObj, obj)
	}

	return res, arrObj, nil
}

// Insert Data
func StoreComment(post_id uint64, content string, user_id uint64) (Response, error) {

	var res Response

	v := validator.New()

	comment := CommentInsert{
		Pid:     post_id,
		Content: content,
		Uid:     user_id,
	}

	err := v.Struct(comment)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT INTO `comments`(`comment_content`, `post_id`, `user_id`, `created_at`, `updated_at`) VALUES (?,?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	result, err := stmt.Exec(content, post_id, user_id)

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
