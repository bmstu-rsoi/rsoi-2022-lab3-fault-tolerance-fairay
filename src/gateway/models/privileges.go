package models

import (
	"gateway/objects"
	"gateway/repository"
)

type PrivilegesM struct {
	rep repository.PrivilegesRep
}

func NewPrivilegesM(rep repository.PrivilegesRep) *PrivilegesM {
	return &PrivilegesM{rep: rep}
}

func (model *PrivilegesM) Fetch(username string) (*objects.PrivilegeInfoResponse, error) {
	return model.rep.GetAll(username)
}

func (model *PrivilegesM) AddTicket(username string, request *objects.AddHistoryRequest) (*objects.AddHistoryResponce, error) {
	return model.rep.Add(username, request)
}

func (model *PrivilegesM) DeleteTicket(username string, ticketUid string) error {
	return model.rep.Delete(username, ticketUid)
}
