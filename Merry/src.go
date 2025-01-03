package merry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var allRoutes = make(map[string]routeDetails)

type Merry struct {
	BasePattern string
}

type routeDetails struct {
	hanlder func(http.ResponseWriter, *http.Request)
}

func Init(basePattern string, filepath string) (merry *Merry) {

	m := Merry{
		BasePattern: basePattern,
	}
	return &m
}

// Route is for all the route defines after the BasePattern and defiend function
func (m *Merry) Route(pattern string, handler func(mr MerryResponseWriter)) {
	test := MerryResponseWriter{
		context: pattern,
	}

	handler(test)
}

// Ship serves in http
func Ship(port string, m *Merry) (error error) {

	fixedPort := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()

	for rou, details := range allRoutes {
		fullRoutePath := m.BasePattern + rou

		mux.HandleFunc(fullRoutePath, details.hanlder)
	}
	log.Printf("Server is listing on Port:%s", port)
	err := http.ListenAndServe(fixedPort, timeMiddleWare(mux))
	if err != nil {
		return err
	}

	return nil
}

type MerryResponseWriter struct {
	context string
}

type MerryRequest struct {
}

// Res is for the server response
func (mr *MerryResponseWriter) Res(statusCode int, data interface{}) {

	byteData, err := json.Marshal(data)

	if err != nil {
		log.Printf("cannot marshal json:%s", err.Error())
		return
	}

	allRoutes[mr.context] = routeDetails{
		hanlder: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(statusCode)
			w.Write(byteData)
		},
	}

}
