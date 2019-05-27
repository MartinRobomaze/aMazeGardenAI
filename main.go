package main

import (
	"aMazeGardenAI/db"
	"aMazeGardenAI/serverUtils"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Json received from aMazeGarden data logger template.
type MeteoData struct {
	AppID          string `json:"app_id"`
	DevID          string `json:"dev_id"`
	HardwareSerial string `json:"hardware_serial"`
	Port           int    `json:"port"`
	Counter        int    `json:"counter"`
	PayloadRaw     string `json:"payload_raw"`
	PayloadFields  struct {
		AirHumidity     int `json:"airHumidity"`
		AirTemperature  int `json:"airTemperature"`
		SoilMoisture    int `json:"soilMoisture"`
		SoilTemperature int `json:"soilTemperature"`
	} `json:"payload_fields"`
	Metadata struct {
		Time       time.Time `json:"time"`
		Frequency  float64   `json:"frequency"`
		Modulation string    `json:"modulation"`
		DataRate   string    `json:"data_rate"`
		CodingRate string    `json:"coding_rate"`
		Gateways   []struct {
			GtwID     string    `json:"gtw_id"`
			Timestamp int       `json:"timestamp"`
			Time      time.Time `json:"time"`
			Channel   int       `json:"channel"`
			Rssi      int       `json:"rssi"`
			Snr       float64   `json:"snr"`
			RfChain   int       `json:"rf_chain"`
			Latitude  float64   `json:"latitude"`
			Longitude float64   `json:"longitude"`
			Altitude  int       `json:"altitude,omitempty"`
		} `json:"gateways"`
	} `json:"metadata"`
	DownlinkURL string `json:"downlink_url"`
}

type ForecastData struct {
	ForecastTemperature              string `json:"forecastTemperature"`
	ForecastPrecipitationPossibility string `json:"forecastPrecipitationPossibility"`
}

type DisplayData struct {
	AirTemperature                   string
	AirHumidity                      string
	SoilMoisture                     string
	SoilTemperature                  string
	ForecastTemperature              string
	ForecastPrecipitationPossibility string
}

type PayloadData struct {
	DevID         string `json:"dev_id"`
	PayloadFields PayloadFields `json:"payload_fields"`
}

type PayloadFields struct {
	SoilMoisture int `json:"soilMoisture"`
	PosX         int `json:"posX"`
	PosY         int `json:"posY"`
}

// Database handler object.
var dbHandler = db.DatabaseHandler{
	DriverName: "mysql",
	User:       "sql7293414",
	Password:   "JQQu7m6z8m",
	Database:   "sql7293414",
}

var maxSoilTemp = 25

var meteoData MeteoData

var forecastData ForecastData

var addPlantFormLoader = serverUtils.FileLoader{Path: "html/addPlantForm.html"}
var setGardenFormLoader = serverUtils.FileLoader{Path: "html/gardenSettingsForm.html"}

var editPlantFormLoader = serverUtils.FileTemplateLoader{
	Path:      "html/editPlantForm.html",
	DbHandler: dbHandler,
}

var deletePlantFormLoader = serverUtils.FileTemplateLoader{
	Path:      "html/removePlantForm.html",
	DbHandler: dbHandler,
}

func main() {
	port := os.Getenv("PORT")
	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", handle)
	http.HandleFunc("/addPlant", addPlantFormLoader.LoadFile)
	http.HandleFunc("/removePlant", deletePlantFormLoader.LoadFileTemplate)
	http.HandleFunc("/editPlant", editPlantFormLoader.LoadFileTemplate)
	http.HandleFunc("/gardenSetting", setGardenFormLoader.LoadFile)
	http.HandleFunc("/addPlantDb", addPlantToDb)
	http.HandleFunc("/removePlantDb", deletePlantDb)
	http.HandleFunc("/editPlantDb", editPlantDb)
	http.HandleFunc("/dataLoggerData", handleMeteoDataRequest)
	http.HandleFunc("/forecastData", handleForecastRequest)

	fmt.Println("Http server listening...")

	// Connecting todatabase.
	err := dbHandler.Begin()

	// Error handling.
	if err != nil {
		panic(err)
	}
	// Listen on port 8080.
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		panic(err)
	}
}

func wateringAI() {
	fmt.Println("Got data: ", meteoData)
	plantsWateredSoilMoisture, err := dbHandler.GetAllPlantsSoilMoisture()
	fmt.Println(plantsWateredSoilMoisture)

	// Get number of plants in database.
	if err != nil {
		panic(err)
	}

	// Scan for all plants.
	for i := 0; i < len(plantsWateredSoilMoisture); i++ {
		soilTemperature := meteoData.PayloadFields.SoilTemperature

		// Get watered soil moisture of plant from the database.
		wateredSoilMoisture, err := strconv.Atoi(plantsWateredSoilMoisture[i])

		// Error handling.
		if err != nil {
			panic(err)
		}

		// If soil is dry.
		if wateredSoilMoisture < wateredSoilMoisture {
			if soilTemperature < maxSoilTemp {
				plantID, err := dbHandler.GetPlantID(wateredSoilMoisture)

				plantX, err := dbHandler.GetPlantX(plantID);
				plantY, err := dbHandler.GetPlantX(plantID);

				forecastPrecipitation, err := strconv.Atoi(forecastData.ForecastPrecipitationPossibility)

				if err != nil {
					water(wateredSoilMoisture, plantX, plantY)
					break
				}

				if forecastPrecipitation < 70 {
					water(wateredSoilMoisture, plantX, plantY)
					break
				}
			}
		}
	}
}

func water(soilMoisture int, posX int, posY int) {
	fmt.Println("Watering needed")

	var payloadData = PayloadData{
		DevID:"meteo_data_logger",
		PayloadFields:PayloadFields{
			SoilMoisture:soilMoisture,
			PosX:posX,
			PosY:posY,
		},
	}

	payload, err := json.Marshal(payloadData)

	if err != nil {
		panic(err)
	}

	fmt.Println(payloadData)

	_, err = http.Post("https://en236k5hsg4bxj.x.pipedream.net", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		panic(err)
	}
}

func handle(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("html/index.html")

	if err != nil {
		panic(err)
	}
	displayData := DisplayData{
		AirTemperature:                   strconv.Itoa(meteoData.PayloadFields.AirTemperature),
		AirHumidity:                      strconv.Itoa(meteoData.PayloadFields.AirHumidity),
		SoilMoisture:                     strconv.Itoa(meteoData.PayloadFields.SoilMoisture),
		SoilTemperature:                  strconv.Itoa(meteoData.PayloadFields.SoilTemperature),
		ForecastTemperature:              forecastData.ForecastTemperature,
		ForecastPrecipitationPossibility: forecastData.ForecastPrecipitationPossibility,
	}
	if err := t.Execute(writer, displayData); err != nil {
		panic(err)
	}
}

func addPlantToDb(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			panic(err)
		}

		plantName := request.FormValue("plantName")
		plantWateredSoilMoisture, err := strconv.Atoi(request.FormValue("wateredSoilMoisture"))
		plantPositionX, err := strconv.Atoi(request.FormValue("positionX"))
		plantPositionY, err := strconv.Atoi(request.FormValue("positionY"))

		if err != nil {
			http.Redirect(writer, request, "/addPlant", 303)
		}

		err = dbHandler.Write(plantName, plantWateredSoilMoisture, plantPositionX, plantPositionY)

		if err != nil {
			panic(err)
		}

		http.Redirect(writer, request, "/", 303)
	}
}

func editPlantDb(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			panic(err)
		}

		plantName := request.FormValue("plantsList")

		wateredSoilMoisture, err := strconv.Atoi(request.FormValue("wateredSoilMoisture"))

		if err != nil {
			panic(err)
		}

		err = dbHandler.Update(plantName, wateredSoilMoisture)

		if err != nil {
			http.Redirect(writer, request, "/editPlant", 303)
		}

		http.Redirect(writer, request, "/", 303)
	}
}

func deletePlantDb(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		if err := request.ParseForm(); err != nil {
			panic(err)
		}

		plantName := request.FormValue("plantsList")

		err := dbHandler.DeletePlant(plantName)

		if err != nil {
			http.Redirect(writer, request, "/removePlant", 303)
		}

		http.Redirect(writer, request, "/", 303)
	}
}

func handleForecastRequest(writer http.ResponseWriter, request *http.Request) {
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
		if err := json.Unmarshal([]byte(requestMessage), &forecastData); err != nil {
			panic(err)
		}
	}

	wateringAI()
}

/*
	Handling data logger request.
*/
func handleMeteoDataRequest(writer http.ResponseWriter, request *http.Request) {
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
		wateringAI()
	}
}
