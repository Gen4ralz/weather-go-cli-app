package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location 	Location	`json:"location"`
	Current		CurrentTemp	`json:"current"`
	Forecast	Forecast	`json:"forecast"`
}

type Location struct {
	Name	string	`json:"name"`
	Country	string	`json:"country"`
}

type CurrentTemp struct {
	TempC		float64		`json:"temp_c"`
	Condition	Condition	`json:"condition"`
}

type Condition struct {
	Text	string	`json:"text"`
}

type Forecast struct {
	Forecastday []Forecastday	`json:"forecastday"`
}

type Forecastday struct {
	Hour	[]Hour	`json:"hour"`
}

type Hour struct {
	TimeEpoch		int64 		`json:"time_epoch"`
	TempC			float64		`json:"temp_c"`
	Condition		Condition	`json:"condition"`
	ChanceOfRain	float64		`json:"chance_of_rain"`
}

func main() {
	q := "Bangkok"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}
	
	response, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=2ea4d1ddc816428492a95618231608&q="+ q +"&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather

	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf("%s - %.0fC, %.0f, %s\n", date.Format("15.04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}