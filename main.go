package main

import (
	"encoding/json"

	merry "github.com/Mondal-Prasun/BloodBank/Merry"
	"github.com/google/uuid"
)

func main() {

	staticDir := "static"

	udb, cdb := initDatabase()

	defer udb.Close()
	defer cdb.Close()

	dbCfg := DbCfg{
		udb: udb,
		cdb: cdb,
	}

	m := merry.Init("/v1", &staticDir)

	m.Route(merry.GET, "/hi", func(mr merry.MerryContext) {
		mr.Res(200, struct {
			Msg string `json:"msg"`
		}{
			Msg: "this is for get method",
		})
	})

	m.Route(merry.GET, "/home", func(mr merry.MerryContext) {
		mr.ServeHtml(200, "/home.html")
	})
	m.Route(merry.GET, "/apply", func(mr merry.MerryContext) {
		mr.ServeHtml(200, "/apply.html")
	})

	m.Route(merry.GET, "/donors", func(mr merry.MerryContext) {
		mr.ServeHtml(200, "/donars.html")
	})

	m.Route(merry.POST, "/insertUser", func(mr merry.MerryContext) {

		data := struct {
			Phone     string `json:"phone"`
			Name      string `json:"fullName"`
			Gender    string `json:"gender"`
			BloodType string `json:"bloodType"`
			CampId    string `json:"campId"`
			Weight    int    `json:"weight"`
		}{}

		body := mr.ReqBody()

		decoder := json.NewDecoder(body)

		err := decoder.Decode(&data)

		if err != nil {
			mr.Err(300, err.Error())
			return
		}

		done, err := dbCfg.insertUserTableDonor(&User{
			ID:        uuid.New(),
			Phone:     data.Phone,
			Name:      data.Name,
			Gender:    data.Gender,
			BloodType: data.BloodType,
			CampId:    data.CampId,
			Weight:    data.Weight,
		})

		if err != nil {
			mr.Err(300, err.Error())
			return
		}

		if done {
			mr.Res(200, struct {
				Created bool `json:"created"`
			}{
				Created: done,
			})
		}

	})

	m.Route(merry.GET, "/allCamp", func(mr merry.MerryContext) {
		data, err := dbCfg.getAllTableCamp()

		if err != nil {
			mr.Err(500, err.Error())
			return
		}

		mr.Res(200, data)

	})

	merry.Ship("8080", m)

}
