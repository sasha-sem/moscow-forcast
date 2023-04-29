package forecast

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocarina/gocsv"
	"github.com/sasha-sem/moscow-forcast/cugms/internal/parser"
)

const (
	forecast_url = "http://cugms.ru/pogoda-i-klimat/prognoz-pogody/"
)

func NewParser(filePath string) parser.Parser {
	return &forecastParser{
		lastForecast:    "",
		currentForecast: []*forecast{},
		filePath:        filePath,
	}
}

type forecastParser struct {
	lastForecast    string
	currentForecast []*forecast
	filePath        string
}

func (p *forecastParser) Parse() error {

	forecasts, err := p.getForecast()
	if err != nil {
		return err
	}

	p.currentForecast = forecasts

	return nil
}

func (p *forecastParser) getForecast() ([]*forecast, error) {
	resp, err := http.Get(forecast_url)
	if err != nil {
		return nil, fmt.Errorf("couldn't get forecast from CUGMS: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read body of forecast from CUGMS: %w", err)
	}

	forecasts := []*forecast{}
	updatedTime := doc.Find("figure.wp-block-table + p").Find("strong").Text()
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		if i == 0 {
			return
		}

		forecast := &forecast{}

		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			block := td.Text()
			switch j {
			case 0:
				{
					forecast.Date = block
				}
			case 1:
				{
					forecast.Сloudiness = strings.TrimSpace(block)
				}
			case 2:
				{
					forecast.TemperatureMoscow = strings.TrimSpace(block)
				}
			case 3:
				{
					forecast.TemperatureMoscowRegion = strings.TrimSpace(block)
				}
			case 4:
				{
					forecast.Wind = block
				}
			}
		})

		forecast.UpdatedTime = updatedTime
		forecasts = append(forecasts, forecast)
	})

	return forecasts, nil
}

func (p *forecastParser) Write() error {
	currentForecastBytes, err := json.Marshal(p.currentForecast)
	if err != nil {
		return fmt.Errorf("couldn't marshal current weather: %w", err)
	}

	currentForecastStr := string(currentForecastBytes)
	if currentForecastStr == p.lastForecast {
		return nil
	}

	exists, err := fileExists(p.filePath)
	if err != nil {
		return fmt.Errorf("couldn't check existance of forecast file \"%s\": %w", p.filePath, err)
	}

	if !exists {
		err = p.writeNewFile()
	} else {
		err = p.writeAppend()
	}
	if err != nil {
		return err
	}

	p.lastForecast = currentForecastStr
	return nil
}

func (p *forecastParser) writeNewFile() error {
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

	err = gocsv.MarshalFile(p.currentForecast, file)
	if err != nil {
		return fmt.Errorf("couldn't write to new weather file \"%s\": %w", p.filePath, err)
	}

	return nil
}

func (p *forecastParser) writeAppend() error {
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

	err = gocsv.MarshalCSVWithoutHeaders(p.currentForecast, writer)
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
