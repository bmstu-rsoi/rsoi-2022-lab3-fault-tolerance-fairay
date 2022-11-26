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

	flightsRep := repository.NewCBFlightsRep(utils.Config.FlightsEndpoint)
	privilegesRep := repository.NewCBPrivilegesRep(utils.Config.PrivilegesEndpoint)

	models.Flights = NewFlightsM(flightsRep)
	models.Privileges = NewPrivilegesM(privilegesRep)
	models.Tickets = NewTicketsM(client, flightsRep, privilegesRep)

	return models
}
