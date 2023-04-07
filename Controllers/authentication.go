package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"

	"net/http"
	"os"

	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var twilioClient *twilio.RestClient

func TwilioInit(password string) {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: u.TWILIO_ACCOUNT_SID,
		Password: password,
	})

}

// // twilio client interface
// var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
// 	Username: u.TWILIO_ACCOUNT_SID,
// 	Password: u.TWILIO_AUTH_TOKEN,
// })

// send OTP to user
func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	u.EnableCors(&w)

	var mp = make(map[string]interface{})

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	// Check for number
	var exists bool
	err = db.DB.Raw("SELECT EXISTS(SELECT 1 FROM credentials WHERE phone_number=?)", mp["phoneNumber"]).Scan(&exists).Error
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	// Response
	if !exists {
		u.ShowResponse("Failure", 409, "Number do not exists, please register first", w)
		return
	}
	ok, sid := sendOtp("+91"+mp["phoneNumber"].(string), w)
	if ok {
		u.ShowResponse("Success", 200, sid, w)
	}

}

// function to send OTP while user registration
func sendOtp(to string, w http.ResponseWriter) (bool, *string) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)

	params.SetChannel("sms")

	resp, err := twilioClient.VerifyV2.CreateVerification(os.Getenv("VERIFY_SERVICE_SID"), params)
	if err != nil {
		u.ShowResponse("Failure", 401, "No credentials provided", w)
		return false, nil
	} else {
		return true, resp.Sid
	}

}

// Check OTP status
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	u.EnableCors(&w)

	var mp = make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&mp)

	if CheckOtp("+91"+mp["phoneNumber"].(string), mp["otp"].(string), w) {

		var userDetails models.Credential
		db.DB.Where("phone_number=?", mp["phoneNumber"]).First(&userDetails)
		userDetails.IsLoggedIn = true
		tokenString := u.CreateToken(userDetails)
		db.DB.Where("user_id=?", userDetails.User_ID).Updates(userDetails)
		u.ShowResponse("Success", 200, tokenString, w)

		return
	} else {
		// fmt.Println("Verifictaion failed")
		u.ShowResponse("Not Found", 404, "OTP ERROR", w)
		return
	}
}

// OTP code verification
func CheckOtp(to string, code string, w http.ResponseWriter) bool {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(os.Getenv("VERIFY_SERVICE_SID"), params)

	if err != nil {
		return false
	} else if *resp.Status == "approved" {
		return true
	} else {
		return false
	}
}

// @Description Registers a admin
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param UserDetails body models.Credential true "Registers a admin"
// @Tags Authentication
// @Router /adminRegister [post]
func AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)

	var Credential models.Credential
	err := json.NewDecoder(r.Body).Decode(&Credential)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	Credential.Role = "admin"

	var existRecord models.Credential
	err = db.DB.Where("phone_number=?", Credential.PhoneNumber).First(&existRecord).Error
	if err == nil {
		u.ShowResponse("Failure", 400, "User already register please login to contnue", w)
		return
	}
	err = db.DB.Create(&Credential).Error
	if err != nil {

		u.ShowResponse("Failure", 500, "Internal Server Error", w)
		return
	}
	u.ShowResponse("Success", 200, Credential, w)
}

// @Description Registers a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param UserDetails body models.Credential true "Registers a user"
// @Tags Authentication
// @Router /userRegister [post]
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	u.SetHeader(w)

	var Credential models.Credential
	err := json.NewDecoder(r.Body).Decode(&Credential)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	Credential.Role = "user"
	err = db.DB.Create(&Credential).Error
	if err != nil {
		u.ShowResponse("Failure", 500, "Internal Server Error", w)
		return
	}
	u.ShowResponse("Success", 200, Credential, w)

}

// @Description Logs out a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param token header string true "token generated for the user"
// @Tags Authentication
// @Router /logOut [get]
func LogOut(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("token")

	//Decode the token
	claims, err := u.DecodeToken(tokenString)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now())
	db.DB.Model(&models.Credential{}).Where("user_id=?", claims.UserID).Update("is_logged_in", false)
	u.ShowResponse("Success", 200, "Logged out successfully", w)

}

// @Description Updates the data of the user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param token header string true "token generated for the user"
// @Param userDetails body models.Credential true "user updated datas"
// @Tags Authentication
// @Router /updateProfile [post]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	var creds models.Credential
	json.NewDecoder(r.Body).Decode(&creds)

	tokenString := r.Header.Get("token")

	//Decode the token
	claims, err := u.DecodeToken(tokenString)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	if creds.Role != "" || creds.Role != claims.Role {
		u.ShowResponse("Failure", 403, "Forbidden", w)
		return
	}

	db.DB.Where("u_id=?", claims.UserID).Updates(&creds)
	u.ShowResponse("Success", 200, "User updated successfully", w)

}
