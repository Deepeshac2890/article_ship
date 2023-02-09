package main

import (
	Db "article_ship/Database"
	"article_ship/Helper"
	"article_ship/Models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

var articles Models.Articles

func allArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("End point hit : All articles endpoint")
	articles = Db.GetAllArticles()
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
	var article Models.Article
	err := json.Unmarshal(reqBody, &article)
	if err != nil {
		return
	}
	isInserted := Db.InsertArticle(article)
	if isInserted {
		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			return
		}
	} else {
		err = json.NewEncoder(w).Encode("Article was not inserted")
		if err != nil {
			return
		}
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
	vars := mux.Vars(req)
	key := vars["id"]
	id := Helper.StringToInt32(key)
	reqBody, _ := io.ReadAll(req.Body)
	var article Models.Article
	err := json.Unmarshal(reqBody, &article)
	if err != nil {
		err := json.NewEncoder(w).Encode("Something went wrong !! Problem with request")
		if err != nil {
			panic(err.Error())
		}
	} else {
		if Db.UpdateArticle(article, id) {
			err = json.NewEncoder(w).Encode("Article has been updated !!")
			if err != nil {
				panic(err.Error())
			}
		} else {
			err = json.NewEncoder(w).Encode("Something went wrong !!")
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func returnSingleArticle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["id"]
	id := Helper.StringToInt32(key)

	article, _ := Db.GetSingleArticle(id)
	err := json.NewEncoder(w).Encode(article)
	if err != nil {
		return
	}
}

func deleteAllArticle(w http.ResponseWriter, r *http.Request) {
	isSuccess := Db.DeleteAllArticles()
	if isSuccess {
		err := json.NewEncoder(w).Encode("All articles have been deleted")
		if err != nil {
			return
		}
	} else {
		err := json.NewEncoder(w).Encode("Something went wrong!!")
		if err != nil {
			return
		}
	}
}

func deleteSingleArticle(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	key := vars["id"]

	hasDeleted := Db.DeleteSingleArticle(Helper.StringToInt32(key))
	// we then need to loop through all our articles
	if hasDeleted {
		err := json.NewEncoder(w).Encode("Article has been deleted")
		if err != nil {
			return
		}
	} else {
		err := json.NewEncoder(w).Encode("Something went wrong!!")
		if err != nil {
			return
		}
	}
}

func main() {
	articles = Models.Articles{}
	Db.InitConnection()
	handleRequests()
}
