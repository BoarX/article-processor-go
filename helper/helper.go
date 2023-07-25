package helper

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	HeaderContentType   = "Content-Type"
	MimeApplicationJSON = "application/json"
)

func SendJson(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("error encoding json response: %v", err)
		return fmt.Errorf("error encoding json response: %v", err)
	}
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}

func SendJsonOk(w http.ResponseWriter, obj interface{}) error {
	return SendJson(w, http.StatusOK, obj)
}

func SendJsonError(w http.ResponseWriter, status int, obj interface{}) error {
	w.Header().Set(HeaderContentType, MimeApplicationJSON)
	jsonData, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Send Json Error %v", err)
		return fmt.Errorf("error encoding json response: %v", err)
	}
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}
