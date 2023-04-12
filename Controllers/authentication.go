package controllers

import (
	db "cricHeros/Database"
	models "cricHeros/Models"
	u "cricHeros/Utils"
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

// send OTP to user
// @Description Sends a OTp to the nuber entered
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param phoneNumber body string true "Phone Number of registered user" SchemaExample({\n"phoneNumber":"string"\n})
// @Tags Authentication
// @Router /sendOTP [post]
func SendOtpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	u.EnableCors(&w)

	var mp = make(map[string]interface{})
	var exists bool

	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	err = validation.Validate(mp,
		validation.Map(
			validation.Key("phoneNumber", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}
	// Check for number

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
	fmt.Println("SID is", sid)
	if ok {
		u.ShowResponse("Success", 200, "OTP sent sucessfully", w)
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

// @Description Verifies the OTP sent to the user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param details body string true "Phone Number of registered user and the otp sent to it" SchemaExample({\n"phoneNumber":"string",\n"otp":"string"\n})
// @Tags Authentication
// @Router /verifyOTP [post]
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	u.EnableCors(&w)

	var mp = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&mp)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}

	err = validation.Validate(mp,
		validation.Map(
			validation.Key("phoneNumber", validation.Required),
			validation.Key("otp", validation.Required),
		),
	)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	if CheckOtp("+91"+mp["phoneNumber"].(string), mp["otp"].(string), w) {

		var userDetails models.Credential
		db.DB.Where("phone_number=?", mp["phoneNumber"]).First(&userDetails)
		expirationTime := time.Now().Add(time.Hour * 4)
		fmt.Println("Cookie expiration time: ", expirationTime)
		userDetails.IsLoggedIn = true
		tokenString := u.CreateToken(userDetails)
		cookie := &http.Cookie{
			Name:  "token",
			Value: tokenString,
			// Secure: true,
			HttpOnly: true,
			Expires:  expirationTime,
		}
		http.SetCookie(w, cookie)
		userDetails.Token = tokenString
		db.DB.Where("user_id=?", userDetails.User_ID).Updates(userDetails)
		u.ShowResponse("Success", 200, tokenString, w)

		return
	} else {
		// fmt.Println("Verifictaion failed")
		u.ShowResponse("Failure", 401, "Verifictaion Failed", w)
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
	var existRecord models.Credential

	err := json.NewDecoder(r.Body).Decode(&Credential)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	Credential.Role = "admin"

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
// @Tags Authentication
// @Router /logOut [get]
func LogOut(w http.ResponseWriter, r *http.Request) {
	tokenString, err := r.Cookie("token")

	if err != nil {
		u.ShowResponse("Failure", 403, err.Error(), w)
		return
	}
	var blackList models.Blacklist
	//Decode the token
	payload, err := u.DecodeToken(tokenString.Value, w)
	if err != nil {
		u.ShowResponse("Failure", 400, err.Error(), w)
		return
	}

	db.DB.Model(&models.Credential{}).Where("user_id=?", payload.UserId).Update("is_logged_in", false)
	db.DB.Model(&models.Credential{}).Where("user_id=?", payload.UserId).Update("token", nil)

	//adding the token to blacklist table
	blackList.Token = tokenString.Value
	db.DB.Create(blackList)

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	u.ShowResponse("Success", 200, "Logged out successfully", w)
}

// @Description Updates the data of the user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Param userDetails body models.Credential true "user updated datas"
// @Tags Authentication
// @Router /updateProfile [post]
func UpdateProfile(w http.ResponseWriter, r *http.Request) {

	var creds models.Credential
	json.NewDecoder(r.Body).Decode(&creds)

	tokenString, err := r.Cookie("token")
	if err != nil {
		u.ShowResponse("Failure", 403, err.Error(), w)
		return
	}

	//Decode the token
	payload, err := u.DecodeToken(tokenString.Value, w)
	if err != nil {
		u.ShowResponse("Failure", 400, err, w)
		return
	}
	if creds.Role != "" || creds.Role != payload.Role {
		u.ShowResponse("Failure", 403, "Forbidden", w)
		return
	}

	db.DB.Where("u_id=?", payload.UserId).Updates(&creds)
	u.ShowResponse("Success", 200, "User updated successfully", w)

}
