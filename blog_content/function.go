package blog_content

import (
	// "context"
	"encoding/json"
	// firebase "firebase.google.com/go"
	// "firebase.google.com/go/db"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

type Post struct {
	Slug             string `json:"slug"`
	Thumbnail        string `json:"thumbnail"`
	ThumbnailAltText string `json:"thumbnailAltText"`
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	Content          string `json:"content"`
}

type Posts []Post

var (
	routes = Routes{
		Route{
			"GetAll",
			http.MethodGet,
			"/",
			getAllContent,
		},
		Route{
			"Get",
			http.MethodGet,
			"/{token}",
			getContent,
		},
	}
)

// var client *db.Client
//
// func init() {
// 	ctx := context.Background()
// 	conf := &firebase.Conig{
// 		DatabaseURL: os.Getenv("FIREBASE_URL")
// 	}
// 	app, err := firebase.NewApp(ctx, conf)
// 	if err != nil {
// 		log.Fatalf("firebase.NewApp: %v", err)
// 	}
// 	client, err = app.Database(ctx)
// 	if err != nil {
// 		log.FatalF("app.Firestore: %v", err)
// 	}
// }

func HandleContent(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	router.ServeHTTP(w, r)
}

func getAllContent(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("../data/data.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var posts Posts

	if err = json.Unmarshal(byteValue, &posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

}

func getContent(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	jsonFile, err := os.Open("../data/data.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var posts Posts

	if err = json.Unmarshal(byteValue, &posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}
	for _, post := range posts {
		if post.Slug == token {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(post); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(err)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}
}
