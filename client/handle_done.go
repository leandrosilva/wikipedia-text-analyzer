package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// DoneRequest is the payload received from analyzer saying it's done
type DoneRequest struct {
	TargetURL string `json:"targetURL"`
	ResultURL string `json:"resultURL"`
}

// DoneResponse is what we say when they pull our hook
type DoneResponse struct {
	Acknowledge bool `json:"acknowledge"`
}

func handleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := getDoneRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := done(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getDoneRequest(r *http.Request) (DoneRequest, error) {
	var request DoneRequest

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return request, err
	}

	err = json.Unmarshal(body, &request)

	return request, err
}

func done(req DoneRequest) (DoneResponse, error) {
	log.Println("Hook was pulled:", req.TargetURL, "=>", req.ResultURL)
	response := DoneResponse{Acknowledge: true}

	return response, nil
}
