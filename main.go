package main

import (
	merry "github.com/Mondal-Prasun/BloodBank/Merry"
)

func main() {

	staticDir := "static"

	m := merry.Init("/v1", &staticDir)

	m.Route(merry.GET, "/hi", func(mr merry.MerryContext) {
		mr.Res(200, struct {
			Msg string `json:"msg"`
		}{
			Msg: "this is for get method",
		})
	})

	m.Route(merry.POST, "/posthi", func(mr merry.MerryContext) {
		mr.Res(201, struct {
			Msg string `json:"msg"`
		}{
			Msg: "this is for post",
		})
	})

	m.Route(merry.GET, "/home", func(mr merry.MerryContext) {
		mr.ServeHtml(202, "/home.html")
	})

	merry.Ship("8080", m)

}
