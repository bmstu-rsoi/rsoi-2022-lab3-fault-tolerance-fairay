package models

import (
	"gateway/repository"
	"gateway/utils"
)

type Models struct {
	Flights    *FlightsM
	Privileges *PrivilegesM
	Tickets    *TicketsM
}

func InitModels() *Models {
	flightsRep := repository.NewCBFlightsRep(utils.Config.FlightsEndpoint)
	privilegesRep := repository.NewCBPrivilegesRep(utils.Config.PrivilegesEndpoint)
	ticketsRep := repository.NewCBTicketsRep(utils.Config.TicketsEndpoint)

	return &Models{
		Flights:    NewFlightsM(flightsRep),
		Privileges: NewPrivilegesM(privilegesRep),
		Tickets:    NewTicketsM(ticketsRep, flightsRep, privilegesRep),
	}
}
