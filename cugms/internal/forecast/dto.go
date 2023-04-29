package forecast

type forecast struct {
	Date                    string `csv:"Дата"`
	Сloudiness              string `csv:"Облачность. Осадки. явления."`
	TemperatureMoscow       string `csv:"Температура. °С. г. Москва"`
	TemperatureMoscowRegion string `csv:"Температура. °С. Московская область"`
	Wind                    string `csv:"Ветер. м/с."`
	UpdatedTime             string `csv:"Дата обновления"`
}
