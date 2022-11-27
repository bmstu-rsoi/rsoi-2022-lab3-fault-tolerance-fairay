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

func (model *FlightsM) Fetch(page int, pageSize int) (*objects.PaginationResponse, error) {
	return model.rep.GetAll(page, pageSize)
}

func (model *FlightsM) Find(flightNumber string) (*objects.FlightResponse, error) {
	return model.rep.Find(flightNumber)
}
