package models

import (
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type Post struct {
	Id        int    `json:"id" validate:"numeric"`
	Title     string `json:"post_title" validate:"required,max=255"`
	Content   string `json:"post_content" validate:"required"`
	Image     string `json:"post_image"`
	Anonymous bool   `json:"anonymous"`
	UId       uint64 `json:"-" validate:"required,numeric"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type PostInsert struct {
	Title     string `json:"post_title" validate:"required,max=255"`
	Content   string `json:"post_content" validate:"required"`
	Image     string `json:"post_image"`
	Anonymous bool   `json:"anonymous"`
	UId       uint64 `json:"-" validate:"required,numeric"`
}

// Read All
func FetchAllPosts() (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Read 10
func FetchPosts() (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id LIMIT 10"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Read 10
func FetchOwnedPosts(uid string) (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.user_id = ? LIMIT 10"

	rows, err := conn.Query(sqlStatement, uid)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Search
func SearchPosts(key string) (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.post_title LIKE ? OR posts.post_content LIKE ? OR (users.user_displayname LIKE ? AND posts.anonymous = false)"

	rows, err := conn.Query(sqlStatement, "%" + key + "%", "%" + key + "%", "%" + key + "%")

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Insert Data
func StorePost(title string, content string, image string, anonymous bool, uid uint64) (Response, error) {

	var res Response

	v := validator.New()

	post := PostInsert{
		Title:     title,
		Content:   content,
		Image:     image,
		Anonymous: anonymous,
		UId:       uid,
	}

	err := v.Struct(post)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT INTO `posts`(`post_title`, `post_content`, `post_image`, `anonymous`, `user_id`, `created_at`, `updated_at`) VALUES (?,?,?,?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	result, err := stmt.Exec(title, content, image, anonymous, uid)

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

// Update Data
func UpdatePost(title string, content string, image string, anonymous bool, id uint64) (Response, error) {

	var res Response

	con := db.CreateCon()

	sqlStatement := "UPDATE posts SET post_title=?, post_content=?, post_image=?, anonymous=?, updated_at=NOW() WHERE id=?"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(title, content, image, anonymous, id)

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
		"last_inserted_id": lastInsertedID,
	}

	return res, nil
}

// Delete Data
func DeletePost(pid string) (Response, error) {

	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM `posts` WHERE id=?"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(pid)

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
		"last_inserted_id": lastInsertedID,
	}

	return res, nil
}
