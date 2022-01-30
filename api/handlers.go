package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"location-history/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddLocationData(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	// append location to the history only if order_id is available
	if orderId, exist := requestVars["order_id"]; exist {
		locationReq := util.AddLocationRequest{}
		_ = json.NewDecoder(r.Body).Decode(&locationReq)

		// create new location entry
		loc := locationEntry{
			OrderId:   orderId,
			Lat:       locationReq.Lat,
			Lng:       locationReq.Lng,
			Timestamp: time.Now(),
		}

		prevLocEntries, _, _ := getPreviousLocationData()

		// write the updated location history to the json file
		prevLocEntries = append(prevLocEntries, loc)
		updatedLocHistory, _ := json.Marshal(prevLocEntries)
		err := ioutil.WriteFile(filename, updatedLocHistory, 0644)
		if err != nil {
			log.Printf("Error updating location history file: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Bad request: Invalid order_id"))
	}

}

func GetLocationData(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	prevLocEntries, prevOrderIds, err := getPreviousLocationData()
	if err != nil {
		log.Printf("Error reading location history file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// fetch history only if order_id is available and present in the dataset
	if orderId, exist := requestVars["order_id"]; exist && isValidOrderId(orderId, prevOrderIds) {

		// create location history response
		var locResp locHistoryResp
		locResp.OrderId = orderId
		var locHistory []history
		for _, locEntry := range prevLocEntries {
			if locEntry.OrderId == orderId {
				locHistory = append(locHistory, history{
					Lat: locEntry.Lat,
					Lng: locEntry.Lng,
				})
			}
		}

		// return location history based on max value
		max := r.FormValue("max")
		if max == "" {
			locResp.History = locHistory
		} else {
			entriesToReturn, _ := strconv.Atoi(r.FormValue("max"))
			locResp.History = locHistory[:entriesToReturn]
		}

		// send location history response
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(locResp)
		if err != nil {
			log.Printf("Error writing response %v\n", err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Bad request: Invalid order_id"))
	}
}

func DeleteLocationData(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	prevLocEntries, prevOrderIds, err := getPreviousLocationData()
	if err != nil {
		log.Printf("Error reading location history file: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// delete location history only if order_id is available and present in the dataset
	if orderId, exist := requestVars["order_id"]; exist && isValidOrderId(orderId, prevOrderIds) {
		var updatedLocEntries []locationEntry
		for _, locEntry := range prevLocEntries {
			if locEntry.OrderId != orderId {
				updatedLocEntries = append(updatedLocEntries, locEntry)
			}
		}
		updatedLocHistory, _ := json.Marshal(updatedLocEntries)
		err = ioutil.WriteFile(filename, updatedLocHistory, 0644)
		if err != nil {
			log.Printf("Error updating location history file: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/text")
		w.Write([]byte("Bad request: Invalid order_id"))
	}
}
