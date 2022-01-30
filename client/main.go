package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	values := map[string]float64{"lat": 12.24, "lng": 56.78}
	reqBody, _ := json.Marshal(values)

	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", "http://localhost:8899/location/abc123?max=2", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	fmt.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", string(body))
}
