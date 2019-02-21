package controller

import (
	"net/http"
	"github.com/prithvisagarrao/restaurant-api/models"
	"encoding/json"
	"github.com/prithvisagarrao/restaurant-api/utils"
	"github.com/prithvisagarrao/restaurant-api/auth"
	"github.com/antigloss/go/logger"
)

//CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		logger.Error("Invalid data sent")
		var dbAbs models.DBAbs
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "Invalid Request"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}

	dbAbs := account.CreateAccount()
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
	return

}

//Login is used to allow valid users to login into the application
func Login(w http.ResponseWriter, r *http.Request){

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		logger.Error("Invalid data sent")
		var dbAbs models.DBAbs
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "Invalid Request"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}

	dbAbs := account.Login(account.Email, account.Password)
	if dbAbs.Status != "success"{
		logger.Error("Login unsuccessful")
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}
	tokenString := auth.GenerateToken(account.UserId)
	account.TokenString = tokenString
	data, err := json.Marshal(account)
	dbAbs.Status = "success"
	dbAbs.StatusCode = http.StatusOK
	dbAbs.Data = string(data)
	logger.Info("Token generated")
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
	return
}