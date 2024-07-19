package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type statusResponse struct {
	Status string `json: status`
}

func servErr(w http.ResponseWriter, err error, message string) {
	//trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//log.Println(trace)
	log.Println(message)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func clientErr(w http.ResponseWriter, status int, message string) {
	log.Println(message)
	http.Error(w, message, status)
}

func notFound(w http.ResponseWriter) {
	clientErr(w, http.StatusNotFound, "page not found")
}

func JSONStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
