package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/golang-jwt/jwt/v4"
)

// @Description Registers a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Credential
// @Param UserDetails body models.Credential true "Registers a user"
// @Tags Authentication
// @Router /adminRegister [post]
func AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u.SetHeader(w)
	u.EnableCors(&w)
	var credential models.Credential
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}

	validationErr := u.CheckValidation(credential)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, validationErr.Error(), w)
		return
	}
	if err, ok := u.IsvalidatePass(credential.Password); !ok {
		u.ShowResponse("Failure", 401, err, w)
		return
	}
	credential.Role = "admin"
	err = db.DB.Where("username=?", credential.Username).First(&models.Credential{}).Error
	if err == nil {
		u.ShowResponse("Failure", 400, "User already exists..", w)
		return
	}
	// bs, err := bcrypt.GenerateFromPassword([]byte(credential.Password), 8)
	// if err != nil {
	// 	panic(err)
	// }
	//credential.Password = string(bs)
	err = db.DB.Create(&credential).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	u.ShowResponse("Success", 200, credential, w)
}

// @Description Registers a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Credential
// @Param UserDetails body models.Credential true "Registers a user"
// @Tags Authentication
// @Router /userRegister [post]
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u.SetHeader(w)
	u.EnableCors(&w)
	var credential models.Credential
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}

	validationErr := u.CheckValidation(credential)
	if validationErr != nil {
		u.ShowResponse("Failure", 400, validationErr.Error(), w)
		return
	}
	if err, ok := u.IsvalidatePass(credential.Password); !ok {
		u.ShowResponse("Failure", 401, err, w)
		return
	}
	credential.Role = "user"
	err = db.DB.Where("username=?", credential.Username).First(&models.Credential{}).Error
	if err == nil {
		u.ShowResponse("Failure", 400, "User already exists..", w)
		return
	}

	// hashPass, err := u.GenerateHashPassword(credential.Password)
	// if err != nil {
	// 	u.ShowResponse("Failure", 400, "Unable to hash the password", w)
	// 	return
	// }
	//credential.Password = string(hashPass)
	err = db.DB.Create(&credential).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	u.ShowResponse("Success", 200, credential, w)

}

// @Description Login a user
// @Accept json
// @Success 200 {string} Logged in successfully
// @Param UserDetails body models.Credential true "Log in the user"
// @Tags Authentication
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	u.EnableCors(&w)
	var credential models.Credential
	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}
	var existCred models.Credential
	err = db.DB.Where("email=?", credential.Email).First(&existCred).Error
	if err != nil {
		u.ShowResponse("Failure", http.StatusUnauthorized, "Incorrect details", w)
		return
	}
	// err = bcrypt.CompareHashAndPassword([]byte(existCred.Password), []byte(credential.Password))
	// if err != nil {
	// 	u.ShowResponse("Failure", http.StatusUnauthorized, "Incorrect details", w)
	// 	return
	// }
	if existCred.Password != credential.Password {
		u.ShowResponse("Failure", http.StatusUnauthorized, "Incorrect details", w)
		return
	}
	tokenString := u.CreateToken(existCred)

	u.ShowResponse("Success", 200, tokenString, w)

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)

	u.EnableCors(&w)
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)

	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}

	var cred models.Credential
	err = db.DB.Where("email=?", mp["email"]).First(&cred).Error
	if err != nil {
		u.ShowResponse("Failure", http.StatusUnauthorized, "User with this email does not exists", w)
		return
	}
	//check if the user is valid then only create the token
	tokenString := u.CreateToken(cred)
	from := "abc@example.com"
	to := []string{
		"prajwal1711@gmail.com",
	}
	url := "http://localhost:8000/resetPassword/" + tokenString
	message := []byte("Click <a href=\"" + url + "\"></a> here to reset your password")

	err = smtp.SendMail("0.0.0.0:1025", nil, from, to, message)
	if err != nil {
		panic(err)
	}
	u.ShowResponse("Success", 200, "Mail sent sucessfully", w)

}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)
	u.EnableCors(&w)
	tokenString := r.Header.Get("token")

	claims := &models.Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("helloWorld"), nil
	})
	if err != nil || !parsedToken.Valid {
		u.ShowResponse("Failure", http.StatusUnauthorized, "Invalid Token", w)
		return
	}

	fmt.Println(claims.UserID)
	var password = make(map[string]string)

	var userCred models.Credential
	err = db.DB.Where("user_id=?", claims.UserID).Find(&userCred).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&password)

	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}
	fmt.Println("Password: ", password["password"])
	// hashPass, err := u.GenerateHashPassword(credential.Password)
	// if err != nil {
	// 	u.ShowResponse("Failure", 400, "Unable to hash the password", w)
	// 	return
	// }
	// userCred.Password = string(hashPass)
	userCred.Password = password["password"]
	err = db.DB.Where("user_id=?", claims.UserID).Updates(userCred).Error
	if err != nil {
		http.Error(w, "Failed to update user password", http.StatusInternalServerError)
		return
	}

	u.ShowResponse("Success", 200, "Password updated successfully", w)

}

// @Description Registers a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Credential
// @Param user_id body object true "ID of the user whose passsword is to be changed"
// @Tags Authentication
// @Router /adminRegister [post]
func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	var mp = make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, "Error in decoding the request body", w)
		return
	}
	var creds models.Credential
	db.DB.Where("user_id=?", mp["userId"]).First(&creds)
	if mp["role"] != "" {
		return
	}

	// err = bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(mp["existPassword"]))
	// if err != nil {
	// 	u.ShowResponse("Failure", 401, "Password not matched", w)
	// 	return
	// }
	if mp["existPassword"] != creds.Password {
		u.ShowResponse("Failure", 401, "Password not matched", w)
		return
	}

	if err, ok := u.IsvalidatePass(mp["newPassword"]); !ok {
		u.ShowResponse("Failure", 401, err, w)
		return
	}
	creds.Password = mp["newPassword"]
	err = db.DB.Where("user_id=?", mp["userId"]).Updates(&creds).Error
	if err != nil {
		u.ShowResponse("Failure", 500, "Error in updating the credentials", w)
		return
	}
	u.ShowResponse("Success", 200, "Password updated successfully", w)

}
