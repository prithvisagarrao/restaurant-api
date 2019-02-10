package models

import (
	"github.com/dgrijalva/jwt-go"
	"strings"
	"fmt"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/antigloss/go/logger"
)

type Recipe struct {
	RecipeId string	`json:"recipe_id"`
	Name string		`json:"name"`
	PrepTime string	`json:"prep_time"`
	Difficulty	int	`json:"difficulty"`
	Vegetarian bool	`json:"vegetarian"`
}

type ReqAbs struct {
	HTTPRequestType string
	Payload map[string]interface{}
	Filter	string
}

type DBAbs struct {
	Query string
	Status string
	StatusCode int
	Data string
	DataHTTP	interface{}
	RichData []map[string]interface{}
	Message	string
	Count	uint64
}

type RespAbs struct {
	Status	string	`json:"status"`
	StatusCode int	`json:"status_code"`
	Message	string	`json:"message"`
	ResultData	interface{} `json:"result_data"`
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	UserId		uint	`json:"user_id"`
	Email		string	`json:"email"`
	Password	string	`json:"password"`
	TokenString	string
}

//Validate method checks and validates the user credentials provided
func (user *User) Validate() (ret string, ok bool) {

	if !strings.Contains(user.Email,"@"){
		logger.Error("Email is required. Please provide a valid email id")
		ret = "Email is required. Please provide a valid email id"
		ok = false
		return
	}

	if len(user.Password) <=6 {
		logger.Error("Password of length of minimum 6 characters is required")
		ret = "Password of length of minimum 6 characters is required"
		ok = false
		return
	}



	var dbAbs DBAbs

	selectQuery := fmt.Sprintf("SELECT email FROM user_info WHERE email = '%v' ;",user.Email)

	dbAbs.Query = selectQuery
	Select(&dbAbs)

	if dbAbs.Status != "success"{

		logger.Error("Could not confirm if user email already exists.")
		ret = "Could not confirm if user email already exists."
		ok = false
		return
	}

	if len(dbAbs.RichData) > 0 {
		for k, v := range dbAbs.RichData[0] {
			if k == "email" {
				if v != "" {
					ret = "Email already in use"
					ok = false
					return
				}
			}
		}
	}

	ret = "Requirements passed"
	ok = true

	return
}

//CreateAccount creates a new user account
func (user User)CreateAccount() (dbAbs DBAbs){

	ret, ok := user.Validate()
	if !ok {

		logger.Error("User credentials failed")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = ret

		return
	}


	if len(user.Email) == 0 || len(user.Password) == 0{
		dbAbs.StatusCode=http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "Invalid request"
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	insertQuery := fmt.Sprintf("INSERT INTO user_info (email, password) VALUES ('%v', '%v');",user.Email,user.Password)

	logger.Info("Running query to create user")
	dbAbs.Query = insertQuery
	Insert(&dbAbs)
	return
}

//Login method validates the user from the values in the database
func (user *User)Login (email string, password string) (dbAbs DBAbs){

	selectStr := fmt.Sprintf("SELECT * FROM user_info WHERE email = '%v' ;", email)
	dbAbs.Query = selectStr
	Select(&dbAbs)

	if dbAbs.Status != "success"{
		logger.Error("Could not check if user exists")
		dbAbs.StatusCode = http.StatusInternalServerError
		return
	}

	var pwd string
	var userId	uint

	if len(dbAbs.RichData) >0 {
		for k, v := range dbAbs.RichData[0] {
			if k == "email" && v == "" {
				dbAbs.Message = "Email not found"
				dbAbs.Status = "fail"
				return
			}

			if k == "password"{
				if _, ok := v.(string); ok {
					pwd = v.(string)
				}
			}

			if k == "user_id"{
				if _, ok := v.(int64); ok{
					userId = uint(v.(int64))
				}
			}

		}

		err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
		if err != nil {
			logger.Error("Invalid login credentials.Please check your password")
			dbAbs.Status = "fail"
			dbAbs.StatusCode = http.StatusBadRequest
			dbAbs.Message = "Invalid login credentials"
			return

		}

		user.Password = ""
		user.UserId = userId

		dbAbs.Message = "Login successful"
		logger.Info("Login successful")
	}

	return
}