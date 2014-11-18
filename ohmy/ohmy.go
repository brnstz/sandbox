// Attempt to pull from ohmyrockess.com
package ohmy

import (
	"github.com/PuerkitoBio/goquery"

	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	base = `http://www.ohmyrockness.com/`
	api  = `http://www.ohmyrockness.com/api/shows.json?index=true&page=1&per=50&regioned=1`
)

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
	req.Header.Add(`Authorization`, `Token token="3b35f8a73dabd5f14b1cac167a14c1f6"`)
	apiRes, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer apiRes.Body.Close()

	b, err := ioutil.ReadAll(apiRes.Body)
	fmt.Printf("%s", b)
}
