package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
)

func createUserHandler(r *http.Request, web *Web, ds *Ds) (int, string) {
	name := strings.Replace(r.PostFormValue("name"), " ", "", -1)
	if name == "" {
		return 400, "name required"
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
			return 400, "bad isAddrVerified"
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
			return 400, "已存在该手机号的用户, 请重新填写"
		}

		return 500, fmt.Sprintf("insert user error : %v", err)
	}

	return web.Json(200, user)
}

func listUserHandler(r *http.Request, web *Web, ds *Ds) (int, string) {
	page, err := parseIntParam(r, "page", 1)
	if err != nil {
		return 400, err.Error()
	}

	size, err := parseIntParam(r, "size", 10)
	if err != nil {
		return 400, err.Error()
	}

	SPEC := bson.M{}

	skip := (page - 1) * size
	l, total, err := findUserByQuery(ds, SPEC, skip, size)
	chk(err)
	return web.Json(200, J{"data": l, "total": total, "page": page, "size": size})
}

func delUserHandler(ds *Ds, web *Web, user *User) (int, string) {
	if err := delUserById(ds, user.Id); err != nil {
		return 500, "del user err"
	}

	return 200, "ok"
}

func showUserHandler(ds *Ds, web *Web, user *User) (int, string) {
	return web.Json(200, user)
}

func updateUserHandler(r *http.Request, user *User, web *Web, ds *Ds) (int, string) {
	name := strings.Replace(r.PostFormValue("name"), " ", "", -1)
	if name == "" {
		return 400, "name required"
	}

	address := strings.TrimSpace(r.PostFormValue("address"))
	mobile := strings.TrimSpace(r.PostFormValue("mobile"))
	desc := strings.TrimSpace(r.PostFormValue("desc"))

	sex, err := strconv.ParseBool(strings.TrimSpace(r.PostFormValue("sex")))
	if err != nil {
		return 400, "bad isAddrVerified"
	}

	user.Name = name
	user.Sex = sex
	user.Address = address
	user.Mobile = mobile
	user.Desc = desc
	user.Mt = tick()

	err = ds.se.DB(DB).C(C_USER).UpdateId(user.Id, user)
	if err != nil {
		if dup(err) {
			return 400, "已存在手机号,请重新填写!"
		}
		return 500, "update user err"
	}

	return web.Json(200, user)
}
