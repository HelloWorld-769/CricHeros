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

	"github.com/o1egl/paseto"

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

func CheckValidation(data interface{}) error {
	validationErr := Validate.Struct(data)
	if validationErr != nil {
		return validationErr
	}
	return nil
}

func CreateToken(tokenPayload models.Credential) string {

	expirationTime := time.Now().Add(2 * time.Hour)
	fmt.Println("JWT token Exipiration time is :", expirationTime)
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

func DecodeToken(tokenString string, w http.ResponseWriter) (*models.Claims, error) {
	claims := &models.Claims{}
	fmt.Println("decode token called")
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	//creating a new token when the token expiry time is only 2 min left
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
			return nil, err
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
	return claims, nil
}

func CreatePasetoToken(tokenPayload models.Credential) (string, error) {

	payload := &models.Payload{
		UserId:    tokenPayload.User_ID,
		Role:      tokenPayload.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Minute * 5),
	}

	tokenString, err := paseto.NewV2().Encrypt([]byte(os.Getenv("PASETO_KEY")), payload, nil)
	if err != nil {
		return "", nil
	}

	return tokenString, nil

}

func DecodePasetoToken(token string, w http.ResponseWriter) (*models.Payload, error) {
	payload := &models.Payload{}
	err := paseto.NewV2().Decrypt(token, []byte(os.Getenv("PASETO_KEY")), payload, nil)
	if err != nil {
		return nil, err
	}

	if payload.ExpiredAt.Before(time.Now().Add(2 * time.Minute)) {
		//create a new paseto token for them
		newPayload := &models.Payload{
			UserId:    payload.UserId,
			Role:      payload.Role,
			IssuedAt:  time.Now(),
			ExpiredAt: time.Now().Add(time.Minute * 5),
		}

		newTokenString, err := paseto.NewV2().Encrypt([]byte(os.Getenv("PASETO_KEY")), newPayload, nil)
		if err != nil {
			return nil, err
		}
		db.DB.Where("user_id=?", payload.UserId).Update("token", newTokenString)
		cookie := http.Cookie{
			Name:     "token",
			Value:    newTokenString,
			HttpOnly: true,
			Expires:  time.Now().Add(time.Minute * 4),
		}
		http.SetCookie(w, &cookie)

	}

	return payload, nil
}
