package utils

import (
	"encoding/json"
	"net/http"
)

func ShowErr(msg string, statusCode int64, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, msg, int(statusCode))
	// Err := models.Err{
	// 	Message:    msg,
	// 	StatusCode: statusCode,
	// }

	//json.NewEncoder(w).Encode(&Err)
}

func Encode(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(&data)
}
