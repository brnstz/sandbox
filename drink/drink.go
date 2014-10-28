package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	homepageHTML = `
<html>
    <head>
        <link rel="stylesheet" href="http://cdn.leafletjs.com/leaflet-0.7.3/leaflet.css" />
        <script src="http://cdn.leafletjs.com/leaflet-0.7.3/leaflet.js"></script>
        <script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>

        <style>
            #map { height: 700px; }
        </style>

        <body>
            <div id="map"></div>
        </body>

        <script>
            var map;
            $(document).ready(function() {
                $.getJSON("svc/init", function(data) {
                    map = L.map('map').setView(
                        [data.Lat, data.Lon], data.DefaultZoom
                    );
                    L.tileLayer(data.UrlTemplate, data.Options).addTo(map);
                });

                $.getJSON("svc/drinks", function(data) {
                    for (i = 0; i < data.length; i++) {
                        var marker = L.marker([data[i].lat, data[i].lon]).addTo(map);
                        marker.bindPopup(data[i].licensee_name_business);
                    }
                });
           });
        </script>
    </head>

</html>
`
)

var (
	nomURL       = `http://nominatim.openstreetmap.org/search`
	liquorSearch = `?search[]=@establishment_address_zip+("11222")&conjunction=and`

	enigmaDataset = `enigma.licenses.liquor.us`
	enigmaAPIKey  = os.Getenv("ENGIMA_API_KEY")
	enigmaURL     = fmt.Sprint(
		`https://api.enigma.io/v2/data/`,
		os.Getenv("ENIGMA_API_KEY"),
		`/`, enigmaDataset,
		liquorSearch,
	)

	// cache licenses after loading once
	cachedLicenses []*license
)

// initValues is returned to init the map
type initValues struct {
	Lat         float64
	Lon         float64
	UrlTemplate string
	MapHeight   string
	DefaultZoom int
	Options     struct {
		Attribution string
		MaxZoom     int
	}
}

// enigmaResponse is what enigmaURL returns
type enigmaResponse struct {
	Result []*license
}

// license is a single instance of license from enigma.licenses.liquor.us
type license struct {
	// Data from engima
	EstablishmentName string `json:"licensee_name_business"`
	Address           string `json:"establishment_address_street1"`
	City              string `json:"establishment_address_city"`
	State             string `json:"establishment_address_state"`
	Zip               string `json:"establishment_address_zip"`

	// appended data from nominatim
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// nomResp is the response from http://nominatim.openstreetmap.org/
type nomResp []struct {
	Lat string
	Lon string
}

// getJSON GETs a URL and unmarshals into obj, or returns an error
func getJSON(url string, obj interface{}) (err error) {
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		log.Println(err)
		return
	}

	// Success!
	return
}

func appendGeo(lic *license) (err error) {
	v := url.Values{}
	v.Set("format", "json")
	v.Set("street", lic.Address)
	v.Set("city", lic.City)
	v.Set("state", lic.State)

	nr := nomResp{}
	err = getJSON(fmt.Sprint(nomURL, "?", v.Encode()), &nr)
	if err != nil {
		log.Println(err)
		return
	}
	if len(nr) < 1 {
		log.Println("No response from ", nomURL)

		return
	}

	lic.Lat, err = strconv.ParseFloat(nr[0].Lat, 64)
	if err != nil {
		log.Println("Can't convert latitude", err)
		return
	}

	lic.Lon, err = strconv.ParseFloat(nr[0].Lon, 64)
	if err != nil {
		log.Println("Can't convert longitude", err)
		return
	}

	return
}

// loadDrinks loads places to drink in the background
func loadDrinks() {
	log.Println("starting to load drinks")

	// Get the liceneses from Engima
	eResp := enigmaResponse{}
	err := getJSON(enigmaURL, &eResp)
	if err != nil {
		log.Println(err)
		return
	}

	// Append geo and yelp info
	for _, lic := range eResp.Result {
		err = appendGeo(lic)
		if err != nil {
			continue
		}

		log.Printf("%+v", lic)
		cachedLicenses = append(cachedLicenses, lic)
	}
	log.Println("done loading")
}

// initVals is a handler that returns JSON to initialize the map
func initVals(w http.ResponseWriter, r *http.Request) {
	i := initValues{}
	i.Lat = 40.7263
	i.Lon = -73.9456
	i.DefaultZoom = 15
	i.UrlTemplate = `http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png`

	i.Options.Attribution = `© <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors`
	i.Options.MaxZoom = 18

	b, err := json.Marshal(i)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}

// drinks is a handler that returns data on the drinking establishments
// in the neighborhood.
func drinks(w http.ResponseWriter, r *http.Request) {

	b, err := json.Marshal(cachedLicenses)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
	}
}

// homepage is a handler that renders the homepage
func homepage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, homepageHTML)
}

func main() {
	go loadDrinks()

	http.HandleFunc("/svc/init", initVals)
	http.HandleFunc("/svc/drinks", drinks)
	http.HandleFunc("/", homepage)

	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic(err)
	}
}
