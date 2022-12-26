package helpers

import (
	"moneytracker/pkg/logging"
	"net/http"
)

var log = logging.GetLogger()

func NewErrorResponse(w http.ResponseWriter, statusCode int, massage string) {
	log.Errorln(massage)
	http.Error(w, massage, statusCode)
}

func NewSuccessResponse(w http.ResponseWriter, statusCode int, message string) {
	log.Println(message)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Errorln(err)
		return
	}
	w.WriteHeader(statusCode)
}
