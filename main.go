package main

import (
	"fmt"

	merry "github.com/Mondal-Prasun/BloodBank/Merry"
)

func main() {
	fmt.Println("This is new test")

	m := merry.Init("/v1", "")

	m.Route("/hi", func(mr merry.MerryResponseWriter) {
		mr.Res(200, struct {
			Msg string `json:"msg"`
		}{
			Msg: "this is working",
		})
	})

	merry.Ship("8080", m)

}
