package models

import (
	"net/http"
	"talknity/db"

	"github.com/go-playground/validator"
)

type CommunityMember struct {
	Id   uint64 `json:"id"`
	Cid  uint64 `json:"community_id" validate:"required,numeric"`
	User User   `json:"user" validate:"required,numeric"`
}

type CommunityMemberInsert struct {
	Id  uint64 `json:"id"`
	Cid uint64 `json:"community_id" validate:"required,numeric"`
	Uid uint64 `json:"user" validate:"required,numeric"`
}

type CommunityMembers struct {
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
	Members     []CommunityMember `json:"members"`
	Count       uint64            `json:"member_count"`
}

// Read All
func FetchAllCommunityMember(cid uint64) (Response, error) {
	var com CommunityMembers
	var obj CommunityMember
	var cat CommunityCategory
	var usr User
	var ldr User
	var arrObj []CommunityMember
	var res Response

	conn := db.CreateCon()

	sqlStatement := "SELECT community_members.id, community_members.community_id, users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, '') FROM community_members INNER JOIN users ON community_members.user_id = users.id WHERE community_members.community_id = ?"

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

	sqlStatement2 := "SELECT communities.id, communities.community_name, communities.community_description, communities.community_contact, communities.community_logo, communities.community_category_id, communities.user_id, communities.created_at, (SELECT COUNT(community_members.user_id) FROM community_members WHERE community_members.community_id = communities.id), users.id, users.user_username, users.user_displayname, users.user_email, COALESCE(users.user_image, ''), community_categories.id, community_categories.category_name, community_categories.category_logo, community_categories.category_color1, community_categories.category_color2, community_categories.category_color3 FROM communities INNER JOIN users ON communities.user_id = users.id INNER JOIN community_categories ON communities.community_category_id = community_categories.id WHERE communities.id = ?"

	err = conn.QueryRow(sqlStatement2, cid).Scan(&com.Id, &com.Name, &com.Description, &com.Contact, &com.Logo, &com.Cid, &com.Lid, &com.CreatedAt, &com.Count, &ldr.Id, &ldr.Username, &ldr.Displayname, &ldr.Email, &ldr.Image, &cat.Id, &cat.Name, &cat.Logo, &cat.Color1, &cat.Color2, &cat.Color3)

	if err != nil {
		return res, err
	}

	com.Leader = ldr
	com.Category = cat
	com.Members = arrObj

	res.Status = http.StatusOK
	res.Message = "Success"
	res.Data = com

	return res, nil
}

// Delete Data
func DeleteMember(id uint64) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM community_members WHERE id=?"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(id)

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

func SignoutCommunity(uid uint64, cid uint64) (Response, error) {
	var res Response

	con := db.CreateCon()

	sqlStatement := "DELETE FROM community_members WHERE user_id=? AND community_id=?"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		return res, err
	}

	result, err := stmt.Exec(uid, cid)

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

// Insert Data
func JoinCommunity(cid uint64, uid uint64) (Response, error) {

	var res Response

	v := validator.New()

	member := CommunityMemberInsert {
		Cid: cid,
		Uid: uid,
	}

	err := v.Struct(member)

	if err != nil {
		res.Status = http.StatusBadRequest
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	con := db.CreateCon()

	sqlStatement := "INSERT INTO `community_members`(`community_id`, `user_id`, `created_at`, `updated_at`) VALUES (?,?,NOW(),NOW())"
	stmt, err := con.Prepare(sqlStatement)

	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = "Error"
		res.Data = map[string]string{
			"errors": err.Error(),
		}

		return res, err
	}

	result, err := stmt.Exec(cid, uid)

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
