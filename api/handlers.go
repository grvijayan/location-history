package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"location-history/util"
	"net/http"
	"os"
	"time"
)

// TODO: move to a config file
var filename = "location-data/location_data.json"

type locHistoryResp struct {
	OrderId string    `json:"order_id"`
	History []history `json:"history"`
}

type history struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type locationEntry struct {
	OrderId   string
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Timestamp time.Time
}

func AddLocationData(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)

	if orderId, exist := requestVars["order_id"]; exist {
		locationReq := util.AddLocationRequest{}
		_ = json.NewDecoder(r.Body).Decode(&locationReq)

		loc := locationEntry{
			OrderId:   orderId,
			Lat:       locationReq.Lat,
			Lng:       locationReq.Lng,
			Timestamp: time.Now(),
		}

		prevLocData, err := getPrevLocHistory(filename)

		var prevLocEntries []locationEntry

		err = json.Unmarshal(prevLocData, &prevLocEntries)
		if err != nil {
			fmt.Println(err)
		}

		prevLocEntries = append(prevLocEntries, loc)

		updatedLocHistory, _ := json.Marshal(prevLocEntries)
		fmt.Println(string(updatedLocHistory))

		err = ioutil.WriteFile(filename, updatedLocHistory, 0644)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("Bad request: Order Id not present")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func GetLocationData(w http.ResponseWriter, r *http.Request) {
	requestVars := mux.Vars(r)
	if orderId, exist := requestVars["order_id"]; exist {
		prevLocData, err := getPrevLocHistory(filename)
		var prevLocEntries []locationEntry
		err = json.Unmarshal(prevLocData, &prevLocEntries)
		if err != nil {
			fmt.Println(err)
		}

		// TODO: return entries based on max value
		//entriesToReturn, _ := strconv.Atoi(r.FormValue("max"))
		//if entriesToReturn == 0 {
		//	entriesToReturn = len(prevLocEntries)
		//}

		var locResp locHistoryResp
		locResp.OrderId = orderId
		var locHistory []history
		for _, loc := range prevLocEntries {
			if loc.OrderId == orderId {
				fmt.Println(loc)
				locHistory = append(locHistory, history{
					Lat: loc.Lat,
					Lng: loc.Lng,
				})
			}
		}

		locResp.History = locHistory

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(locResp)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Bad request: Order Id not present")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DeleteLocationData(w http.ResponseWriter, r *http.Request) {
	// to be implemented
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func getPrevLocHistory(filename string) ([]byte, error) {
	err := checkFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return file, err
}
