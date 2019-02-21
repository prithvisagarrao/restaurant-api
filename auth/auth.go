package auth

import (
	"net/http"
	"strings"
	"github.com/prithvisagarrao/restaurant-api/models"
	"github.com/prithvisagarrao/restaurant-api/utils"
	"github.com/prithvisagarrao/restaurant-api/config"
	"github.com/dgrijalva/jwt-go"
	"context"
	"crypto/rsa"
	"io/ioutil"
	"fmt"
	"time"
	"github.com/antigloss/go/logger"
)

var (
	verifyKey *rsa.PublicKey
	signKey *rsa.PrivateKey
)

//init is to initialize the private and public keys for token authentication
func init(){

	var Conf = config.GetConfig()

	signBytes, err := ioutil.ReadFile(Conf.Keys.PrivateKeyPath)
	if err != nil{
		logger.Error(fmt.Sprintf("Error in reading key file %v",err))
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil{
		logger.Error(fmt.Sprintf("Error in parsing private key file %v",err))
	}

	verifyBytes, err := ioutil.ReadFile(Conf.Keys.PublicKeyPath)
	if err != nil{
		logger.Error(fmt.Sprintf("Error in reading key file %v",err))
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil{
		logger.Error(fmt.Sprintf("Error in parsing public key file %v",err))
	}

	
}

//JWTAuthentication is the middleware layer that authenticates the protected endpoints
var JWTAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notProtectedUserRating :=  []string{"/createUser", "/login"}
		requestPath := r.URL.Path

		for _, value := range notProtectedUserRating {
			if value == requestPath{
				next.ServeHTTP(w, r)
				return
			}
		}

		if strings.Contains(requestPath,"rating"){
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(requestPath,"recipes") && r.Method == "GET"{
			next.ServeHTTP(w, r)
			return
		}


		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == ""{
			logger.Error("Missing authorization token")
			var dbAbs models.DBAbs
			dbAbs.StatusCode = http.StatusForbidden
			dbAbs.Status = "fail"
			dbAbs.Message = "Missing auth token"
			respAbs := utils.PrepareResponseFromDB(dbAbs)
			utils.Respond(w,respAbs)
			return
		}

		strSplit := strings.Split(tokenHeader, " ")
		if len(strSplit) != 2{
			logger.Error("Invalid authorization token")
			var dbAbs models.DBAbs
			dbAbs.StatusCode = http.StatusForbidden
			dbAbs.Status = "fail"
			dbAbs.Message = "Invalid auth token"
			respAbs := utils.PrepareResponseFromDB(dbAbs)
			utils.Respond(w, respAbs)
			return
		}

		var tokenPart string
		if len(strSplit) > 0{
			tokenPart = strSplit[1]
		}

		tk := &models.Token{}
		token, err:= jwt.ParseWithClaims(tokenPart,tk, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

		if err != nil || !token.Valid{
			logger.Error("Invalid authorization token")
			var dbAbs models.DBAbs
			dbAbs.StatusCode = http.StatusForbidden
			dbAbs.Status = "fail"
			dbAbs.Message = "Invalid auth token"
			respAbs := utils.PrepareResponseFromDB(dbAbs)
			utils.Respond(w, respAbs)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

//GenerateToken generates a JWT for a particular user_id
func GenerateToken(userID uint) (tokenString string){

	tk := &models.Token{UserId: userID}
	tk.ExpiresAt = time.Now().Add(time.Minute * 4).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), tk)
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in generating token %v",err.Error()))
		tokenString = ""
		return
	}

	return
}