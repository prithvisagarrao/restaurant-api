package main

import (
	"fmt"
	"github.com/antigloss/go/logger"
	"github.com/gorilla/mux"
	"github.com/prithvisagarrao/restaurant-api/auth"
	"github.com/prithvisagarrao/restaurant-api/config"
	"github.com/prithvisagarrao/restaurant-api/controller"
	"net/http"
)

func main() {

	config.InitConfig()

	//Conf := config.GetConfig()
	//targets.Handlehttp()


	r := mux.NewRouter()
	r.Use(auth.JWTAuthentication)
	r.HandleFunc("/recipes", controller.ListAllRecipes).Methods("GET")
	r.HandleFunc("/recipes", controller.CreateRecipe).Methods("POST")
	r.HandleFunc("/recipes/{id}", controller.GetRecipe).Methods("GET")
	r.HandleFunc("/recipes/{id}", controller.UpdateRecipe).Methods("PATCH", "PUT")
	r.HandleFunc("/recipes/{id}", controller.DeleteRecipe).Methods("DELETE")
	r.HandleFunc("/recipes/{id}/rating", controller.RateRecipe).Methods("POST")
	r.HandleFunc("/createUser", controller.CreateUser).Methods("POST")
	r.HandleFunc("/recipes/search/{name}",controller.SearchRecipe).Methods("GET")
	r.HandleFunc("/login", controller.Login).Methods("POST")

	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Error(fmt.Sprintf("Error in ListenAndServe %v", err))
	}

}
