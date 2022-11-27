package controllers

import (
	"gateway/controllers/responses"
	"gateway/errors"
	"gateway/models"
	"gateway/objects"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ticketsCtrl struct {
	tickets *models.TicketsM
}

func InitTickets(r *mux.Router, tickets *models.TicketsM) {
	ctrl := &ticketsCtrl{tickets: tickets}
	r.HandleFunc("/me", ctrl.me).Methods("GET")
	r.HandleFunc("/tickets", ctrl.fetch).Methods("GET")
	r.HandleFunc("/tickets", ctrl.post).Methods("POST")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.get).Methods("GET")
	r.HandleFunc("/tickets/{ticketUid}", ctrl.delete).Methods("DELETE")
}

func (ctrl *ticketsCtrl) me(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-User-Name")
	data, err := ctrl.tickets.FetchUser(username)
	if err != nil {
		responses.InternalError(w)
	} else {
		responses.JsonSuccess(w, data)
	}
}

func (ctrl *ticketsCtrl) fetch(w http.ResponseWriter, r *http.Request) {
	data, err := ctrl.tickets.Fetch()
	if err != nil {
		responses.InternalError(w)
	} else {
		responses.JsonSuccess(w, data)
	}
}

func (ctrl *ticketsCtrl) post(w http.ResponseWriter, r *http.Request) {
	reqBody := new(objects.TicketPurchaseRequest)
	err := json.NewDecoder(r.Body).Decode(reqBody)
	if err != nil {
		responses.ValidationErrorResponse(w, err.Error())
		return
	}

	data, err := ctrl.tickets.Create(reqBody.FlightNumber, r.Header.Get("X-User-Name"), reqBody.Price, reqBody.PaidFromBalance)
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.FlightUnavailable, errors.BonusUnavailable, errors.TicketsUnavailable:
		responses.ServiceUnavailable(w, err.Error())
	default:
		responses.InternalError(w)
	}
}

func (ctrl *ticketsCtrl) get(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	ticketUid := urlParams["ticketUid"]
	username := r.Header.Get("X-User-Name")

	data, err := ctrl.tickets.Find(ticketUid, username)
	switch err {
	case nil:
		responses.JsonSuccess(w, data)
	case errors.ForbiddenTicket:
		responses.Forbidden(w)
	default:
		responses.RecordNotFound(w, ticketUid)
	}
}

func (ctrl *ticketsCtrl) delete(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	ticketUid := urlParams["ticketUid"]
	username := r.Header.Get("X-User-Name")

	err := ctrl.tickets.Delete(ticketUid, username)
	switch err {
	case nil:
		responses.SuccessTicketDeletion(w)
	case errors.ForbiddenTicket:
		responses.Forbidden(w)
	default:
		responses.RecordNotFound(w, ticketUid)
	}
}
