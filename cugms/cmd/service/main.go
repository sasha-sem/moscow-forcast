package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/sasha-sem/moscow-forcast/cugms/internal/forecast"
	"github.com/sasha-sem/moscow-forcast/cugms/internal/scheduler"
	"github.com/sasha-sem/moscow-forcast/cugms/internal/weather"
)

const (
	DATA_DIR         = "data"
	WEATHER_FILENAME = "cugms_weather.csv"
	FORCAST_FILENAME = "cugms_forecast.csv"
	LOC              = "Europe/Moscow"
)

var ExecuteTime = []string{
	"00:00:00",
	"03:00:00",
	"06:00:00",
	"09:00:00",
	"12:00:00",
	"15:00:00",
	"18:00:00",
	"21:00:00",
}

func main() {
	if _, err := os.Stat(DATA_DIR); os.IsNotExist(err) {
		err := os.Mkdir(DATA_DIR, os.ModePerm)
		if err != nil {
			log.Print(fmt.Errorf("couln't create directory for data: %w", err))
			return
		}
	}

	location, err := time.LoadLocation(LOC)
	if err != nil {
		log.Print(fmt.Errorf("couln't parse location: %w", err))
		return
	}
	scheduler := scheduler.NewScheduler(ExecuteTime, location)

	weatherParser := weather.NewParser(path.Join(DATA_DIR, WEATHER_FILENAME))
	forecastParser := forecast.NewParser(path.Join(DATA_DIR, FORCAST_FILENAME))

	for {
		log.Print("[INFO] Starting to update")
		err := weatherParser.Parse()
		if err != nil {
			log.Print(err)
		} else {
			err = weatherParser.Write()
			if err != nil {
				log.Print(err)
			}
		}

		err = forecastParser.Parse()
		if err != nil {
			log.Print(err)
		} else {
			err = forecastParser.Write()
			if err != nil {
				log.Print(err)
			}
		}

		duration, err := scheduler.GetTimeToWait()
		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("[INFO] Updated. Next update time: %s", time.Now().Add(duration).In(location).Format(time.DateTime))
		time.Sleep(duration)
	}

}
