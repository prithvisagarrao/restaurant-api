package utils

import (
	"github.com/prithvisagarrao/restaurant-api/models"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"errors"
	"strconv"
	"github.com/antigloss/go/logger"
)

//PerpareResponseFromDB prepares the response string from the DBAbs
func PrepareResponseFromDB(dbabs models.DBAbs) ( string){

	var v interface{}

	err := json.Unmarshal([]byte(dbabs.Data), &v)
	if err != nil{
		logger.Error(fmt.Sprintf("Error in unmarshalling dbAbs data, error: %v",err))
	}

	respAbs := models.RespAbs{
		Status:dbabs.Status,
		StatusCode:dbabs.StatusCode,
		Message:dbabs.Message,
		ResultData:v,
	}

	jsonString ,err := json.Marshal(respAbs)
	if err != nil {
		logger.Error(fmt.Sprintf("Error in encoding JSON : %v", err.Error()))
	}

	response := string(jsonString)
	return response
}

//Respond writes the response to the http.ResponseWritter
func Respond(w http.ResponseWriter, response string){

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintln(w, response)

}

//MapReqAbs generates a request abstract by mapping the http request accordingly
func MapReqAbs(r *http.Request) models.ReqAbs {
	var reqAbs models.ReqAbs
	reqAbs.HTTPRequestType = r.Method

	if reqAbs.HTTPRequestType == "POST" {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {

			logger.Error(fmt.Sprintf("Error in reading POST Body : %v ", err.Error()))
		}

		err = json.Unmarshal(body, &reqAbs.Payload)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in parsing JSON body in POST -- Error : %v", err))
		}
	} else if reqAbs.HTTPRequestType == "GET" || reqAbs.HTTPRequestType == "DELETE"	{

		reqAbs.Filter = r.RequestURI

	} else if reqAbs.HTTPRequestType == "PUT" || reqAbs.HTTPRequestType == "PATCH" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {

			logger.Error(fmt.Sprintf("Error in reading POST Body : %v ", err.Error()))
		}

		err = json.Unmarshal(body, &reqAbs.Payload)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in parsing JSON body in POST -- Error : %v", err))
		}

		reqAbs.Filter = r.RequestURI
	}

	return  reqAbs
}

//OffsetLimitStr returns the string with the offset and limit specified
func OffsetLimitStr(filterstr string) (ret string, err error){
	filter := strings.Split(filterstr,"?")

	var offset, limit string

	if len(filter) > 1{
		strSplit := strings.Split(filter[1],"&")

		if len(strSplit) > 1{
			for _, v := range strSplit{

				if strings.Contains(strings.ToLower(v),"offset"){
					sOffset := strings.Split(v,"=")
					if len(sOffset) > 1 {
						offset = sOffset[1]
					} else {
						logger.Error("offset provided in wrong format")
						err = errors.New("offset provided in wrong format")
						return "",err
					}
				}
				if strings.Contains(strings.ToLower(v),"limit"){
					sLimit := strings.Split(v, "=")
					if len(sLimit) > 1{
						limit = sLimit[1]
					} else {
						logger.Error("limit provided in wrong format")
						err = errors.New("limit provided in wrong format")
						return "",err
					}
				}
			}
		} else {
			if strings.Contains(strings.ToLower(strSplit[0]), "offset") {
				sOffset := strings.Split(strSplit[0], "=")
				if len(sOffset) > 1 {
					offset = sOffset[1]
				} else {
					logger.Error("offset provided in wrong format")
					err = errors.New("offset provided in wrong format")
					return "",err
				}
			}
			if strings.Contains(strings.ToLower(strSplit[0]), "limit") {
				sLimit := strings.Split(strSplit[0], "=")
				if len(sLimit) > 1 {
					limit = sLimit[1]
				} else {
					logger.Error("limit provided in wrong format")
					err = errors.New("limit provided in wrong format")
					return "",err
				}
			}
		}
	}

	if limit == ""{
		limit = "20"
	}


	if of,er := strconv.Atoi(offset); er != nil{
		if li, erl := strconv.Atoi(limit); erl != nil{
			logger.Error("using default values for offset and limit")
			ret = fmt.Sprintf(" LIMIT 20")
		} else {
			ret = fmt.Sprintf(" LIMIT %v", li)
		}
	} else {
		if li, erl := strconv.Atoi(limit); erl != nil{
			logger.Error("using default values for limit")
			ret = fmt.Sprintf(" OFFSET %v LIMIT 20",of)
		} else {
			ret = fmt.Sprintf(" OFFSET %v LIMIT %v",of, li)
		}
	}

return

}