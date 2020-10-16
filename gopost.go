package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/olivere/elastic.v7"

	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

type Article struct {
	Id      string `json:Id`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	fmt.Println("%s\n", article)
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	if err != nil {
		// Handle error
		panic(err)
	}
	dataJSON, err := json.Marshal(article)
	js := string(dataJSON)
	ind, err := client.Index().
		Index("students").
		BodyJson(js).
		Do(ctx)
	fmt.Println("[Elastic][InsertProduct]Insertion Successful", ind, ctx)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	// NOTE: Ordering is important here! This has to be defined before
	// the other `/article` endpoint.
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/all", returnAllArticles)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	Articles = []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}

	handleRequests()

}
