package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gateway/errors"
	"gateway/objects"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	body, err := rep.cb.Execute(func() (interface{}, error) {
		resp, err := rep.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	})
	if err != nil {
		err = errors.BonusUnavailable
	}
	return body, err
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
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/history", rep.endpoint), bytes.NewBuffer(reqBody))
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
	go retry.Do(func() error {
		fmt.Println("Trying to delete")
		_, err := rep.cbExecute(req)
		return err
	}, retry.Delay(15*time.Second))

	return nil
}
