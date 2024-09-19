package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	web "web/func"
)

func main() {
	// serve static folder
	http.HandleFunc("/", homeHandle)
	http.HandleFunc("/generate", gHandler)

	//serve the file, no need :)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//don't use read
	//html, err := os.ReadFile("static/index.html")

	// if err != nil {
	// 	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Fprintf(w, "%s", html)

	templ, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, nil)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func gHandler(w http.ResponseWriter, r *http.Request) {
	//check that the method is post ONLY
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Request Method", http.StatusMethodNotAllowed)
		return
	}

	wordTofind := r.FormValue("inputText")
	banner := r.FormValue("banner")

	// Debug output
	fmt.Printf("Received: wordTofind='%s', banner='%s'\n", wordTofind, banner)

	wordTofind, file, err := web.Validation(wordTofind, banner, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Debug output
	fmt.Printf("After validation: wordTofind='%s', file='%s'\n", wordTofind, file)

	fileArray, err := web.Convert(file)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	splitstring := strings.Split(wordTofind, "\r\n")
	result := ""

	for _, word := range splitstring {
		if word == "" {
			result += "\n"
		} else {
			asciiArtWeb, err := web.GenerateAscii(word, fileArray)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			result += asciiArtWeb
		}
	}

	// Debug output
	fmt.Printf("Generated ASCII art:\n%s\n", result)

	//create the temp
	templ, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	//print the result on the same page of index.html
	err = templ.Execute(w, result)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}
