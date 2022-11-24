package repository

import (
	"encoding/json"
	"fmt"
	"gateway/objects"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sony/gobreaker"
)

type FlightsRep interface {
	GetAll(page int, page_size int) (*objects.PaginationResponse, error)
	Find(flight_number string) (*objects.FlightResponse, error)
}

type CBFlightsRep struct {
	cb       *gobreaker.CircuitBreaker
	endpoint string
	client   *http.Client
}

func NewCBFlightsRep(endpoint string) *CBFlightsRep {
	settings := gobreaker.Settings{Name: "Flights"}
	cb := gobreaker.NewCircuitBreaker(settings)
	client := &http.Client{Timeout: 2 * time.Second}
	return &CBFlightsRep{cb, endpoint, client}
}

func (rep *CBFlightsRep) GetAll(page int, page_size int) (*objects.PaginationResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/flights", rep.endpoint), nil)
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", page))
	q.Add("size", fmt.Sprintf("%d", page_size))
	req.URL.RawQuery = q.Encode()

	body, err := rep.cb.Execute(func() (interface{}, error) {
		resp, err := rep.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
	if err != nil {
		return nil, err
	}

	data := &objects.PaginationResponse{}
	err = json.Unmarshal(body.([]byte), data)
	return data, err
}

func (rep *CBFlightsRep) Find(flight_number string) (*objects.FlightResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/flights/%s", rep.endpoint, flight_number), nil)

	body, err := rep.cb.Execute(func() (interface{}, error) {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
	if err != nil {
		return nil, err
	}

	data := &objects.FlightResponse{}
	err = json.Unmarshal(body.([]byte), data)
	return data, err
}
