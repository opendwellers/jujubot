package commands

import (
	"errors"

	ud "github.com/dpatrie/urbandictionary"
)

func GetUrbanDictionaryDefinition(term string) (def *ud.Result, err error) {
	res, err := ud.Query(term)
	if err != nil {
		return
	}
	if len(res.Results) == 0 {
		err = errors.New("No results found")
		return
	}

	def = &res.Results[0]
	return
}
