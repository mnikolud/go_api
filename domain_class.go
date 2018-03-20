package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Domain struct
type Domain struct {
	ID         int      `json:"id"` //only int accepted!!!
	Name       string   `json:"name"`
	Expiration string   `json:"expiration"` //time.Time??? how to get from scan?
	Owner      *Contact `json:"owner"`
}

//Contact struct
type Contact struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

//Get All Domains
func (a *App) getDomains(w http.ResponseWriter, r *http.Request) {
	domains, err := getDomains(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, domains)
}

//Get Specific Domain
func (a *App) getDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}
	dom := Domain{ID: id, Owner: new(Contact)} //init pointer!!!
	if err := dom.getDomain(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Domain not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, dom)
}

//Create a Domain
func (a *App) createDomain(w http.ResponseWriter, r *http.Request) {
	var dom Domain
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dom); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request input")
		return
	}
	defer r.Body.Close()
	if err := dom.createDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, dom)
}

//Update a Domain
func (a *App) updateDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}
	dom := Domain{ID: id, Owner: new(Contact)} //init pointer!!!
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&dom); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest input")
		return
	}
	defer r.Body.Close()
	if err := dom.updateDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, dom)
}

//delete a Domain
func (a *App) deleteDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid domain ID")
		return
	}
	dom := Domain{ID: id}
	if err := dom.deleteDomain(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
