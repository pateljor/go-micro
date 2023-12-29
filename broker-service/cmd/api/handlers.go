package main

import (
	"broker/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
	Mail MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}


func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker, nice ur sick",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action{
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown acrtion"))

	}
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
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
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload){
	// create some json we'll send to auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	authServiceURL := "http://authentication-service/authenticate"
	request, err := http.NewRequest("POST", authServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		log.Println("Errror creating request", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		log.Println("Errror getting resposne", err)
		return
	}
	defer response.Body.Close()


	log.Println("Error:", err)
	// make sure we get the correct status code 
	if response.StatusCode == http.StatusUnauthorized {
		log.Println("Wrong Status Code", response.StatusCode)
		app.errorJSON(w, errors.New("invalid creds"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		log.Println("Wrong Status Code 2", response.StatusCode)
		app.errorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode the json from the auth service 

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println("Error creating decoder", err)
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error{
		log.Println("Error from jsonFromService", )
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// call mail service
	mailServiceURL := "http://mail-service/send"
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		log.Println("Errror creating request", err)
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")

	// set client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		log.Println("Errror getting response (mail)", err)
		return
	}
	defer response.Body.Close()
	
	// error catch for non accepted status codes
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mailer service"))
		return
	}
	
	// setup & send  payload
	var payload jsonResponse
	payload.Error = false
	payload.Message = "sent mail"
	payload.Data = msg

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil{
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload {
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
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
	client, err := rpc.Dial("tpc", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := RPCPayload {
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: result,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}