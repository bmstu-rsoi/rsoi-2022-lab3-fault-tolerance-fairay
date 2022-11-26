package models

import (
	"gateway/objects"
	"gateway/repository"
)

type FlightsM struct {
	rep repository.FlightsRep
}

func NewFlightsM(rep repository.FlightsRep) *FlightsM {
	return &FlightsM{rep}
}

func (model *FlightsM) Fetch(page int, page_size int) (*objects.PaginationResponse, error) {
	return model.rep.GetAll(page, page_size)
}

func (model *FlightsM) Find(flight_number string) (*objects.FlightResponse, error) {
	return model.rep.Find(flight_number)
}
