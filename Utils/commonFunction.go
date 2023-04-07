package utils

import (
	models "cricHeros/Models"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

var Validate = validator.New()

func ShowResponse(status string, statusCode int64, data interface{}, w http.ResponseWriter) {
	SetHeader(w)
	w.WriteHeader(int(statusCode))
	response := models.Response{
		Status: status,
		Code:   statusCode,
		Data:   data,
	}

	json.NewEncoder(w).Encode(&response)
}
func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Encode(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(&data)
}

func SetHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func CreateToken(tokenPayload models.Credential) string {
	fmt.Println("payload: ", tokenPayload)

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := models.Claims{
		UserID: tokenPayload.User_ID,
		Role:   tokenPayload.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		fmt.Println("error is :", err)
		return ""
	}

	return tokenString
}

func CheckValidation(data interface{}) error {
	validationErr := Validate.Struct(data)
	if validationErr != nil {
		return validationErr
	}
	return nil
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func CheckMapValidation(data interface{}) error {
	fmt.Println("dadtaaskh: ", data)
	validationErr := Validate.Var(data, "required")
	if validationErr != nil {
		return validationErr
	}
	return nil
}
