package web

import (
	"net/http"
	"text/template"
)

func ErrorDisplay(w http.ResponseWriter, errorCode int, errorMessage string) {
	//set HTTP status
	w.WriteHeader(errorCode)

	//open the error html
	tmpl, err := template.ParseFiles("static/errorPage.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	//execute which fills the space in the html depending of th index, 0 for the code, and 1 for the message
	err = tmpl.Execute(w, []interface{}{errorCode, errorMessage})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
