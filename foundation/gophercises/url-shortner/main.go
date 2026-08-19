package main

import (
	"fmt"
	"net/http"
)

func main() {
	urls := map[string]string{
		"/prem":     "https://portfolio.premgowda.in",
		"/github":   "https://github.com/premgowda98",
		"/linkedin": "https://linkedin.com/in/premgowda98",
	}

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home Page")
	})

	handler := redirectHandler(urls, router)

	fmt.Println("Server is listening at port 8090")
	http.ListenAndServe(":8090", handler)
}

func redirectHandler(urls map[string]string, fallbach http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		dest, ok := urls[path]

		if !ok {
			fallbach.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, dest, http.StatusFound)
	}
}
