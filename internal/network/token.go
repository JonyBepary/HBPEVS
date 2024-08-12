package network

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func parsePort(port uint64, seed string, start string, duration string) (string, error) {

	baseURL := fmt.Sprintf("http://localhost:%d/", port)
	// seed := "def95e06826ffb028c97aa85096078e44c488e79e405626f2858a8b070c761c1"
	// start := "1681083456"
	// duration := "10800000000000"

	url := baseURL + "?seed=" + seed + "&start=" + start + "&duration=" + duration
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	var response struct {
		Digest string `json:"digest"`
	}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {

		fmt.Println(err)
		return "", err
	}

	return response.Digest, nil
}
