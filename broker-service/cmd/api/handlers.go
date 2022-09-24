package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"

	"github.com/raykrishardi/broker/event"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
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
		// app.logItem(w, requestPayload.Log) // SENDING POST TO log-service
		// app.logEventViaRabbit(w, requestPayload.Log) // SENDING MESSAGE TO RABBITMQ, then the listener-service will be notified if there's messages and consume it
		app.logItemViaRPC(w, requestPayload.Log) // SENDING to another go microservice via RPC
	case "mail":
		app.sendMail(w, requestPayload.Mail)
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

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	// Marshall/convert to json the payload
	jsonData, err := json.MarshalIndent(&msg, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Call the mail microservice
	request, err := http.NewRequest("POST", "http://mail-service/send", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	// Create http client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		// Something went wrong on the server side
		app.errorJSON(w, errors.New("error calling mail service"))
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
		app.errorJSON(w, errors.New("get error response from the mail service"), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromService)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	jsonResponse := jsonResponse{
		Error:   false,
		Message: "Logged via rabbitmq",
	}

	app.writeJSON(w, http.StatusAccepted, jsonResponse)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO") // Severity -> topic name
	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string

	// RPCServer -> the struct that's created on the server end
	// Naming of the function very important here!
	// Then just pass the arguments of the RPCServer.LogInfo function
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
