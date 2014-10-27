package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const homepageHTML = `
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
            var blah;
            $(document).ready(function() {
                $.getJSON("svc/init", function(data) {
                    var map = L.map('map').setView(
                        [data.Lat, data.Long], data.DefaultZoom
                    );
                    L.tileLayer(data.UrlTemplate, data.Options).addTo(map);
                });

                $.getJSON("svc/drinks", function(data) {
                    for (i = 0; i < data.results; i++) {
                        console.log(data.results[i]);
                    }
                });
           });
        </script>
    </head>

</html>
`

const apiFmt = `https://api.enigma.io/v2/data/%s/enigma.licenses.liquor.us?search[]=@establishment_address_zip+("11222")&conjunction=and`

type initValues struct {
	Lat         float32
	Long        float32
	UrlTemplate string
	MapHeight   string
	DefaultZoom int
	Options     struct {
		Attribution string
		MaxZoom     int
	}
}

/*
type liquorRow struct {
	Address string `json:"establishment_address_street1"`
}

type engimaResults struct {
	Result []liquorRow
}
*/

func initVals(w http.ResponseWriter, r *http.Request) {
	i := initValues{}
	i.Lat = 40.7263
	i.Long = -73.9456
	i.DefaultZoom = 15
	i.UrlTemplate = `http://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png`
	//i.UrlTemplate = `http://otile{s}.mqcdn.com/tiles/1.0.0/osm/{z}/{x}/{y}.png`
	//i.UrlTemplate = `http://{s}.tile.cloudmade.com/{key}/{styleId}/256/{z}/{x}/{y}.png`

	i.Options.Attribution = `© <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors`
	i.Options.MaxZoom = 18

	b, err := json.Marshal(i)

	if err != nil {
		log.Println(err)
		return
	}

	w.Write(b)
}

func drinks(w http.ResponseWriter, r *http.Request) {
	apiURL := fmt.Sprintf(apiFmt, os.Getenv("ENIGMA_API_KEY"))

	resp, err := http.Get(apiURL)

	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(apiURL)
	log.Printf("%s", b)

	/*
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Println(err)
		}
	*/
}

func homepage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, homepageHTML)
}

func main() {
	http.HandleFunc("/svc/init", initVals)
	http.HandleFunc("/svc/drinks", drinks)
	http.HandleFunc("/", homepage)

	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic(err)
	}
}
