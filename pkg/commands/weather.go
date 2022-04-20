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

func (w Weather) GetCurrentWeather(location string) (message string, err error) {
	err = w.current.CurrentByName(location)
	if err != nil {
		return
	}

	temp := strconv.FormatFloat(w.current.Main.Temp, 'f', 0, 64)
	feel := strconv.FormatFloat(w.current.Main.FeelsLike, 'f', 0, 64)
	wind := strconv.FormatFloat(w.current.Wind.Speed, 'f', 1, 64)
	// dayT := strconv.FormatFloat(day.Temp.Day, 'f', 0, 64)

	message = fmt.Sprintf(`### Current weather in %s

| Description | Temperature | Feels Like | Humidity | Wind |
|:---------------------------|:------------------------------------|:--------|:--------|:--------|:--------|
| %s | %s °C | %s °C | %d%% | %s km/h |`,
		w.current.Name,
		getDescriptionWithIcon(w.current.Weather[0].Icon, w.current.Weather[0].Description),
		temp,
		feel,
		w.current.Main.Humidity,
		wind)

	return
}

func (w Weather) GetWeather(location string) (message string, err error) {
	w.forecast.DailyByName(location, 5)

	message = fmt.Sprintf(`### Weather in %s for the next few days

| Day | Description | High | Low | Humidity | Day |
|:---------------------------|:------------------------------------|:--------|:--------|:--------|:--------|`, location)

	forecast := w.forecast.ForecastWeatherJson.(*owm.Forecast16WeatherData)
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
