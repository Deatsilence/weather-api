package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// ! ---------CURRENT WEATHER DATA----------
type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type wind struct {
	Speed float64 `json:"speed"`
	Deg   float64 `json:"deg"`
}

type clouds struct {
	All int `json:"all"`
}

type sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

type weatherData struct {
	Coord      coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       wind      `json:"wind"`
	Clouds     clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

//! ---------FORECAST WEATHER DATA----------

type BaseForecestData struct {
	Cod     string            `json:"cod"`
	Message int               `json:"message"`
	Cnt     int               `json:"cnt"`
	List    []WeatherForecast `json:"list"`
	City    City              `json:"city"`
}

type WeatherForecast struct {
	Dt         int                `json:"dt"`
	Main       WeatherMain        `json:"main"`
	Weather    []WeatherCondition `json:"weather"`
	Clouds     WeatherClouds      `json:"clouds"`
	Wind       WeatherWind        `json:"wind"`
	Visibility int                `json:"visibility"`
	Pop        float64            `json:"pop"`
	Sys        WeatherSys         `json:"sys"`
	DtTxt      string             `json:"dt_txt"`
}

type WeatherMain struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
	Humidity  int     `json:"humidity"`
	TempKf    float64 `json:"temp_kf"`
}

type WeatherCondition struct {
	Id          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherClouds struct {
	All int `json:"all"`
}

type WeatherWind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

type WeatherSys struct {
	Pod string `json:"pod"`
}

type City struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Coord      CityCoord `json:"coord"`
	Country    string    `json:"country"`
	Population int       `json:"population"`
	Timezone   int       `json:"timezone"`
	Sunrise    int       `json:"sunrise"`
	Sunset     int       `json:"sunset"`
}

type CityCoord struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

//! ---------FORECAST WEATHER DATA----------

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData

	err = json.Unmarshal(bytes, &c)

	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from go weather! \n"))
}

func query(city string) (weatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")

	if err != nil {
		return weatherData{}, err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?" + "q=" + city + "&appid=" + apiConfig.OpenWeatherMapApiKey)

	if err != nil {
		return weatherData{}, err
	}
	defer resp.Body.Close()

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}

	return d, nil
}

func weather(w http.ResponseWriter, r *http.Request) {
	// city := strings.SplitN(r.URL.Path, "/", 3)[2]
	city := r.URL.Query().Get("name")
	data, err := query(city)
	data.convertKelvinToCelsius()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func (kelvinToCelsius *weatherData) convertKelvinToCelsius() {
	(*kelvinToCelsius).Main.Temp -= 273.15
}
func (kelvinToCelsius *BaseForecestData) ConvertKelvinToCelsius() {
	for i := range kelvinToCelsius.List {
		kelvinToCelsius.List[i].Main.Temp -= 273.15
		fmt.Println(kelvinToCelsius.List[i].Main.Temp)
	}
}

func forecastQuery(lat string, lon string) (BaseForecestData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")

	if err != nil {
		return BaseForecestData{}, err
	}

	resp, err := http.Get("https://api.openweathermap.org/data/2.5/forecast?lat=" + lat + "&lon=" + lon + "&cnt=40&appid=" + apiConfig.OpenWeatherMapApiKey)

	if err != nil {
		return BaseForecestData{}, err
	}
	defer resp.Body.Close()

	var d BaseForecestData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return BaseForecestData{}, err
	}

	return d, nil
}

func forecast(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	data, err := forecastQuery(lat, lon)
	data.ConvertKelvinToCelsius()
	fmt.Println(data.List[0].Main.Temp)
	fmt.Println(data.List[1].Main.Temp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", weather)
	http.HandleFunc("/forecast/", forecast)

	fmt.Println("Server is listening on port 8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)

	}
}
