package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/objects"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type TicketsRep interface {
	GetAll(username string) (objects.TicketArr, error)
	Create(flightNumber string, price int, username string) (*objects.TicketCreateResponse, error)
	Find(ticketUid string) (*objects.Ticket, error)
	Delete(ticketUid string) error
}

type CBTicketsRep struct {
	cb       *gobreaker.CircuitBreaker
	endpoint string
	client   *http.Client
}

func NewCBTicketsRep(endpoint string) *CBTicketsRep {
	settings := gobreaker.Settings{Name: "Tickets"}
	cb := gobreaker.NewCircuitBreaker(settings)
	client := &http.Client{Timeout: 2 * time.Second}
	return &CBTicketsRep{cb, endpoint, client}
}

func (rep *CBTicketsRep) cbExecute(req *http.Request) (interface{}, error) {
	return rep.cb.Execute(func() (interface{}, error) {
		resp, err := rep.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
}

func (rep *CBTicketsRep) GetAll(username string) (objects.TicketArr, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tickets", rep.endpoint), nil)
	if username != "" {
		req.Header.Set("X-User-Name", username)
	}

	body, err := rep.cbExecute(req)
	if err != nil {
		return nil, err
	}

	data := new(objects.TicketArr)
	json.Unmarshal(body.([]byte), data)
	return *data, nil
}

func (rep *CBTicketsRep) Create(flightNumber string, price int, username string) (*objects.TicketCreateResponse, error) {
	reqBody, err := json.Marshal(&objects.TicketCreateRequest{FlightNumber: flightNumber, Price: price})
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/tickets", rep.endpoint), bytes.NewBuffer(reqBody))
	req.Header.Add("X-User-Name", username)

	body, err := rep.cbExecute(req)
	if err != nil {
		return nil, err
	}

	data := &objects.TicketCreateResponse{}
	err = json.Unmarshal(body.([]byte), data)
	return data, err
}

func (rep *CBTicketsRep) Find(ticketUid string) (*objects.Ticket, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tickets/%s", rep.endpoint, ticketUid), nil)

	body, err := rep.cbExecute(req)
	if err != nil {
		return nil, err
	}

	data := &objects.Ticket{}
	json.Unmarshal(body.([]byte), data)
	return data, nil
}

func (rep *CBTicketsRep) Delete(ticketUid string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/tickets/%s", rep.endpoint, ticketUid), nil)
	_, err := rep.client.Do(req)
	return err
}
