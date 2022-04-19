package commands

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const baseUrl = "https://frankfurter.app/latest"

type apiResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func Convert(from, to string, amount float64) (float64, error) {
	to = strings.ToUpper(to)
	from = strings.ToUpper(from)
	url := baseUrl + "?from=" + from + "&to=" + to + "&amount=" + strconv.FormatFloat(amount, 'f', 2, 64)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("Invalid response code from convert service: " + resp.Status)
	}

	var response apiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}
	return response.Rates[to], err
}
