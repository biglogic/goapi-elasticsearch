package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/olivere/elastic.v7"
	"gopkg.in/yaml.v2"

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

type Configuration struct {
	EShost     string `yaml:"EShost"`
	ESPort     string `yaml:"ESPort"`
	ESUsername string `yaml:"ESUsername"`
	ESPassword string `yaml:"ESPassword"`
}

//global variables
var obj Configuration
var Articles []Article
var article Article

// var configuration Configuration

func createclient(host, Port, Username, Password string) {
	ctx := context.Background()
	fmt.Println(host)
	client, err := elastic.NewClient(elastic.SetURL("http://"+host+":"+Port), elastic.SetSniff(true), elastic.SetBasicAuth(Username, Password))
	if err != nil {
		// Handle error
		panic(err)
	}
	dataJSON, err := json.Marshal(article)
	js := string(dataJSON)
	ind, err := client.Index().
		Index("json").
		BodyJson(js).
		Do(ctx)
	fmt.Println("[Elastic][InsertProduct]Insertion Successful", ind, ctx)
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	//var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include
	// our new Article
	fmt.Println("%s\n", article)
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
	createclient(obj.EShost, obj.ESPort, obj.ESUsername, obj.ESPassword)
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

	var filename string
	flag.StringVar(&filename, "f", "config.yml", "YAML file to parse.")
	flag.Parse()

	if filename == "" {
		fmt.Println("Please provide yaml file by using -f option")
		return
	}
	yamlfile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	err = yaml.Unmarshal(yamlfile, &obj)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
	}

	handleRequests()

}
