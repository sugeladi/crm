package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
)

func createUserHandler(r *http.Request, web *Web, ds *Ds) (int, string) {
	rsp := &Rsp{
		Code: 0,
	}

	name := strings.Replace(r.PostFormValue("name"), " ", "", -1)
	if name == "" {
		rsp.Data = "name required"
		return web.Json(200, rsp)
	}

	address := strings.TrimSpace(r.PostFormValue("address"))
	mobile := strings.TrimSpace(r.PostFormValue("mobile"))
	desc := strings.TrimSpace(r.PostFormValue("desc"))

	var err error
	sex := false
	isAdd := r.PostFormValue("isAddrVerified")
	if isAdd == "" {
		sex = false
	} else {
		sex, err = strconv.ParseBool(strings.TrimSpace(isAdd))
		if err != nil {
			rsp.Data = "bad isAddrVerified"
			return web.Json(200, rsp)
		}
	}

	user := &User{
		Id:       newId().Hex(),
		Name:     name,
		Password: "123456",
		Sex:      sex,
		Address:  address,
		Mobile:   mobile,
		Desc:     desc,
		Ct:       tick(),
		Mt:       tick(),
	}

	err = addUser(ds, user)
	if err != nil {
		if dup(err) {
			rsp.Data = "已存在该手机号的用户, 请重新填写"
			return web.Json(200, rsp)
		}

		rsp.Data = fmt.Sprintf("insert user error : %v", err)
		return web.Json(200, rsp)
	}

	rsp.Data = "insert user success"
	rsp.Code = 1
	return web.Json(200, rsp)
}

func listUserHandler(r *http.Request, web *Web, ds *Ds) (int, string) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		return 400, "bad page: " + r.FormValue("page")
	}

	size, err := strconv.Atoi(r.FormValue("size"))
	if err != nil {
		return 400, "bad size: " + r.FormValue("size")
	}

	SPEC := bson.M{}

	countryId := r.FormValue("countryId")
	if countryId != "" && bson.IsObjectIdHex(countryId) {
		SPEC["ex.country"] = countryId
	}

	skip := (page - 1) * size
	l, total, err := findUserByQuery(ds, SPEC, skip, size)
	chk(err)
	return web.Json(200, J{"users": l, "total": total, "page": page, "size": size})
}

func delUserHandler(ds *Ds, web *Web, user *User) (int, string) {
	rsp := &Rsp{
		Code: 0,
	}

	if err := delUserById(ds, user.Id); err != nil {
		rsp.Data = fmt.Sprintf("del user error : %v", err)
		return web.Json(200, rsp)
	}

	rsp.Code = 1
	rsp.Data = "delete user successfully"
	return web.Json(200, rsp)
}

func showUserHandler(ds *Ds, web *Web, user *User) (int, string) {
	return web.Json(200, user)
}

func updateUserHandler(r *http.Request, user *User, web *Web, ds *Ds) (int, string) {
	rsp := &Rsp{
		Code: 0,
	}

	name := strings.Replace(r.PostFormValue("name"), " ", "", -1)
	if name == "" {
		rsp.Data = "name required"
		return web.Json(200, rsp)
	}

	address := strings.TrimSpace(r.PostFormValue("address"))
	mobile := strings.TrimSpace(r.PostFormValue("mobile"))
	desc := strings.TrimSpace(r.PostFormValue("desc"))

	sex, err := strconv.ParseBool(strings.TrimSpace(r.PostFormValue("sex")))
	if err != nil {
		rsp.Data = "bad isAddrVerified"
		return web.Json(200, rsp)
	}

	user.Name = name
	user.Sex = sex
	user.Address = address
	user.Mobile = mobile
	user.Desc = desc
	user.Mt = tick()

	err = ds.se.DB(DB).C(C_USER).UpdateId(user.Id, user)
	chk(err)

	rsp.Code = 1
	rsp.Data = "update user successfully"
	return web.Json(200, rsp)
}
