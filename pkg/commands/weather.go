package commands

import (
	"fmt"
	"strconv"
	"time"

	owm "github.com/briandowns/openweathermap"
)

type Weather struct {
	apiKey   string
	current  *owm.CurrentWeatherData
	forecast *owm.ForecastWeatherData
}

func NewWeatherClient(apiKey string) (weather Weather, err error) {
	weather.apiKey = apiKey
	c, err := owm.NewCurrent("C", "en", apiKey)
	if err != nil {
		return
	}
	f, err := owm.NewForecast("16", "C", "en", apiKey)
	if err != nil {
		return
	}
	weather.current = c
	weather.forecast = f
	return
}

func (weather Weather) GetWeather(location string) (string, error) {
	// weather.current.CurrentByName(location)
	weather.forecast.DailyByName(location, 5)

	message := fmt.Sprintf(`### Weather in %s for the next few days

| Day | Description | High | Low | Humidity | Day |
|:---------------------------|:------------------------------------|:--------|:--------|:--------|:--------|`, location)

	forecast := weather.forecast.ForecastWeatherJson.(*owm.Forecast16WeatherData)
	for _, day := range forecast.List {
		text := getWeatherLine(day)
		message += text
	}
	return message, nil
}

// Return a markdown formatted line for a weather day
func getWeatherLine(day owm.Forecast16WeatherList) string {

	time := time.Unix(int64(day.Dt), 0)

	min := strconv.FormatFloat(day.Temp.Min, 'f', 0, 64)
	max := strconv.FormatFloat(day.Temp.Max, 'f', 0, 64)
	dayT := strconv.FormatFloat(day.Temp.Day, 'f', 0, 64)

	return fmt.Sprintf(`
| %s, %s. %d | %s | %s °C | %s °C | %d%% | %s °C |`,
		time.Weekday().String(),
		time.Month().String(),
		time.Day(),
		getDescriptionWithIcon(day.Weather[0].Icon, day.Weather[0].Description),
		max,
		min,
		day.Humidity,
		dayT)
}

func getDescriptionWithIcon(iconCode string, description string) string {
	return fmt.Sprintf("%s ![%s](http://openweathermap.org/img/w/%s.png \"%s\")", description, description, iconCode, description)
}
