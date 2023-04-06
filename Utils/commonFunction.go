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
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()
var valid *validator.Validate

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

func IsvalidatePass(password string) (string, bool) {

	if len(password) < 8 {

		return "Password is too short", false
	}
	hasUpperCase := false
	hasLowerCase := false
	hasNumbers := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		} else if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		} else if char >= '0' && char <= '9' {
			hasNumbers = true
		} else if char >= '!' && char <= '/' {
			hasSpecial = true
		} else if char >= ':' && char <= '@' {
			hasSpecial = true
		}
	}

	if !hasUpperCase {
		return "Password do not contain upperCase Character", false

	}

	if !hasLowerCase {
		return "Password do not contain lowerCase Character", false

	}

	if !hasNumbers {
		return "Password do not contain any numbers", false

	}

	if !hasSpecial {
		return "Password do not contain any special character", false

	}
	return "", true
}

func GenerateHashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
