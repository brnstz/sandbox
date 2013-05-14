package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {

	http.Handle("www.wholok.com/", http.FileServer(http.Dir("/home/bseitz/www/wholok.com")))
	http.Handle("wholok.com/", http.FileServer(http.Dir("/home/bseitz/www/wholok.com")))
	http.Handle("www.brnstz.com/", http.FileServer(http.Dir("/home/bseitz/www/brnstz.com")))
	http.Handle("brnstz.com/", http.FileServer(http.Dir("/home/bseitz/www/brnstz.com")))
	http.Handle("themeasurements.com/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir("/home/bseitz/www/themeasurements.com/logo/"))))
	http.Handle("www.themeasurements.com/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir("/home/bseitz/www/themeasurements.com/logo/"))))
	http.HandleFunc("themeasurements.com/",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://themeasurements.bandcamp.com/", http.StatusFound)
		})
	http.HandleFunc("www.themeasurements.com/",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://themeasurements.bandcamp.com/", http.StatusFound)
		})

	highlit_url, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	http.Handle("highlit.brnstz.com/", httputil.NewSingleHostReverseProxy(highlit_url))

	inbox_url, err := url.Parse("http://localhost:8081")
	if err != nil {
		panic(err)
	}
	http.Handle("inbox.brnstz.com/", httputil.NewSingleHostReverseProxy(inbox_url))

	err = http.ListenAndServe(":80", nil)

	if err != nil {
		panic(err)
	}
}
