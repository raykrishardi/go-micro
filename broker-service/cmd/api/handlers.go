package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	// Unmarshal the json
	requestPayload := RequestPayload{}
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	default:
		if err != nil {
			app.errorJSON(w, errors.New("unknown action"))
			return
		}
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// Marshall/convert to json the auth payload which consists of email and password
	jsonData, err := json.MarshalIndent(&a, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Call the auth microservice
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		// Something went wrong on the server side
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	jsonFromService := jsonResponse{}

	// decode/unmarshal the json from the auth service (basically this is an alternative to UNMARSHALLING)
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New("get error response from the auth service"), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	// Marshall/convert to json the payload
	jsonData, err := json.MarshalIndent(&entry, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Call the logger microservice
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		// Something went wrong on the server side
		app.errorJSON(w, errors.New("error calling logger service"))
		return
	}

	// create a variable we'll read response.Body into
	jsonFromService := jsonResponse{}

	// decode/unmarshal the json from the auth service (basically this is an alternative to UNMARSHALLING)
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New("get error response from the logger service"), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromService)
}
