package webengage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Webengage interface
type Webengage interface {
	Track(e *Event) error
}

type webengage struct {
	Client *http.Client
	LicenseCode string
	ApiKey string
	ApiUrl string
}

// Event for Webengage
type Event struct {
	UserID    string
	AnonymousID string
	EventName string
	EventTime string
	EventData map[string]interface{}
}

// Constant for WebEngage anonymousID
const AnonID = "anon-user"

// Track sends the event to webengage
func (w *webengage) Track(e *Event) error {
	body := map[string]interface{}{
		"userId":      e.UserID,
		"anonymousId": e.AnonymousID,
		"eventName":   e.EventName,
		"eventTime":   e.EventTime,
		"eventData":   e.EventData,
	}
	return w.send(body)
}

func (w *webengage) send(body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := w.ApiUrl + "/" + w.LicenseCode + "/events"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer " + w.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := w.Client
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if 201 != resp.StatusCode {
		return fmt.Errorf("%s", resBody)
	}
	return nil
}

// New client instance
func New(licenseCode, apiKey, apiUrl string) Webengage {
	return &webengage{
		Client: http.DefaultClient,
		LicenseCode: licenseCode,
		ApiKey: apiKey,
		ApiUrl: apiUrl,
	}
}