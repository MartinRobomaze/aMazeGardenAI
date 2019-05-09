package ttn

import (
	"io/ioutil"
	"net/http"
)

//type MeteoData struct {
//	AirTemperature  string `json:"airTemperature"`
//	AirHumidity     string `json:"airHumidity"`
//	SoilMoisture    string `json:"soilMoisture"`
//	SoilTemperature string `json:"soilTemperature"`
//}

type HttpRequestHandler struct {
	Url  string
	Port string
}

var requestMessage string

func (handler HttpRequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != handler.Url {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "POST" {
		jsn, err := ioutil.ReadAll(r.Body);
		requestMessage = string(jsn);

		if err != nil {
			return
		}

		//err = json.Unmarshal(jsn, &meteoData)
		//
		//if (err != nil) {
		//	log.Fatal("Error parsing message ", err)
		//}
		//
		//log.Printf("Data received: %s", meteoData)
	}
}

func (handler HttpRequestHandler) Begin() {
	http.HandleFunc(handler.Url, handler.HandleRequest)
}