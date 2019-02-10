package controller

import (
	"testing"
	"net/http/httptest"
	"bytes"
	"encoding/json"
)

func TestCreateUser(t *testing.T) {

	requestURIs := []string{
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
	}

	requestBodies := []map[string]interface{}{
		{"email":"testcase1@gmail.com","password":"test@123"},
		{"email":"testgmail","password":"test@123"},
		{"email":"test@gmail.com","password":"test@1"},
		{"password":"test@123"},
		{"email":"test@gmail.com"},
		{"email":"test@gmail.com","password":"test@123"},

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


	// TODO: Add test cases.
	//success
	//invalid username
	//invalid password
	//missing username
	//missing password
	//existing username--fail

	tests := []testStruct{
		{"Test_1_Success",argsz[0],"Insert successful"},
		{"Test_2_invalid_email",argsz[1],"Email is required. Please provide a valid email id"},
		{"Test_3_wrong_length_pwd",argsz[2],"Password of length of minimum 6 characters is required"},
		{"Test_4_missing_email",argsz[3],"Email is required. Please provide a valid email id"},
		{"Test_5_missing_pwd",argsz[4],"Password of length of minimum 6 characters is required"},
		{"Test_6_in_use_email",argsz[5],"Email already in use"},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateUser(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}


func TestLogin(t *testing.T) {

	requestURIs := []string{
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
		"http://localhost:80/createUser",
	}

	requestBodies := []map[string]interface{}{
		{"email":"testcase@gmail.com","password":"test@123"},
		{"email":"testgmail","password":"test@123"},
		{"email":"test@gmail.com","password":"test@1"},
		{"password":"test@123"},
		{"email":"test@gmail.com"},
		{"email":"testcase@gmail.com","password":"test@123"},

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
		{"Test_1_Success",argsz[0],"Login successful"},
		{"Test_2_invalid_email",argsz[1],"No records found"},
		{"Test_3_wrong_length_pwd",argsz[2],"Invalid login credentials"},
		{"Test_4_missing_email",argsz[3],"No records found"},
		{"Test_5_missing_pwd",argsz[4],"Invalid login credentials"},
		{"Test_6_in_use_email",argsz[5],"Login successful"},
	}



	// TODO: Add test cases.
	//correct username pwd
	//incorrect username
	//inc pwd
	//missing username
	//missing pwd

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Login(tt.argums.w, tt.argums.r)
			erm := checkResponseRecorder(tt)
			if erm != nil{
				t.Errorf("Error is : %v",erm.Error())
			}
		})
	}
}
