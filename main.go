package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"Desc"`
	Content string `json:"Content"`
}

type Articles []Article

var articles Articles

func allArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("End point hit : All articles endpoint")
	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		_, err := fmt.Fprintln(w, http.StatusNotFound)
		if err != nil {
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Homepage Endpoint Hit")
	if err != nil {
		return
	}
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := io.ReadAll(r.Body)
	var article Article
	err := json.Unmarshal(reqBody, &article)
	if err != nil {
		return
	}
	articles = append(articles, article)
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		return
	}
}

func handleRequests() {
	// This mux router helps us determine what verbs we can use to access endpoints.
	// like we can also specify what method GET,POST etc can be used to access a particular endpoint
	// We can have 2 endpoints with same name but have different methods and functions. Like we can have articles for GET and articles for POST
	// This gives use sort of polymorphism.
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// If we try to use POST here then it show 405 method not allowed !!
	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
	// NOTE: Ordering is important here! This has to be defined before
	// the other `/article` endpoint.
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	// This method name we put is case in-sensitive
	myRouter.HandleFunc("/delete/{id}", deleteSingleArticle).Methods("Delete")
	myRouter.HandleFunc("/deleteAll", deleteAllArticle).Methods("Delete")
	myRouter.HandleFunc("/updateArticle/{id}", updateArticle).Methods("put")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func updateArticle(w http.ResponseWriter, req *http.Request) {
	deleteSingleArticle(w, req)
	createNewArticle(w, req)
}

func returnSingleArticle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["id"]

	_, err := fmt.Fprintln(w, "Key : "+key)
	if err != nil {
		return
	}

	for _, article := range articles {
		if article.Id == key {
			err := json.NewEncoder(w).Encode(article)
			if err != nil {
				return
			}
		}
	}
}

func deleteAllArticle(w http.ResponseWriter, r *http.Request) {
	articles = Articles{}
	err := json.NewEncoder(w).Encode(articles)
	if err != nil {
		return
	}
}

func deleteSingleArticle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["id"]

	_, err := fmt.Fprintln(w, "Key : "+key)
	if err != nil {
		return
	}

	// we then need to loop through all our articles
	for index, article := range articles {
		// if our id path parameter matches one of our
		// articles
		if article.Id == key {
			// updates our Articles array to remove the
			// article
			articles = append(articles[:index], articles[index+1:]...)
		}
	}
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		return
	}
}

func main() {
	articles = Articles{}
	handleRequests()
}
