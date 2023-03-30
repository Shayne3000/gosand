package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

var API_KEY string = "8150e50dcd5b50bb05c7a227bae36aaa"

// Represents the data we need returned by the WeatherAPI
type WeatherData struct {
	Name string `json:"name"` // JSON tag for unmarshalling Weather API response data directly into the struct
	Main struct {
		TempKelvin float64 `json:"temp"`
	} `json:"main"`
}

func main() {
	// Assigns a handler function to a url pattern/endpoint in the ServeMux
	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/weather/", weatherRequestHandler)

	http.ListenAndServe(":8081", nil)
}

// Handler uses the http.ResponseWriter to write a response to the Client.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// Handles any HTTP request that comes to the Weather endpoint
func weatherRequestHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := queryWeather(city)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Encoder is used for serialising/marshalling json responses
	// from the WeatherData struct for return to the client
	json.NewEncoder(w).Encode(data)
}

// Function with a standard error handling idiom
func queryWeather(city string) (WeatherData, error) {
	var d WeatherData
	weatherResponse, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + API_KEY + "&q=" + city)

	if err != nil {
		return WeatherData{}, err
	}

	// closing the response body after we exit the function scope with
	// "defer" is an elegant form of resource management.
	defer weatherResponse.Body.Close()

	// Decoder is used for deserialising/unmarshalling the json response
	// directly into the WeatherData struct variable, d.
	if err := json.NewDecoder(weatherResponse.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}

	return d, nil
}
