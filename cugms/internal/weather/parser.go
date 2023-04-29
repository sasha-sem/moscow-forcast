package weather

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/sasha-sem/moscow-forcast/cugms/internal/parser"
)

const (
	current_weather_url = "http://cugms.ru/wp-content/uploads/2022/01/Yandex/meteocsdn.csv"
)

func NewParser(filePath string) parser.Parser {
	return &weatherParser{
		lastWeather:    "",
		currentWeather: []*weather{},
		filePath:       filePath,
	}
}

type weatherParser struct {
	lastWeather    string
	currentWeather []*weather
	filePath       string
}

func (p *weatherParser) Parse() error {
	weather, err := p.getWeather()
	if err != nil {
		return err
	}

	p.currentWeather = weather

	return nil
}

func (p *weatherParser) getWeather() ([]*weather, error) {
	resp, err := http.Get(current_weather_url)
	if err != nil {
		return nil, fmt.Errorf("couldn't get current weather from CUGMS: %w", err)
	}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'
		return r
	})

	weather := []*weather{}

	if err := gocsv.Unmarshal(resp.Body, &weather); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal body of current weather from CUGMS: %w", err)
	}

	for _, w := range weather {
		w.LastUpdateTime = time.Now().Format(time.DateOnly) + " " + w.LastUpdateTime
	}

	return weather, nil
}

func (p *weatherParser) Write() error {
	currentWeatherBytes, err := json.Marshal(p.currentWeather)
	if err != nil {
		return fmt.Errorf("couldn't marshal current weather: %w", err)
	}

	currentWeatherStr := string(currentWeatherBytes)
	if currentWeatherStr == p.lastWeather {
		return nil
	}

	exists, err := fileExists(p.filePath)
	if err != nil {
		return fmt.Errorf("couldn't open weather file \"%s\": %w", p.filePath, err)
	}

	if !exists {
		err = p.writeNewFile()
	} else {
		err = p.writeAppend()
	}
	if err != nil {
		return err
	}

	p.lastWeather = currentWeatherStr
	return nil
}

func (p *weatherParser) writeNewFile() error {
	file, err := os.OpenFile(p.filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("couldn't create new weather file \"%s\": %w", p.filePath, err)
	}
	defer file.Close()

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	err = gocsv.MarshalFile(p.currentWeather, file)
	if err != nil {
		return fmt.Errorf("couldn't write to new weather file \"%s\": %w", p.filePath, err)
	}

	return nil
}

func (p *weatherParser) writeAppend() error {
	file, err := os.OpenFile(p.filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open existing weather file \"%s\": %w", p.filePath, err)
	}
	defer file.Close()

	writer := gocsv.NewSafeCSVWriter(func() *csv.Writer {
		csvWriter := csv.NewWriter(file)
		csvWriter.Comma = ';'
		return csvWriter
	}())

	err = gocsv.MarshalCSVWithoutHeaders(p.currentWeather, writer)
	if err != nil {
		return fmt.Errorf("couldn't write to existing weather file \"%s\": %w", p.filePath, err)
	}

	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
