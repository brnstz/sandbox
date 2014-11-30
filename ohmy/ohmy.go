// Attempt to pull from ohmyrockess.com
package ohmy

import (
	"github.com/PuerkitoBio/goquery"

	"fmt"
	"io/ioutil"
	"net/http"
)

type Show struct {
	Bands []*Band `json:"cached_bands"`
}

type Band struct {
	Name string
}

type Track struct {
}

const (
	// The base URL of the API, and the index page. We use index page
	// to get the CSRF token.
	base = "http://www.ohmyrockness.com/"

	// The API endpoint
	//api = "api/shows.json"

	// Value for NY region (parameter name: regioned)
	regioned = 1

	// Number of records per page (parameter name: per)
	per = 50

	api = `http://www.ohmyrockness.com/api/shows.json?index=true&page=1&per=50&regioned=1`

	// FIXME: when does this change?
	token = `Token token="3b35f8a73dabd5f14b1cac167a14c1f6"`
)

/*
func GetShows(page int) *[]Show {
    v := url.Values{}

    v.Add(

}
*/

func Doit() {

	res, err := http.Get(base)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		panic(err)
	}

	var content string
	doc.Find(`meta[name="csrf-token"]`).Each(func(i int, s *goquery.Selection) {
		fmt.Println("hello")
		fmt.Println(s.Attr("content"))
		content, _ = s.Attr("content")
	})

	client := &http.Client{}
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		panic(err)
	}

	for _, c := range res.Cookies() {
		req.AddCookie(c)
		fmt.Println(c.Name, c.Value)
	}
	req.Header.Add(`X-CSRF-Token`, content)
	req.Header.Add(`Authorization`, token)
	apiRes, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer apiRes.Body.Close()

	b, err := ioutil.ReadAll(apiRes.Body)
	fmt.Printf("%s", b)
}
