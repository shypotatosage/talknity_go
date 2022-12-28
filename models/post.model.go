package models

import (
	"database/sql"
	"fmt"
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
	Count     uint64 `json:"comment_count" validate:"required"`
}

type PostInsert struct {
	Title     string `json:"post_title" validate:"required,max=255"`
	Content   string `json:"post_content" validate:"required"`
	Image     string `json:"post_image"`
	Anonymous bool   `json:"anonymous"`
	UId       uint64 `json:"-" validate:"required,numeric"`
}

type PostComment struct {
	Id        int       `json:"id" validate:"numeric"`
	Title     string    `json:"post_title" validate:"required,max=255"`
	Content   string    `json:"post_content" validate:"required"`
	Image     string    `json:"post_image"`
	Anonymous bool      `json:"anonymous"`
	UId       uint64    `json:"-" validate:"required,numeric"`
	CreatedAt string    `json:"created_at"`
	User      User      `json:"user"`
	Comments  []Comment `json:"comments"`
}

// Read All
func FetchAllPosts() (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, (SELECT COUNT(comments.user_id) FROM comments WHERE comments.post_id = posts.id), users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, (SELECT COUNT(comments.user_id) FROM comments WHERE comments.post_id = posts.id) as comments_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id ORDER BY comments_count DESC LIMIT 10"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Read Owned
func FetchOwnedPosts(uid string) (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, (SELECT COUNT(comments.user_id) FROM comments WHERE comments.post_id = posts.id), users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.user_id = ?"

	rows, err := conn.Query(sqlStatement, uid)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

// Read Individual Post
func FetchPost(pid uint64) (Response, error) {
	var obj PostComment
	var usr User
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.id = ? LIMIT 1"

	err := conn.QueryRow(sqlStatement, pid).Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)

	if err != nil {
		return res, err
	}

	obj.User = usr

	res, com, err := FetchAllComment(pid)

	if err != nil {
		return res, err
	}

	obj.Comments = com

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = obj

	return res, nil
}

// Search
func SearchPosts(key string) (Response, error) {
	var obj Post
	var usr User
	var arrObj []Post
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT posts.id, posts.post_title, posts.post_content, posts.post_image, posts.anonymous, posts.user_id, posts.created_at, (SELECT COUNT(comments.user_id) FROM comments WHERE comments.post_id = posts.id), users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM posts INNER JOIN users ON posts.user_id = users.id WHERE posts.post_title LIKE ? OR posts.post_content LIKE ? OR (users.user_displayname LIKE ? AND posts.anonymous = false)"

	rows, err := conn.Query(sqlStatement, "%"+key+"%", "%"+key+"%", "%"+key+"%")

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Title, &obj.Content, &obj.Image, &obj.Anonymous, &obj.UId, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
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

	sqlStatement := ""

	if image == "" {
		sqlStatement = "UPDATE posts SET post_title=?, post_content=?, anonymous=?, updated_at=NOW() WHERE id=?"
	} else {
		sqlStatement = "UPDATE posts SET post_title=?, post_content=?, post_image=?, anonymous=?, updated_at=NOW() WHERE id=?"
	}

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	var result sql.Result

	if image == "" {
		result, err = stmt.Exec(title, content, anonymous, id)
	} else {
		result, err = stmt.Exec(title, content, image, anonymous, id)
	}

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
func DeletePost(pid uint64) (Response, error) {

	var res Response

	con := db.CreateCon()

	sqlStatement2 := "DELETE FROM `comments` WHERE post_id=?"
	stmt2, err := con.Prepare(sqlStatement2)

	if err != nil {
		return res, err
	}

	result2, err := stmt2.Exec(pid)

	if err != nil {
		return res, err
	}

	fmt.Println(result2.RowsAffected())

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

func GetPostImage(pid uint64) (string, error) {
	img := ""

	con := db.CreateCon()

	sqlStatement := "SELECT post_image FROM posts WHERE id = ?"
	err := con.QueryRow(sqlStatement, pid).Scan(
		&img,
	)

	if err != nil {
		return "", err
	}
	
	return img, nil
}