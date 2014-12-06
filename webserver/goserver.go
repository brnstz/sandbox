package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"strings"
)

func measurementsRedir(w http.ResponseWriter, r *http.Request) {
	lowerPath := strings.ToLower(r.URL.Path)

	if strings.HasPrefix(lowerPath, "/nodana") {
		http.Redirect(w, r, "https://www.youtube.com/watch?v=CcRk-zdn8Ps", http.StatusFound)
	} else {
		http.Redirect(w, r, "http://themeasurements.bandcamp.com/", http.StatusFound)
	}
}

func main() {

	http.Handle("www.wholok.com/", http.FileServer(http.Dir("/home/bseitz/www/wholok.com")))
	http.Handle("wholok.com/", http.FileServer(http.Dir("/home/bseitz/www/wholok.com")))
	http.Handle("www.brnstz.com/", http.FileServer(http.Dir("/home/bseitz/www/brnstz.com")))
	http.Handle("brnstz.com/", http.FileServer(http.Dir("/home/bseitz/www/brnstz.com")))

	http.Handle("themeasurements.com/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir("/home/bseitz/www/themeasurements.com/logo/"))))
	http.Handle("www.themeasurements.com/logo/", http.StripPrefix("/logo/", http.FileServer(http.Dir("/home/bseitz/www/themeasurements.com/logo/"))))

	http.HandleFunc("themeasurements.com/", measurementsRedir)
	http.HandleFunc("www.themeasurements.com/", measurementsRedir)

	highlitUrl, err := url.Parse("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	http.Handle("highlit.brnstz.com/", httputil.NewSingleHostReverseProxy(highlitUrl))

	inboxUrl, err := url.Parse("http://localhost:8081")
	if err != nil {
		panic(err)
	}
	http.Handle("inbox.brnstz.com/", httputil.NewSingleHostReverseProxy(inboxUrl))

	clusterUrl, err := url.Parse("http://localhost:8082")
	if err != nil {
		panic(err)
	}
	http.Handle("cluster.brnstz.com/", httputil.NewSingleHostReverseProxy(clusterUrl))

	drinkUrl, err := url.Parse("http://localhost:8083")
	if err != nil {
		panic(err)
	}
	http.Handle("drink.brnstz.com/", httputil.NewSingleHostReverseProxy(drinkUrl))

	listenUrl, err := url.Parse("http://localhost:8084")
	if err != nil {
		panic(err)
	}
	http.Handle("listen.brnstz.com/", httputil.NewSingleHostReverseProxy(listenUrl))

	err = http.ListenAndServe(":80", nil)

	if err != nil {
		panic(err)
	}
}
