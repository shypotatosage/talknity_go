package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type Community struct {
	Id          uint64            `json:"id" validate:"numeric"`
	Name        string            `json:"community_name" validate:"required,max=100"`
	Description string            `json:"community_description" validate:"required"`
	Contact     string            `json:"community_contact" validate:"required"`
	Logo        string            `json:"community_logo" validate:"required"`
	Cid         uint64            `json:"-" validate:"required,numeric"`
	Lid         uint64            `json:"-" validate:"required,numeric"`
	CreatedAt   string            `json:"created_at"`
	Leader      User              `json:"leader" validate:"required"`
	Category    CommunityCategory `json:"category" validate:"required"`
	Count       uint64            `json:"member_count"`
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

	sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id ORDER BY members_count DESC"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
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

// Read 10
func FetchCommunity() (Response, error) {
	var obj Community
	var usr User
	var cat CommunityCategory
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id ORDER BY members_count DESC LIMIT 10"

	rows, err := conn.Query(sqlStatement)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
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

// Read Owned
func FetchOwnedCommunity(uid string) (Response, error) {
	var obj Community
	var usr User
	var cat CommunityCategory
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id WHERE communities.user_id = ? ORDER BY communities.id DESC"

	rows, err := conn.Query(sqlStatement, uid)

	if err != nil {
		return res, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
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

// Read Based On Category
func FetchCommunitiesCategory(cid string) (Response, error) {
	var obj Community
	var usr User
	var cat CommunityCategory
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	if cid == "0" {
		sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id ORDER BY members_count DESC"

		rows, err := conn.Query(sqlStatement)

		if err != nil {
			return res, err
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
			obj.Leader = usr
			obj.Category = cat

			if err != nil {
				return res, err
			}

			arrObj = append(arrObj, obj)
		}
	} else {
		sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id WHERE communities.community_category_id = ? ORDER BY members_count DESC"

		rows, err := conn.Query(sqlStatement, cid)

		if err != nil {
			return res, err
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
			obj.Leader = usr
			obj.Category = cat

			if err != nil {
				return res, err
			}

			arrObj = append(arrObj, obj)
		}
	}

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = arrObj

	return res, nil
}

// Search Based On Category
func SearchCommunitiesCategory(cid string, key string) (Response, error) {
	var obj Community
	var usr User
	var cat CommunityCategory
	var arrObj []Community
	var res Response

	conn := db.CreateCon()

	if cid == "0" {
		sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id WHERE communities.community_name LIKE ? OR communities.community_description LIKE ? ORDER BY members_count DESC"

		rows, err := conn.Query(sqlStatement, "%"+key+"%", "%"+key+"%")

		if err != nil {
			return res, err
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
			obj.Leader = usr
			obj.Category = cat

			if err != nil {
				return res, err
			}

			arrObj = append(arrObj, obj)
		}
	} else {
		sqlStatement := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id) as members_count, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id WHERE communities.community_category_id = ? AND (communities.community_name LIKE ? OR communities.community_description LIKE ?) ORDER BY members_count DESC"

		rows, err := conn.Query(sqlStatement, cid, "%"+key+"%", "%"+key+"%")

		if err != nil {
			return res, err
		}

		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&obj.Id, &obj.Name, &obj.Description, &obj.Contact, &obj.Logo, &obj.Cid, &obj.Lid, &obj.CreatedAt, &obj.Count, &usr.Id, &usr.Username, &usr.Displayname, &usr.Email, &usr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)
			obj.Leader = usr
			obj.Category = cat

			if err != nil {
				return res, err
			}

			arrObj = append(arrObj, obj)
		}
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
		Name:        name,
		Description: description,
		Contact:     contact,
		Logo:        logo,
		Cid:         cid,
		Lid:         lid,
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

	sqlStatement := "INSERT INTO `communities`(`community_name`, `community_description`, `community_contact`, `community_logo`, `community_category_id`, `user_id`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,NOW(),NOW())"
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
func UpdateCommunity(name string, description string, contact string, logo string, cid uint64, id uint64) (Response, error) {

	var res Response

	con := db.CreateCon()

	sqlStatement := ""

	if logo == "" {
		sqlStatement = "UPDATE communities SET community_name=?, community_description=?, community_contact=?, community_category_id=?, updated_at=NOW() WHERE id=?"
	} else {
		sqlStatement = "UPDATE communities SET community_name=?, community_description=?, community_contact=?, community_logo=?, community_category_id=?, updated_at=NOW() WHERE id=?"
	}

	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	var result sql.Result

	if logo == "" {
		result, err = stmt.Exec(name, description, contact, cid, id)
	} else {
		result, err = stmt.Exec(name, description, contact, logo, cid, id)
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
func DeleteCommunity(cid uint64) (Response, error) {

	var res Response

	con := db.CreateCon()

	sqlStatement2 := "DELETE FROM `community_members` WHERE community_id=?"
	stmt2, err := con.Prepare(sqlStatement2)

	if err != nil {
		return res, err
	}

	result2, err := stmt2.Exec(cid)

	if err != nil {
		return res, err
	}

	fmt.Println(result2.RowsAffected())

	sqlStatement := "DELETE FROM `communities` WHERE id=?"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(cid)

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

func GetCommunityImage(cid uint64) (string, error) {
	img := ""

	con := db.CreateCon()

	sqlStatement := "SELECT community_logo FROM communities WHERE id = ?"
	err := con.QueryRow(sqlStatement, cid).Scan(
		&img,
	)

	if err != nil {
		return "", err
	}
	
	return img, nil
}