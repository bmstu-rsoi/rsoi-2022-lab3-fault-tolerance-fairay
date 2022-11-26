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

type PrivilegesRep interface {
	GetAll(username string) (*objects.PrivilegeInfoResponse, error)
	Add(username string, request *objects.AddHistoryRequest) (*objects.AddHistoryResponce, error)
	Delete(username string, ticketUid string) error
}

type CBPrivilegesRep struct {
	cb       *gobreaker.CircuitBreaker
	endpoint string
	client   *http.Client
}

func NewCBPrivilegesRep(endpoint string) *CBPrivilegesRep {
	settings := gobreaker.Settings{Name: "Privileges"}
	cb := gobreaker.NewCircuitBreaker(settings)
	client := &http.Client{Timeout: 2 * time.Second}
	return &CBPrivilegesRep{cb, endpoint, client}
}

func (rep *CBPrivilegesRep) cbExecute(req *http.Request) (interface{}, error) {
	return rep.cb.Execute(func() (interface{}, error) {
		resp, err := rep.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
}

func (rep *CBPrivilegesRep) GetAll(username string) (*objects.PrivilegeInfoResponse, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/privilege", rep.endpoint), nil)
	req.Header.Add("X-User-Name", username)

	body, err := rep.cbExecute(req)
	if err != nil {
		return nil, err
	}

	data := &objects.PrivilegeInfoResponse{}
	err = json.Unmarshal(body.([]byte), data)
	return data, err
}

func (rep *CBPrivilegesRep) Add(username string, request *objects.AddHistoryRequest) (*objects.AddHistoryResponce, error) {
	req_body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/history", rep.endpoint), bytes.NewBuffer(req_body))
	req.Header.Add("X-User-Name", username)

	body, err := rep.cbExecute(req)
	if err != nil {
		return nil, err
	}

	data := &objects.AddHistoryResponce{}
	json.Unmarshal(body.([]byte), data)
	return data, nil
}

func (rep *CBPrivilegesRep) Delete(username string, ticketUid string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/history/%s", rep.endpoint, ticketUid), nil)
	req.Header.Add("X-User-Name", username)
	_, err := rep.cbExecute(req)
	return err
}
