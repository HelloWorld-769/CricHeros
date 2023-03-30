package utils

import (
	"encoding/json"
	"math"
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
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Encode(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(&data)
}
