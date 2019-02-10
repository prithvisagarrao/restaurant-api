package models

import (
	"fmt"
	"strings"
	"github.com/hellofreshdevtests/prithvisagarrao-api-test/config"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"reflect"
	"net/http"
	"github.com/antigloss/go/logger"
	"encoding/json"
)

//InsertQuery generates the Insert query string
func InsertQuery(reqAbs ReqAbs, table string) (query string){

	column := ""
	value := ""
	tmpval := ""
	tmpcol := ""


	for c, v:= range reqAbs.Payload{

		switch v.(type){

		case string:
			value = fmt.Sprintf(`'%v',`,v)
			column = fmt.Sprintf("%v,",c)

		case int, int32, int64, float32, float64:
			value = fmt.Sprintf("%v,",v)
			column = fmt.Sprintf("%v,",c)

		case bool:
			value = fmt.Sprintf("%v,",v)
			column = fmt.Sprintf("%v,",c)

		default:
			value = fmt.Sprintf(`'%v',`,v)
			column = fmt.Sprintf("%v,",c)
			fmt.Printf("\ndefault --reflect %v\n",reflect.TypeOf(v))
		}

		tmpval += value
		tmpcol += column


	}

	tmpval = strings.Trim(tmpval,",")
	tmpcol = strings.Trim(tmpcol,",")
	query = fmt.Sprintf("INSERT INTO %v ( %v ) VALUES ( %v );",table,tmpcol,tmpval)

	return
}

//UpdateQuery generates the update query string
func UpdateQuery(reqAbs ReqAbs, recipeId string, table string ) (query string){

	query = ""
	name := ""
	where := ""

	for k,v := range reqAbs.Payload{

		switch v.(type) {

		case string:
			name += fmt.Sprintf(`%v = '%v',`,k,v)

		case int, int32, int64, float32, float64:
			name += fmt.Sprintf("%v = %v,",k,v)

		case bool:
			name += fmt.Sprintf("%v = %v,",k,v)

		default:
			fmt.Printf("\ndefault --reflect %v\n",reflect.TypeOf(v))
			logger.Error("found unsupported type")
		}
	}

	name = strings.Trim(name,",")
	where = fmt.Sprintf("WHERE recipe_id = %v ;",recipeId)

	query = fmt.Sprintf("UPDATE %v SET %v %v",table, name, where)

	return
}

//Insert function performs the insert record operation on postgres
func Insert(dbAbs *DBAbs) {

	dbConn , err := GetDbConn()
	if err != nil {
		logger.Error(fmt.Sprintf("ErrorType : INFRA_ERROR, Not able to connect to DB Server. Got error : %v", err.Error()))
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Status = "fail"
		dbAbs.Message = "Unable to connect to Postgres"

	}

	logger.Info("Running insert query on postgres")
	_, err = dbConn.Exec(dbAbs.Query)

	if err != nil{
		if strings.Contains(err.Error(),"duplicate key value violates unique constraint"){
			logger.Error("Postgres insert query failed, please provide a unique recipe name")
			//log.Fatal(err)
			dbAbs.Status = "fail"
			dbAbs.StatusCode = http.StatusForbidden
			dbAbs.Data = "{}"
			dbAbs.Message = "Insert failed"
		} else {
			logger.Error("Postgres insert query failed")
			dbAbs.Status = "fail"
			dbAbs.StatusCode = http.StatusInternalServerError
			dbAbs.Data = "{}"
			dbAbs.Message = "Insert failed"
		}
	} else {
		logger.Info("Postgres insert query success")
		dbAbs.Status = "success"
		dbAbs.StatusCode = http.StatusOK
		dbAbs.Data = "{}"
		dbAbs.Message = "Insert successful"
	}

}

//Select function performs select operation on postgres and returns the records selected
func Select (dbAbs *DBAbs){

	var postoResult []map[string]interface{}
	dbConn , err := GetDbConn()
	if err != nil {
		logger.Error(fmt.Sprintf("ErrorType : INFRA_ERROR, Not able to connect to DB Server. Got error : %v", err.Error()))
		dbAbs.Status = "fail"
		dbAbs.Message = "Error at in getting connection"
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Data = "{}"
		dbAbs.Count = 0
		return
	}


	logger.Info("Running select query on postgres")
	rows, err := dbConn.Query(dbAbs.Query)
	if err != nil {
		logger.Error(fmt.Sprintf("ErrorType : QUERY_ERROR, Postgres query failed with error : %v", err.Error()))
		dbAbs.Status = "fail"
		dbAbs.Message = "Error at running query"
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Data = "{}"
		dbAbs.Count = 0
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		logger.Error(fmt.Sprintf("ErrorType : QUERY_ERROR, Postgres query failed"))
		dbAbs.Status = "fail"
		dbAbs.Message = "Error at retrieving data"
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Data = "{}"
		dbAbs.Count = 0
		return
	}

	data := make([]interface{}, len(cols))
	args := make([]interface{}, len(data))

	for i := range data {
		args[i] = &data[i]
	}

	for rows.Next() {

		var rowData = make(map[string]interface{})

		if err := rows.Scan(args...); err != nil {
			logger.Error(fmt.Sprintf("ErrorType : QUERY_ERROR, Postgres query failed, Error when fetching data. Error : %v ", err.Error()))
			return
		}

		for i := range data {

			rowData[cols[i]] = data[i]
		}

		postoResult = append(postoResult, rowData)

	}

	rows.Close()

	if err != nil {
		dbAbs.Status = "fail"
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Message = "Error connecting to endpoint"
		dbAbs.Data = "{}"
		dbAbs.Count = 0
	} else {
		data, err := json.Marshal(postoResult)
		if err != nil {
			logger.Error(fmt.Sprintf("Error in encoding json",err.Error()))
		}

		logger.Info("Select query successful on postgres")
		dbAbs.Status = "success"
		dbAbs.StatusCode = http.StatusOK
		dbAbs.Data = string(data)
		dbAbs.RichData = postoResult



		dbAbs.Message = "Select successful"
		if len(postoResult) <= 0{
			dbAbs.Message = "No records found"
		}
		dbAbs.Count = uint64(len(postoResult))
	}


}

//Update function performs the update operation on postgres
func Update(dbAbs *DBAbs){
	dbConn , err := GetDbConn()
	if err != nil {

		logger.Error(fmt.Sprintf("ErrorType : INFRA_ERROR, Not able to connect to DB Server. Got error : %v", err.Error()))
		os.Exit(1)
	}


	r, err := dbConn.Exec(dbAbs.Query)

	if err != nil{

		logger.Error(fmt.Sprintf("Could not perform update query. Got error %v",err.Error()))
		dbAbs.Status = "fail"
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Message = "Update query failed"
		dbAbs.Count = 0

	} else {
		logger.Info("Update query successful")
		dbAbs.StatusCode = http.StatusOK
		dbAbs.Status = "success"
		dbAbs.Message = "update successful"
		rc, _ := r.RowsAffected()
		if rc <= 0{
			dbAbs.Message = "no recipes found to update"
		}
		dbAbs.Count = uint64(rc)
	}
}

//Delete function performs the delete operation on postgres and deletes a record
func Delete(dbAbs *DBAbs) {

	dbConn , err := GetDbConn()
	if err != nil {

		logger.Error(fmt.Sprintf("ErrorType : INFRA_ERROR, Not able to connect to DB Server. Got error : %v", err.Error()))
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Status = "fail"
		dbAbs.Message = "Unable to connect to Postgres"
	}

	ks, err := dbConn.Exec(dbAbs.Query)


	if err != nil{
		logger.Error(fmt.Sprintf("delete query failed with error:%v",err))
		dbAbs.Status = "fail"
		dbAbs.StatusCode = http.StatusNotFound
		dbAbs.Data = "{}"
		dbAbs.Message = "Delete failed"
	} else {
		cnt, err1 := ks.RowsAffected()
		if err1 != nil{
			logger.Error(fmt.Sprintf("Getting row count for delete query failed with error: %v",err))
			dbAbs.Status = "fail"
			dbAbs.StatusCode = http.StatusInternalServerError
			dbAbs.Data = "{}"
			dbAbs.Message = "Getting row count for delete failed"
		}
		dbAbs.Count = uint64(cnt)

		if dbAbs.Count > 0 {
			logger.Info("Delete query successful")
			dbAbs.Status = "success"
			dbAbs.StatusCode = http.StatusOK
			dbAbs.Data = "{}"
			dbAbs.Message = "Delete successful"
		} else {
			logger.Info("Delete query successful")
			dbAbs.Status = "success"
			dbAbs.StatusCode = http.StatusOK
			dbAbs.Data = "{}"
			dbAbs.Message = "No records found to delete"
		}
	}

}

//GetDbConn is used to get sql db connections to the postgres database
func GetDbConn() (dbConn *sql.DB, err error){

	Conf := config.GetConfig()

	dbName := Conf.Postgres.Database
	dbUser := Conf.Postgres.Username
	dbPass := Conf.Postgres.Password
	dbHost := Conf.Postgres.Host
	dbPort := Conf.Postgres.Port

	var dbInfo string

	if len(dbPass) == 0 {

		dbInfo = fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s",
			dbUser, dbName, dbHost, dbPort)
	} else {

		dbInfo = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
			dbUser, dbPass, dbName, dbHost, dbPort)
	}

	dbConn, err = sql.Open("postgres", dbInfo)
	return
}