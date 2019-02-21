package utils

import (
	"testing"

	"github.com/prithvisagarrao/restaurant-api/models"
)

func TestPrepareResponseFromDB(t *testing.T) {
	type args struct {
		dbabs models.DBAbs
	}
	var tmpdb models.DBAbs
	tmpdb.StatusCode = 200
	tmpdb.Status = "success"
	tmpdb.Message = "Insert successful"
	tmpdb.Count = 0
	tmpdb.Data = "{}"

	var tmpdb2 models.DBAbs
	tmpdb2.Status = "success"
	tmpdb2.StatusCode = 200
	tmpdb2.Message = "Select successful"
	tmpdb2.Data = `{"difficulty":3,"name":"test2","prep_time":"30 mins","recipe_id":9,"vegetarian":true}`
	tmpdb2.Count = 1


	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	{name:"Test_1_succes_dbAbs",args:args{tmpdb},want:`{"status":"success","status_code":200,"message":"Insert successful","result_data":{}}`},
	{name:"Test_2_success_get",args:args{tmpdb2},want:`{"status":"success","status_code":200,"message":"Select successful","result_data":{"difficulty":3,"name":"test2","prep_time":"30 mins","recipe_id":9,"vegetarian":true}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareResponseFromDB(tt.args.dbabs); got != tt.want {
				t.Errorf("PrepareResponseFromDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOffsetLimitStr(t *testing.T) {
	type args struct {
		filterstr string
	}
	tests := []struct {
		name    string
		args    args
		wantRet string
		wantErr bool
	}{
	// TODO: Add test cases.
	{name:"Test_1",args:args{"/recipe?offset=5&limit=10"},wantRet:" OFFSET 5 LIMIT 10",wantErr:false},
	{name:"Test_2",args:args{"/recipe?offset=6"},wantRet:" OFFSET 6 LIMIT 20",wantErr:false},
	{name:"Test_3",args:args{"/recipe?limit=5"},wantRet:" LIMIT 5",wantErr:false},
	{name:"Test_4",args:args{"/recipe?limit=hh"},wantRet:" LIMIT 20",wantErr:false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := OffsetLimitStr(tt.args.filterstr)
			if (err != nil) != tt.wantErr {
				t.Errorf("OffsetLimitStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRet != tt.wantRet {
				t.Errorf("OffsetLimitStr() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}
