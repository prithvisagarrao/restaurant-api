package controller

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"encoding/json"
	"github.com/gorilla/mux"
	"strings"
	"errors"
	"strconv"
	"bytes"
	"github.com/prithvisagarrao/restaurant-api/config"
)

func TestGetRecipe(t *testing.T) {
	config.InitConfig()
	requestURIs := []string{
		"http://localhost:80/recipes/16",
		"http://localhost:80/recipes/g",
		"http://localhost:80/recipes/19",
	}

	var argsz []args

	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		testData := make(map[string]string)
		tempR := httptest.NewRequest("GET",requestURIs[i],http.NoBody)
		testData["id"] = strings.Split(requestURIs[i],"recipes/")[1]
		testArgvTemp.r = mux.SetURLVars(tempR,testData)
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}

	tests := []testStruct{
		{"Test_1_success",argsz[0],"Select successful"},
		{"Test_2_failed_wrong_id",argsz[1],"recipe_id invalid"},
		{"Test_3_NON_Existent_id",argsz[2],"No records found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

type args struct {
	w *httptest.ResponseRecorder
	r *http.Request
}

type testStruct struct{
	name string
	argums args
	expect string
}


func getMessageFromBody(bod *bytes.Buffer) (ret string,err error){
	var jsonBody map[string]interface{}
	err = json.Unmarshal(bod.Bytes(),&jsonBody)
	if err != nil{
		return "",err
	}

	if msg,ok := jsonBody["message"]; !ok{
		return "",errors.New("No message found in body")
	} else {
		if msgStr,ok1 := msg.(string); !ok1{
			return "",errors.New("Message found but was not string")
		} else {
			return msgStr,nil
		}
	}
}

func TestDeleteRecipe(t *testing.T) {
	config.InitConfig()

	requestURIs := []string{
		"http://localhost:80/recipes/28",
		"http://localhost:80/recipes/g",
		"http://localhost:80/recipes/102450",
	}

	var argsz []args

	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		testData := make(map[string]string)
		tempR := httptest.NewRequest("DELETE",requestURIs[i],http.NoBody)
		testData["id"] = strings.Split(requestURIs[i],"recipes/")[1]
		testArgvTemp.r = mux.SetURLVars(tempR,testData)
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}


	tests := []testStruct{
		{"Test_1_success_request",argsz[0],"Delete successful"},
		{"Test_2_fail_request",argsz[1],"recipe_id invalid"},
		{"Test_3_no_records_deleted",argsz[2],"No records found to delete"},
	}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

func checkResponseRecorder(testInst testStruct) (err error){

	//1. Status Code check
	if testInst.argums.w.Code != 200 {
		return errors.New("Status code not matched : Expected : 200 ; Got : "+strconv.Itoa(testInst.argums.w.Code))
	}
	//2.  Body Length > 0
	if testInst.argums.w.Body.Len() <= 0{
		return errors.New("Body not found : Length of Body : "+strconv.Itoa(testInst.argums.w.Body.Len()))
	}
	//3. ExpectVals must match the RespAbs.Message
	gotBody,err := getMessageFromBody(testInst.argums.w.Body)
	if err != nil{
		return errors.New("Check Failed because no message could be found from body "+err.Error())
	} else {
		if !strings.EqualFold(testInst.expect,gotBody) {
			return errors.New("Expected Message : "+testInst.expect+" did not match the message we got as "+gotBody)
		}
	}

	return
}




func TestListAllRecipes(t *testing.T) {
	config.InitConfig()
	requestURIs := []string{
		"http://localhost:80/recipes?offset=3&limit=5",
		"http://localhost:80/recipes?limit=5",
		"http://localhost:80/recipes?offset=3",
		"http://localhost:80/recipes?offset=a&limit=5",
		"http://localhost:80/recipes?offset=3&limit=b",
		"http://localhost:80/recipes?offset=3&limit=7.7",
		"http://localhost:80/recipes",
	}


	var argsz []args

	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		tempR := httptest.NewRequest("GET",requestURIs[i],http.NoBody)
		testArgvTemp.r = tempR
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}

	tests := []testStruct{
		{"Success Case with offset and limit",argsz[0],"Select successful"},
		{"Success Case with only limit",argsz[1],"Select successful"},
		{"Success Case with only offset",argsz[2],"Select successful"},
		{"URI with incorrect Datatype for offset",argsz[3],"Select successful"},
		{"URI with incorrect Datatype for limit",argsz[4],"Select successful"},
		{"URI with incorrect datatype",argsz[5],"Select successful"},
		{"URI without ?",argsz[6],"Select successful"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ListAllRecipes(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}


}



func TestUpdateRecipe(t *testing.T) {
	config.InitConfig()


	requestURIs := []string{
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/13",
		"http://localhost:80/recipes/g",
		"http://localhost:80/recipes/10250150",
	}

	requestBodies := []map[string]interface{}{
		{"name":"type911","prep_time":"6 hours","difficulty":2,"vegetarian":false},
		{"name":"type912","prep_time":"6 hours","difficulty":"a","vegetarian":true},
		{"name":"type913","prep_time":"6 hours","difficulty":3,"vegetarian":"asfaf"},
		{"name":"type914","prep_time":"6 hours","difficulty":3},
		{"name":"type915","prep_time":"6 hours","vegetarian":true},
		{"name":"type916","difficulty":3,"vegetarian":true},
		{"name":"type917","prep_time":"6 hours","difficulty":2,"vegetarian":false},
		{"name":"type918","prep_time":"6 hours","difficulty":2,"vegetarian":false},
	}

	var argsz []args



	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		testData := make(map[string]string)
		tempReader,_ := json.Marshal(requestBodies[i])
		tempR := httptest.NewRequest("PUT",requestURIs[i],bytes.NewReader(tempReader))
		testData["id"] = strings.Split(requestURIs[i],"recipes/")[1]
		testArgvTemp.r = mux.SetURLVars(tempR,testData)
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}


	tests := []testStruct{
		{"Test_1_Success",argsz[0],"update successful"},
		{"Test_2_difficulty_data_format_mistake",argsz[1],"data type invalid for difficulty"},
		{"Test_3_veg_data_format_mistake",argsz[2],"data invalid for vegetarian ; must be true or false"},
		{"Test_4_veg_missing",argsz[3],"update successful"},
		{"Test_5_difficulty_missing",argsz[4],"update successful"},
		{"Test_6_prep_time_missing",argsz[5],"update successful"},
		{"Test_7_invalid_id",argsz[6],"recipe_id invalid"},
		{"Test_8_missing_id",argsz[7],"no recipes found to update"},
	}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

func TestRateRecipe(t *testing.T) {
	config.InitConfig()
	requestURIs := []string{
		"http://localhost:80/recipes/13/rating",
		"http://localhost:80/recipes/a/rating",
		"http://localhost:80/recipes/17/rating",
		"http://localhost:80/recipes/13555",
		"http://localhost:80/recipes/17/rating",
	}

	requestBodies := []map[string]interface{}{
		{"rate":3},
		{"rate":5},
		{"rate":"c"},
		{"rate":3},
		{"rate":7},
	}

	var argsz []args


	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		testData := make(map[string]string)
		tempReader,_ := json.Marshal(requestBodies[i])
		tempR := httptest.NewRequest("POST",requestURIs[i],bytes.NewReader(tempReader))
		strSplit := strings.Split(requestURIs[i],"recipes/")[1]
		testData["id"] = strings.Split(strSplit,"/")[0]
		testArgvTemp.r = mux.SetURLVars(tempR,testData)
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}

	tests := []testStruct{
		{"Test_1_Success",argsz[0],"Insert successful"},
		{"Test_2_id_mistake",argsz[1],"recipe_id invalid"},
		{"Test_3_invalid_rate",argsz[2],"data type invalid for rate"},
		{"Test_4_wrong_id",argsz[3],"Insert failed"},
		{"Test_5_difficulty_missing",argsz[4],"data invalid"},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RateRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

func TestCreateRecipe(t *testing.T) {
	config.InitConfig()
	requestURIs := []string{
		"http://localhost:80/recipes",
		"http://localhost:80/recipes",
		"http://localhost:80/recipes",
		"http://localhost:80/recipes",
		"http://localhost:80/recipes",
		"http://localhost:80/recipes",
	}

	requestBodies := []map[string]interface{}{
		{"name":"type911","prep_time":"6 hours","difficulty":2,"vegetarian":false},
		{"name":"type912","prep_time":"6 hours","difficulty":"a","vegetarian":true},
		{"name":"type913","prep_time":"6 hours","difficulty":3,"vegetarian":"asfaf"},
		{"name":"type914","prep_time":"6 hours","difficulty":3},
		{"name":"type915","prep_time":"6 hours","vegetarian":true},
		{"name":"type916","difficulty":3,"vegetarian":true},
	}

	var argsz []args



	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		tempReader,_ := json.Marshal(requestBodies[i])
		tempR := httptest.NewRequest("POST",requestURIs[i],bytes.NewReader(tempReader))
		testArgvTemp.r = tempR
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}


	tests := []testStruct{
		{"Test_1_Success",argsz[0],"Insert successful"},
		{"Test_2_difficulty_data_format_mistake",argsz[1],"data type invalid for difficulty"},
		{"Test_3_veg_data_format_mistake",argsz[2],"data invalid for vegetarian ; must be true or false"},
		{"Test_4_veg_missing",argsz[3],"Insert successful"},
		{"Test_5_difficulty_missing",argsz[4],"Insert successful"},
		{"Test_6_prep_time_missing",argsz[5],"Insert successful"},
	}



	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

func TestSearchRecipe(t *testing.T) {
	config.InitConfig()
	requestURIs := []string{
		"http://localhost:80/recipes/test9",
		"http://localhost:80/recipes/notPresent",
	}

	var argsz []args

	for i:=0;i<len(requestURIs);i++{
		var testArgvTemp args
		testData := make(map[string]string)
		tempR := httptest.NewRequest("GET",requestURIs[i],http.NoBody)
		testData["name"] = strings.Split(requestURIs[i],"recipes/search/")[1]
		testArgvTemp.r = mux.SetURLVars(tempR,testData)
		testArgvTemp.w = httptest.NewRecorder()
		argsz = append(argsz,testArgvTemp)
	}

	tests := []testStruct{
		{"Test_1_success",argsz[0],"Select successful"},
		{"Test_2_invalid_name",argsz[1],"No records found"},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SearchRecipe(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}

