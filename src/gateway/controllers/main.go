package controllers

import (
	"fmt"
	"gateway/models"
	"gateway/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func initControllers(r *mux.Router, models *models.Models) {
	r.Use(utils.LogHandler)

	r.HandleFunc("/manage/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	api1R := r.PathPrefix("/api/v1/").Subrouter()
	InitFlights(api1R, models.Flights)
	InitPrivileges(api1R, models.Privileges)
	InitTickets(api1R, models.Tickets)
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()
	models := models.InitModels()

	initControllers(router, models)
	return router
}

func RunRouter(r *mux.Router, port uint16) error {
	c := cors.New(cors.Options{})
	handler := c.Handler(r)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), handler)
}
