package controller

import (
	"net/http"
	"github.com/hellofreshdevtests/prithvisagarrao-api-test/models"
	"github.com/hellofreshdevtests/prithvisagarrao-api-test/utils"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"github.com/antigloss/go/logger"
	"math"
)

//CreateRecipe creates a new recipe
func CreateRecipe(w http.ResponseWriter, r *http.Request){

	var dbAbs models.DBAbs
	reqAbs := utils.MapReqAbs(r)

	for k, v := range reqAbs.Payload{
		if k=="vegetarian"{
			if _,ok := v.(bool); !ok{
				logger.Error("Invalid data sent to create recipe")
				dbAbs.StatusCode = http.StatusBadRequest
				dbAbs.Status = "fail"
				dbAbs.Message = "data invalid for vegetarian ; must be true or false"

				respAbs := utils.PrepareResponseFromDB(dbAbs)
				utils.Respond(w, respAbs)
				return
			}
		}
		if k=="difficulty"{
			if _,ok := v.(float64); ok{
				//TODO need to change this(for decimals)
				if (v.(float64) < 1.0 || v.(float64) > 3.0) && math.Mod(v.(float64) *10,5 )== 0 {
					logger.Error("Invalid data sent to create recipe")
					dbAbs.StatusCode = http.StatusBadRequest
					dbAbs.Status = "fail"
					dbAbs.Message = "data invalid for difficulty ; must be between 1 and 4"

					respAbs := utils.PrepareResponseFromDB(dbAbs)
					utils.Respond(w, respAbs)
					return
				}
			} else {
				logger.Error("Invalid data sent to create recipe")
				dbAbs.StatusCode = http.StatusBadRequest
				dbAbs.Status = "fail"
				dbAbs.Message = "data type invalid for difficulty"

				respAbs := utils.PrepareResponseFromDB(dbAbs)
				utils.Respond(w, respAbs)
				return
			}
		}
	}

	query := models.InsertQuery(reqAbs,"recipes")

	dbAbs.Query = query

	logger.Info("Running query to create user")
	models.Insert(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)


}

//ListAllRecipes lists all the available recipes
func ListAllRecipes( w http.ResponseWriter, r *http.Request){
	var dbAbs models.DBAbs
	reqAbs := utils.MapReqAbs(r)

	filter, err := utils.OffsetLimitStr(reqAbs.Filter)

	if err != nil {
		logger.Error("filter provided in wrong format")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "filter invalid"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}

	selectStr := fmt.Sprintf("SELECT * FROM recipes %v",filter)
	dbAbs.Query = selectStr
	logger.Info("Running query to list all recipes")
	models.Select(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)

}

//GetRecipe gets a recipe associated with the given recipe_id
func GetRecipe(w http.ResponseWriter, r *http.Request){
	var dbAbs models.DBAbs
	var selectStr string

	vars := mux.Vars(r)
	recipeId := vars["id"]

	if rId, err := strconv.Atoi(recipeId); err != nil{
		logger.Error("Invalid Id sent in request")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "recipe_id invalid"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	} else {

		selectStr = fmt.Sprintf("SELECT * FROM recipes WHERE recipe_id = %v;", rId)
	}


	dbAbs.Query = selectStr
	logger.Info("Running query to select recipe")
	models.Select(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
}

//DeleteRecipe deletes a recipe associated with the given recipe_id
func DeleteRecipe(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	recipeId := vars["id"]
	var dbAbs models.DBAbs

	var deleteStrRate string

	if rId, err := strconv.Atoi(recipeId); err != nil{
		logger.Error("recipe id provided in wrong format")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "recipe_id invalid"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	} else {


		deleteStrRate = fmt.Sprintf("DELETE FROM ratings WHERE recipe_id = %v;", rId)
	}


	dbAbs.Query = deleteStrRate
	logger.Info("Running query to delete ratings first")

	models.Delete(&dbAbs)

	if dbAbs.Status != "success"{
		logger.Error(fmt.Sprintf("ratings for the recipe id %v could not be deleted",recipeId))
		dbAbs.StatusCode = http.StatusInternalServerError
		dbAbs.Status = "fail"
		dbAbs.Message = "ratings could not be deleted"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return

	}

	logger.Info(fmt.Sprintf("Ratings for recipe id: %v deleted successfully",recipeId))
	deleteStr := fmt.Sprintf("DELETE FROM recipes where recipe_id = %v;",recipeId)
	dbAbs.Query = deleteStr
	logger.Info("Running query to delete recipe")

	models.Delete(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
}

//UpdateRecipe update a recipe associated with the given recipe_id
func UpdateRecipe(w http.ResponseWriter, r *http.Request){

	var dbAbs models.DBAbs
	vars := mux.Vars(r)
	recipeId := vars["id"]
	reqAbs := utils.MapReqAbs(r)


	if _, err := strconv.Atoi(recipeId); err != nil{
		logger.Error("recipe id provided in wrong format")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "recipe_id invalid"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}


	for k, v := range reqAbs.Payload{
		if k == "vegetarian"{
			if _,ok := v.(bool); !ok{
				logger.Error("Invalid data sent")
				dbAbs.StatusCode = http.StatusBadRequest
				dbAbs.Status = "fail"
				dbAbs.Message = "data invalid for vegetarian ; must be true or false"

				respAbs := utils.PrepareResponseFromDB(dbAbs)
				utils.Respond(w, respAbs)
				return
			}
		}

		if k=="difficulty"{
			if _,ok := v.(float64); ok{
				if v.(float64) < 1.0 || v.(float64) > 3.0 {
					logger.Error("Invalid data sent")
					dbAbs.StatusCode = http.StatusBadRequest
					dbAbs.Status = "fail"
					dbAbs.Message = "data invalid for difficulty ; must be between 1 and 4"

					respAbs := utils.PrepareResponseFromDB(dbAbs)
					utils.Respond(w, respAbs)
					return
				}
			} else {
				logger.Error("Invalid data sent")
				dbAbs.StatusCode = http.StatusBadRequest
				dbAbs.Status = "fail"
				dbAbs.Message = "data type invalid for difficulty"

				respAbs := utils.PrepareResponseFromDB(dbAbs)
				utils.Respond(w, respAbs)
				return
			}
		}
	}

	query := models.UpdateQuery(reqAbs, recipeId, "recipes")

	dbAbs.Query = query
	models.Update(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
}

//RateRecipe rates a recipe associated with the given recipe_id
func RateRecipe(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	recipeId := vars["id"]
	var dbAbs models.DBAbs
	reqAbs := utils.MapReqAbs(r)

	if _, err := strconv.Atoi(recipeId); err != nil{
		logger.Error("recipe id provided in wrong format")
		dbAbs.StatusCode = http.StatusBadRequest
		dbAbs.Status = "fail"
		dbAbs.Message = "recipe_id invalid"
		respAbs := utils.PrepareResponseFromDB(dbAbs)
		utils.Respond(w, respAbs)
		return
	}


	for k, v := range reqAbs.Payload{
		if k=="rate"{
			if _,ok := v.(float64); ok{
				if (v.(float64) < 1.0 || v.(float64) > 5.0) && math.Mod(v.(float64) *10,10 )== 0  {
					logger.Error("Invalid value sent for rating the recipe")
					dbAbs.StatusCode = http.StatusBadRequest
					dbAbs.Status = "fail"
					dbAbs.Message = "data invalid"

					respAbs := utils.PrepareResponseFromDB(dbAbs)
					utils.Respond(w, respAbs)
					return
				}
			} else {
				logger.Error("Invalid data sent")
				dbAbs.StatusCode = http.StatusBadRequest
				dbAbs.Status = "fail"
				dbAbs.Message = "data type invalid for rate"

				respAbs := utils.PrepareResponseFromDB(dbAbs)
				utils.Respond(w, respAbs)
				return
			}
		}
	}

	reqAbs.Payload["recipe_id"] = recipeId

	query := models.InsertQuery(reqAbs,"ratings")
	dbAbs.Query = query
	logger.Info("Running query to rate recipe")

	models.Insert(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)

	utils.Respond(w, respAbs)

}

func SearchRecipe(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]
	var dbAbs models.DBAbs


	selectStr := fmt.Sprintf("SELECT * FROM recipes WHERE name LIKE '%%%v%%';",name)

	dbAbs.Query = selectStr
	logger.Info("Running query to select recipe")
	models.Select(&dbAbs)
	respAbs := utils.PrepareResponseFromDB(dbAbs)
	utils.Respond(w, respAbs)
}
