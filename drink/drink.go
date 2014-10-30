package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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
            var colorIcon = L.Icon.extend({
                options: {
                    iconSize: [25, 41],
                    iconAnchor: [12, 41],
                    popupAnchor: [1, -34],
                    shadowSize: [41, 41] 
                }
            });

            var redIcon = new colorIcon({iconUrl: 'http://brnstz.com/images/marker-icon-red.png'});
            var greenIcon = new colorIcon({iconUrl: 'http://brnstz.com/images/marker-icon-green.png'});
            var yellowIcon = new colorIcon({iconUrl: 'http://brnstz.com/images/marker-icon-yellow.png'});
            var blueIcon = new colorIcon({iconUrl: 'http://brnstz.com/images/marker-icon.png'});
            var greyIcon = new colorIcon({iconUrl: 'http://brnstz.com/images/marker-icon-grey.png'});

            $(document).ready(function() {
                $.getJSON("svc/init", function(data) {
                    map = L.map('map').setView(
                        [data.Lat, data.Lon], data.DefaultZoom
                    );
                    L.tileLayer(data.UrlTemplate, data.Options).addTo(map);
                });

                $.getJSON("svc/drinks", function(data) {
                    for (var key in data) {
                        console.log(data[key]);
                        var curIcon = redIcon;
                        if (data[key].currentgrade == "A") {
                            curIcon = greenIcon;
                        } else if (data[key].currentgrade == "B") {
                            curIcon = blueIcon;
                        } else if (data[key].currentgrade == "C") {
                            curIcon = yellowIcon;
                        } else if (data[key].currentgrade == "") {
                            curIcon = greyIcon;
                        }
                             
                        var marker = L.marker([data[key].lat, data[key].lon], {icon: curIcon}).addTo(map);
                        marker.bindPopup("Licensee: " + data[key].licensee_name_business +
                            "<br>Business: " + data[key].dba +
                            "<br>Inspection Grade: " + data[key].currentgrade);
                    }
                });
           });
        </script>
    </head>

</html>
`

	// URL to search for address and retrieve lat, lon
	nomURL = `http://nominatim.openstreetmap.org/search`
)

var (
	// cache licenses after loading once, also allow to look up my
	// license.key()
	cachedLicenses = map[string]*license{}

	// AVE or ST at end of string or word, case insensitive, optionally
	// followed by .
	aveR = regexp.MustCompile(`(?i)\sAVE\.?(?:\s|$)`)
	stR  = regexp.MustCompile(`(?i)\sST\.?(?:\s|$)`)
)

// initValues is returned to init the map
type initValues struct {
	Lat         float64
	Lon         float64
	UrlTemplate string
	MapHeight   string
	DefaultZoom int
	Options     struct {
		MaxZoom int
	}
}

// licenseResponse is what calls to enigma.licenses.liquor.us return
type licenseResponse struct {
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

	// appended from restaurant inspection
	BusinessName   string     `json:"dba"`
	CurrentGrade   string     `json:"currentgrade"`
	InspectionDate *time.Time `json:"inspdate"`
}

// key returns a slightly normalized address to match the license to other
// datasets
func (l *license) key() string {
	// The license dataset has mixed case and "AVE" or "ST", while the
	// restaurant inspection dataset consistently has "AVENUE" and "STREET" in
	// all caps. Do some regexs to normalize. This should match
	// us.states.ny.cities.nyc.dohmh.restaurants.inspections building+street
	a := strings.ToUpper(l.Address)

	a = strings.TrimSpace(aveR.ReplaceAllString(a, " AVENUE "))
	a = strings.TrimSpace(stR.ReplaceAllString(a, " STREET "))

	return a
}

// inspectionResponse is what calls to
// us.states.ny.cities.nyc.dohmh.restaurants.inspections return
type inspectionResponse struct {
	Result []*inspection
}

// inspection is a single result from inspectionResponse
type inspection struct {
	BusinessName   string     `json:"dba"`
	Building       string     `json:"building"`
	Street         string     `json:"street"`
	CurrentGrade   string     `json:"currentgrade"`
	InspectionDate *time.Time `json:"inspdate"`
}

// nomResp is the response from http://nominatim.openstreetmap.org/
type nomResp []struct {
	Lat string
	Lon string
}

// enigmaSearchValue is a list of values for this key (column)
type enigmaSearchVal struct {
	key    string
	values []string
}

// QuotedValues returns each value in esv.values quoted for use in String()
func (esv *enigmaSearchVal) QuotedValues() []string {
	quoted := make([]string, len(esv.values))

	for i, _ := range esv.values {
		var b bytes.Buffer
		// Starting quote
		b.WriteRune('"')

		// Go through each character to catch quotes we must escape
		for _, char := range esv.values[i] {
			switch char {
			case '"':
				// Escape quotes with a \
				b.WriteString(`\"`)

			default:
				// Other characters just write verbatim
				b.WriteRune(char)
			}
		}

		// Ending quote
		b.WriteRune('"')

		quoted[i] = b.String()
	}

	return quoted
}

// String stringifies an enigmaSearchVal.
// Example: @cuisine_type ("Polish"|"Chinese")
func (esv *enigmaSearchVal) String() string {
	return fmt.Sprint("@", esv.key, " (",
		strings.Join(esv.QuotedValues(), "|"), ")",
	)
}

// getEnigmaURL returns a full URL to search dataset for the values passed.
// Lets API infer default conjunction of "and". FIXME: May break in edge cases.
func getEnigmaURL(dataset string, search []enigmaSearchVal, page, sort string) string {
	v := url.Values{}
	for _, s := range search {
		v.Add("search[]", s.String())
	}
	v.Set("page", page)

	if len(sort) > 0 {
		v.Set("sort", sort)
	}

	baseURL := fmt.Sprint(
		`https://api.enigma.io/v2/data/`, os.Getenv("ENIGMA_API_KEY"), `/`,
		dataset,
	)

	return fmt.Sprint(baseURL, "?", v.Encode())
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

// saveLic serializes access to cachedLicenses
func saveLic(licChan chan *license) {
	for lic := range licChan {
		log.Println(lic)
		cachedLicenses[lic.key()] = lic
	}
}

// processGeo receive licenses from a channel and retrieves the lat and lon
// from nomURL
func processGeo(geoChan, licChan chan *license, wg *sync.WaitGroup) {
	for lic := range geoChan {
		v := url.Values{}
		v.Set("format", "json")
		v.Set("street", lic.Address)
		v.Set("city", lic.City)
		v.Set("state", lic.State)

		nr := nomResp{}
		err := getJSON(fmt.Sprint(nomURL, "?", v.Encode()), &nr)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(nr) < 1 {
			log.Println("No response from ", nomURL)
			continue
		}

		lic.Lat, err = strconv.ParseFloat(nr[0].Lat, 64)
		if err != nil {
			log.Println("Can't convert latitude", err)
			continue
		}

		lic.Lon, err = strconv.ParseFloat(nr[0].Lon, 64)
		if err != nil {
			log.Println("Can't convert longitude", err)
			continue
		}

		log.Println("Got response from", nomURL)

		licChan <- lic
	}
	wg.Done()
}

// loadDrinks loads places to drink in the background
func loadDrinks() {
	log.Println("starting to load drinks")

	// Get all places in Greenpoint with a liquor license. First page (1)
	// should be enough.
	esv := enigmaSearchVal{
		key:    "establishment_address_zip",
		values: []string{"11222", "11211"},
	}
	enigmaURL := getEnigmaURL(
		`enigma.licenses.liquor.us`, []enigmaSearchVal{esv}, `1`, ``,
	)

	// Get the licenses from Engima
	eResp := licenseResponse{}
	err := getJSON(enigmaURL, &eResp)
	if err != nil {
		return
	}

	geoChan := make(chan *license, 5000)
	licChan := make(chan *license, 5000)

	wg := &sync.WaitGroup{}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go processGeo(geoChan, licChan, wg)
	}

	// Append geo info
	for _, lic := range eResp.Result {
		geoChan <- lic
	}

	// Close the input channel when we've added al
	close(geoChan)

	// Close the output channel once all processGeos have finished
	wg.Wait()
	close(licChan)

	saveLic(licChan)

	log.Println("done loading")
}

// loadInspections loads all inspections for 11222 and appends them to
// licenses in cachedLicenses
func loadInspections() {
	esv := enigmaSearchVal{
		key:    "zipcode",
		values: []string{"11222", "11211"},
	}

	// Call enigma until no more results
	pageInt := 1
	for {
		enigmaURL := getEnigmaURL(
			`us.states.ny.cities.nyc.dohmh.restaurants.inspections`,
			[]enigmaSearchVal{esv},
			fmt.Sprint(pageInt),
			// Sort by ascending date, so the final values are most recent
			"inspdate+",
		)

		log.Println(enigmaURL)

		iResp := inspectionResponse{}
		err := getJSON(enigmaURL, &iResp)
		if err != nil {
			return
		}

		// If no results, stop the loop
		if len(iResp.Result) < 1 {
			break
		}

		// Go over each result and update new values
		for _, resp := range iResp.Result {
			// The key is the addresss
			key := fmt.Sprint(resp.Building, " ", resp.Street)

			if cachedLicenses[key] != nil {
				log.Printf("%+v", resp)
				cachedLicenses[key].CurrentGrade = resp.CurrentGrade
				cachedLicenses[key].BusinessName = resp.BusinessName
				cachedLicenses[key].InspectionDate = resp.InspectionDate
			}
		}

		pageInt++

	}

}

// initVals is a handler that returns JSON to initialize the map
func initVals(w http.ResponseWriter, r *http.Request) {
	i := initValues{}
	i.Lat = 40.7263
	i.Lon = -73.9456
	i.DefaultZoom = 15
	i.UrlTemplate = `http://otile1.mqcdn.com/tiles/1.0.0/map/{z}/{x}/{y}.png`
	i.Options.MaxZoom = 20

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
	log.SetFlags(log.Lshortfile)
	loadDrinks()
	loadInspections()

	http.HandleFunc("/svc/init", initVals)
	http.HandleFunc("/svc/drinks", drinks)
	http.HandleFunc("/", homepage)

	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		panic(err)
	}
}
