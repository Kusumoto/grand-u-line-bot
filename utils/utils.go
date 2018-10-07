package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// GetJSON for get all json data and mapping to object
func GetJSON(url string, target interface{}) error {
	var myClient = &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
