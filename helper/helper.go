package helper

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
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
		return errors.Wrap(err, fmt.Sprintf("Error encoding json response: %v", obj))
	}
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}
