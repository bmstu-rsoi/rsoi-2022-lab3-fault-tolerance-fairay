package models

import (
	"gateway/repository"
	"gateway/utils"
	"net/http"
)

type Models struct {
	Flights    *FlightsM
	Privileges *PrivilegesM
	Tickets    *TicketsM
}

func InitModels() *Models {
	models := new(Models)
	client := &http.Client{}

	models.Flights = NewFlightsM(repository.NewCBFlightsRep(utils.Config.FlightsEndpoint))
	models.Privileges = NewPrivilegesM(client)
	models.Tickets = NewTicketsM(client, models.Flights)

	return models
}
