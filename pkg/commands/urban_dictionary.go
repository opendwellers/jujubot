package commands

import (
	ud "github.com/dpatrie/urbandictionary"
)

func GetUrbanDictionaryDefinition(term string) (def *ud.Result, err error) {
	res, err := ud.Query(term)
	if err != nil {
		return
	}

	def = &res.Results[0]
	return
}
