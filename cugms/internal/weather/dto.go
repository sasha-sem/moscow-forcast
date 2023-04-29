package weather

type weather struct {
	StationID      string `csv:"№ п/п"`
	Lat            string `csv:"Широта"`
	Lon            string `csv:"Долгота"`
	Index          string `csv:"Индекс"`
	Name           string `csv:"Наименование"`
	Temperature    string `csv:"Температура"`
	Humidity       string `csv:"Влажность"`
	Pressure       string `csv:"Давление"`
	WindSpeed      string `csv:"Скорость ветра"`
	WindDirection  string `csv:"Направление ветра"`
	LastUpdateTime string `csv:"Последнее обновление"`
}
