package main

import (
	"aMazeGardenAI/db"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Json received from aMazeGarden data logger template.
type MeteoData struct {
	AirTemperature string `json:"airTemperature"`
	AirHumidity string `json:"airHumidity"`
	SoilMoisture string `json:"soilMoisture"`
	SoilTemperature string `json:"soilTemperature"`
}


// Database handler object.
var dbHandler = db.DatabaseHandler{
	DriverName:	"mysql",
	User:		"aMazeGardenAI",
	Password:	"Tvlwxfcg+1q",
	Database:	"plants",
}

// Url to send data from data logger.
var meteoDataUrl = "/dataLoggerData"

var maxSoilTemp = 25
var useForecast = true

var meteoData MeteoData

func main() {
	// Handle function for receiving data from data logger.
	http.HandleFunc("/", handle)
	http.HandleFunc("/addPlant", addPlantForm)
	http.HandleFunc(meteoDataUrl, handleMeteoDataRequest)

	fmt.Println("Http server listening...")

	// Connecting todatabase.
	err := dbHandler.Begin()

	// Error handling.
	if err != nil {
		panic(err)
	}
	// Listen on port 8080.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func wateringAI(meteoData MeteoData) () {
	fmt.Println("Got data: ", meteoData)

	var numberOfPlants int
	var err error

	// Get number of plants in database.
	if numberOfPlants, err = dbHandler.GetNumberOfPlants(); err != nil {
		fmt.Println("Error getting plants number.")
		panic(err)
	}

	// Scan for all plants.
	for  i := 1; i < numberOfPlants; i++ {
		// Convert string to int.
		soilMoisture, err := strconv.Atoi(meteoData.SoilMoisture)

		if err != nil {
			fmt.Println("Error getting soil moisture.")
			panic(err)
		}

		soilTemperature, err := strconv.Atoi(meteoData.SoilTemperature)

		if err != nil {
			fmt.Println("Error getting soil temperature.")
			panic(err)
		}


		// Get watered soil moisture of plant from the database.
		wateredSoilMoisture, err := dbHandler.GetWateredSoilMoistureFromId(i)

		// Error handling.
		if err != nil {
			fmt.Println("Error getting watered soil moisture.")
			panic(err)
		}

		// If soil is dry.
		if soilMoisture < wateredSoilMoisture {
			if soilTemperature < maxSoilTemp {
				fmt.Println("Watering needed")
			}
		}
	}
}

func addPlantForm(writer http.ResponseWriter, request *http.Request) {
	form, err := ioutil.ReadFile("addPlantForm.html")

	if err != nil {
		panic(err)
	}

	if _, err := writer.Write(form); err != nil {
		panic(err)
	}
}

func handle(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("index.html")

	if err != nil {
		panic(err)
	}

	if err := t.Execute(writer, meteoData); err != nil {
		panic(err)
	}
}

/*
	Handling data logger request.
*/
func handleMeteoDataRequest(writer http.ResponseWriter, request *http.Request) {
	// If url is bad.
	if request.URL.Path != meteoDataUrl {
		http.Error(writer, "404 not found.", http.StatusNotFound)
		return
	}

	// If request method is POST.
	if request.Method == "POST" {
		// Read request.
		jsn, err := ioutil.ReadAll(request.Body)

		// Parsing request to string.
		requestMessage := string(jsn)

		// Error handling.
		if err != nil {
			panic(err)
		}

		// Decoding json.
		if err := json.Unmarshal([]byte(requestMessage), &meteoData); err != nil {
			panic(err)
		}

		// Calling AI fuction.
		wateringAI(meteoData)
	}
}


