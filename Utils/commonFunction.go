package utils

import (
	db "cricHeros/Database"
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

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func CreateToken(tokenPayload models.Credential) string {

	expirationTime := time.Now().Add(3 * time.Minute)
	fmt.Println("token Exipiration time is :", expirationTime)
	claims := models.Claims{
		UserId: tokenPayload.User_ID,
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

func CheckMapValidation(data interface{}) error {
	fmt.Println("dadtaaskh: ", data)
	validationErr := Validate.Var(data, "required")
	if validationErr != nil {
		return validationErr
	}
	return nil
}

func DecodeToken(tokenString string, w http.ResponseWriter) (models.Claims, error) {
	claims := &models.Claims{}
	fmt.Println("decode token called")
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		return *claims, fmt.Errorf("invalid or expired token")
	}

	if claims.ExpiresAt.Before(time.Now().Add(2 * time.Minute)) {
		fmt.Println("refresh handler called")
		//generate new token and update to user table
		newClaims := models.Claims{
			UserId: claims.UserId,
			Role:   claims.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
			},
		}
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
		newTokenString, err := newToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			fmt.Println("error is :", err)
		}
		db.DB.Where("user_id=?", claims.UserId).Update("token", newTokenString)
		cookie := http.Cookie{
			Name:     "token",
			Value:    newTokenString,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Minute * 4),
		}
		http.SetCookie(w, &cookie)
	}
	return *claims, nil
}

// func GenrateTokenPair(tokenPayload models.Credential) (map[string]string, error) {
// 	// expirationTime := time.Now().Add(10 * time.Hour)
// 	//creating access token
// 	atPayload := models.AcessToken{
// 		UserId: tokenPayload.User_ID,
// 		Role:   tokenPayload.Role,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
// 		},
// 	}
// 	aToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atPayload)
// 	aTokenString, err := aToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
// 	if err != nil {
// 		fmt.Println("error is :", err)
// 		return nil, err
// 	}

// 	//creating refresh token
// 	rtPayload := models.RefreshToken{
// 		UserId: tokenPayload.User_ID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(
// 				time.Now().Add(15 * time.Hour)),
// 		},
// 	}

// 	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtPayload)
// 	rTokenString, err := rToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
// 	if err != nil {
// 		fmt.Println("error is :", err)
// 		return nil, err
// 	}

// 	return map[string]string{
// 		"accessToken":  aTokenString,
// 		"refreshToken": rTokenString,
// 	}, nil
// }
