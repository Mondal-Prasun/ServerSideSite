package merry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	UPDATE = "UPDATE"
	DELETE = "DELETE"
)

var allRoutes = make(map[string]routeDetails)

type Merry struct {
	basePattern string
	filepath    string
}

type MerryContext struct {
	context       string
	requestMethod string
	fileDir       string
}

type routeDetails struct {
	hanlder func(http.ResponseWriter, *http.Request)
}

// this initilize the merry package witch return merry struct
func Init(basePattern string, filepath *string) (merry *Merry) {

	if filepath == nil {
		panic("You need a static folder which contents the files and assets")
	}

	m := Merry{
		basePattern: basePattern,
		filepath:    *filepath,
	}
	return &m
}

// Route is for all the route defines after the BasePattern and defiend function
func (m *Merry) Route(requestMethod string, pattern string, handler func(mr MerryContext)) {
	test := MerryContext{
		context:       pattern,
		requestMethod: requestMethod,
		fileDir:       m.filepath,
	}

	handler(test)
}

// Ship serves in http and all the static files
func Ship(port string, m *Merry) (error error) {

	fixedPort := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ext := strings.ToLower(filepath.Ext(r.URL.Path))

		if ext == ".js" {
			w.Header().Set("Content-Type", "application/javascript")
		}

		http.FileServer(http.Dir(m.filepath+"/assets")).ServeHTTP(w, r)
	})))

	for rou, details := range allRoutes {
		fullRoutePath := m.basePattern + rou

		mux.HandleFunc(fullRoutePath, details.hanlder)
	}
	log.Printf("Server is listing on Port:%s", port)
	err := http.ListenAndServe(fixedPort, timeMiddleWare(mux))
	if err != nil {
		return err
	}

	return nil
}

func (mr *MerryContext) ServeHtml(statusCode int, htmlFilePath string) {

	filePath := mr.fileDir + htmlFilePath

	log.Println(filePath)

	htmlFile, err := os.ReadFile(filePath)

	if err != nil {
		mr.Err(501, err.Error())
		return
	}

	allRoutes[mr.context] = routeDetails{
		hanlder: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(statusCode)
			w.Write(htmlFile)
		},
	}

}

// Res is for the server response
func (mr *MerryContext) Res(statusCode int, data interface{}) {

	byteData, err := json.Marshal(data)

	if err != nil {
		log.Printf("cannot marshal json:%s", err.Error())
		return
	}

	allRoutes[mr.context] = routeDetails{
		hanlder: func(w http.ResponseWriter, r *http.Request) {
			if mr.requestMethod == r.Method {

				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(statusCode)
				w.Write(byteData)
			}

		},
	}

}

func (mr *MerryContext) Err(statusCode int, msg string) {
	if statusCode < 300 {
		panic("Status code for error cannot be lower than 300")
	}

	allRoutes[mr.context] = routeDetails{
		hanlder: func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(statusCode)
			w.Write([]byte(msg))
		},
	}

}
