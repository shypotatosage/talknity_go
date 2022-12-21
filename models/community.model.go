package models

import (
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type Community struct {
	Id          uint64               `json:"id" validate:"numeric"`
	Name        string            `json:"community_name" validate:"required,max=100"`
	Description string            `json:"community_description" validate:"required"`
	Contact     string            `json:"community_contact" validate:"required"`
	Logo        string            `json:"community_logo" validate:"required"`
	Cid         uint64            `json:"-" validate:"required,numeric"`
	Lid         uint64            `json:"-" validate:"required,numeric"`
	Leader      User              `json:"leader" validate:"required"`
	CreatedAt   string            `json:"created_at"`
	Category    CommunityCategory `json:"category" validate:"required"`
}

type CommunityInsert struct {
	Name        string `json:"community_name" validate:"required,max=100"`
	Description string `json:"community_description" validate:"required"`
	Contact     string `json:"community_contact" validate:"required"`
	Logo        string `json:"community_logo" validate:"required"`
	Cid         uint64 `json:"-" validate:"required,numeric"`
	Lid         uint64 `json:"-" validate:"required,numeric"`
}

// Read All
func FetchAllCommunities() (Response, error) {
	var obj Community
	var usr User
	var cat CommunityCategory
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, users.id, users.user_username, users.user_displayname, users.user_email, users.user_image, community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_bg FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Bg)
		obj.Leader = usr
		obj.Category = cat

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

// Read All
func FetchCommunity() (Response, error) {
	var obj Community
	var usr User
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, posts.created_at, users.id, users.user_username, users.user_displayname, users.user_email, users.user_image users FROM communities INNER JOIN users ON communities.user_id = users.id"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image)
		obj.Leader = usr

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
func StoreCommunity(name string, description string, contact string, logo string, cid uint64, lid uint64) (Response, error) {

	var res Response

	v := validator.New()

	post := CommunityInsert{
		Name: name,
		Description: description,
		Contact: contact,
		Logo: logo,
		Cid: cid,
		Lid: lid,
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

	sqlStatement := "INSERT INTO `communities`(`community_name`, `community_description`, `community_contact`, `community_logo`, `community_category_id`, `user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	result, err := stmt.Exec(name, description, contact, logo, cid, lid)

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
func UpdateCommunity(title string, content string, image string, anonymous bool, id uint64) (Response, error) {

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
func DeleteCommunity(pid string) (Response, error) {

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
