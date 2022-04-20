package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type JishoResponse struct {
	Data []Datum `json:"data"`
}

type Datum struct {
	Slug     string     `json:"slug"`
	IsCommon bool       `json:"is_common"`
	Tags     []string   `json:"tags"`
	Japanese []Japanese `json:"japanese"`
	Senses   []Sense    `json:"senses"`
}

type Japanese struct {
	Word    string `json:"word"`
	Reading string `json:"reading"`
}

type Sense struct {
	EnglishDefinitions []string      `json:"english_definitions"`
	PartsOfSpeech      []string      `json:"parts_of_speech"`
	Links              []Link        `json:"links"`
	Tags               []string      `json:"tags"`
	Restrictions       []interface{} `json:"restrictions"`
	SeeAlso            []string      `json:"see_also"`
	Antonyms           []string      `json:"antonyms"`
	Source             []interface{} `json:"source"`
	Info               []string      `json:"info"`
	Sentences          []interface{} `json:"sentences"`
}

type Link struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func GetWotdJapanese() (string, error) {
	// Get a random generator that stays the same for a given day
	randomGenerator := rand.New(rand.NewSource(int64(time.Now().YearDay())))
	wotdUrl := fmt.Sprintf("https://jisho.org/api/v1/search/words?keyword=%%23common&page=%d", randomGenerator.Intn(29)+1)
	resp, err := http.Get(wotdUrl)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()

	var jishoResp JishoResponse
	err = json.NewDecoder(resp.Body).Decode(&jishoResp)
	if err != nil {
		return "", err
	}
	random := jishoResp.Data[randomGenerator.Intn(len(jishoResp.Data))]

	message := fmt.Sprintf(`
#### Japanese word of the day for %s

# **%s**
*%s*
Meanings:`, time.Now().Format("Monday, January 2, 2006"),
		random.Japanese[0].Word,
		random.Japanese[0].Reading)

	for _, sense := range random.Senses {
		message += "\n- "
		for _, definition := range sense.EnglishDefinitions {
			message += fmt.Sprintf("%s, ", definition)
		}

		for _, link := range sense.Links {
			message += "\n"
			message += fmt.Sprintf("  - [%s](%s)", link.Text, link.URL)
		}
	}
	return message, nil
}
