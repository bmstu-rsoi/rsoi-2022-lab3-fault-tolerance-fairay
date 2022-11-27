package controllers

import (
	"gateway/controllers/responses"
	"gateway/errors"
	"gateway/models"

	"net/http"

	"github.com/gorilla/mux"
)

type privilegeCtrl struct {
	privileges *models.PrivilegesM
}

func InitPrivileges(r *mux.Router, privileges *models.PrivilegesM) {
	ctrl := &privilegeCtrl{privileges}
	r.HandleFunc("/privilege", ctrl.fetch).Methods("GET")
}

func (ctrl *privilegeCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.privileges.Fetch(r.Header.Get("X-User-Name"))
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.BonusUnavailable:
		responses.ServiceUnavailable(w, err.Error())
	default:
		responses.InternalError(w)
	}
}
