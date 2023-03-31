package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u.SetHeader(w)
	var credential models.Credential
	json.NewDecoder(r.Body).Decode(&credential)

	err := db.DB.Where("user_name=?", credential.Username).First(&models.Credential{}).Error
	if err == nil {
		fmt.Println("User already exist please login to move forward...")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(credential.Password), 8)
	if err != nil {
		panic(err)
	}
	credential.Password = string(bs)
	db.DB.Create(&credential)
	w.Write([]byte("User Registerd sucessfully"))
	json.NewEncoder(w).Encode(credential)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var credential models.Credential
	json.NewDecoder(r.Body).Decode(&credential)
	var existCred models.Credential
	err := db.DB.Where("user_name=?", credential.Username).First(&existCred).Error
	if err != nil {
		fmt.Println("User do not exists please register first...")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(existCred.Password), []byte(credential.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Incorrect Password")
		return
	}
	fmt.Println("Logged In Successfully....")

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	expirationTime := time.Now().Add(5 * time.Minute)

	fmt.Println("expiration time is: ", expirationTime)
	var credential models.Credential
	username := r.URL.Query().Get("username")
	err := db.DB.Where("user_name=?", username).First(&credential).Error
	if err != nil {
		w.Write([]byte("User with given username do not exists....."))
		return
	}
	//check if the user is valid then only create the token
	claims := models.Claims{
		UserId:   credential.User_ID,
		Username: credential.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(os.Getenv("SECRET_KEY"))
	if err != nil {
		fmt.Println("error is :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprint("Token is:", tokenString)))
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	tokenString := r.URL.Query().Get("token")

	claims := &models.Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return os.Getenv("SECRET_KEY"), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	fmt.Println(claims.Username)
	var password = make(map[string]string)

	var userCred models.Credential
	err = db.DB.Where("user_id=?", claims.UserId).Find(&userCred).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	json.NewDecoder(r.Body).Decode(&password)
	fmt.Println("Password: ", password["password"])
	bs, err := bcrypt.GenerateFromPassword([]byte(password["password"]), 8)
	if err != nil {
		panic(err)
	}
	userCred.Password = string(bs)
	userCred.Password = password["password"]
	err = db.DB.Where("user_id=?", claims.UserId).Updates(userCred).Error
	if err != nil {
		http.Error(w, "Failed to update user password", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password updated successfully"))

}
