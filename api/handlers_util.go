package api

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// TODO: move to a config file
var filename = filepath.Join("../", "location-data", "location_data.json")

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

// This function returns location history and previous location order ids
func getPreviousLocationData() ([]locationEntry, []string, error) {
	// get previous location history
	prevLocData, err := readLocHistoryFile(filename)
	var prevLocEntries []locationEntry
	err = json.Unmarshal(prevLocData, &prevLocEntries)

	var prevOrderIds []string
	for _, loc := range prevLocEntries {
		prevOrderIds = append(prevOrderIds, loc.OrderId)
	}
	return prevLocEntries, prevOrderIds, err
}

// This function checks whether a location history file is present
// if not, creates a new location history file
func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if err := os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return err
	}
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// This function reads the location history json file and returns its contents
func readLocHistoryFile(filename string) ([]byte, error) {
	err := checkFile(filename)
	file, err := ioutil.ReadFile(filename)
	return file, err
}

// This function checks whether the order_id from an incoming request is present in location history
// (used in GET and DELETE operations)
func isValidOrderId(orderId string, prevOrderIds []string) bool {
	for _, id := range prevOrderIds {
		if orderId == id {
			return true
		}
	}
	return false
}
