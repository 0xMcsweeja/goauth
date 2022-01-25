package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"main/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Status struct {
	Status string
	Time   time.Time
}

type Handlers struct {
	store map[string]Status
}
type ClientStore struct {
	Store map[string]ClientStruct
}

func main() {
	Store := generateUsers()

	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Printf("Starting server  on http://localhost:%s", port)

	router := chi.NewRouter()

	//middleware stack
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

	//routes
	router.Route("/token", func(router chi.Router) {
		log.Printf("Route: http://localhost:%s", string(port)+"/token")

		router.Get("/", Token)
	})

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {

		jsonValue, _ := json.Marshal(Store)
		w.Write(jsonValue)

	})
	router.Get("/users/1", func(w http.ResponseWriter, r *http.Request) {
		VisibleClientDetails(w, r, Store.Store["1"])

	})

	//router.Route("/startup", func(w http.ResponseWriter, r *http.Request) {
	//	router.Get("/", Startup(w, r, Store))
	//})
	http.ListenAndServe(":8080", router)
}

func generateUsers() *ClientStore {
	clients := new(ClientStore)
	clients.Store = make(map[string]ClientStruct)

	connectionString := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Credentials{})
	if err != nil {
		panic("failed to connect database")
	}

	fnames := []string{"Sleve", "Onson", "Darryl", "Anatoli", "Mario", "Kevin", "Mike", "Raul", "Willie", "Jeromy", "Scott", "Karl"}
	lnames := []string{"McDichael", "Sweeney", "Smorin", "Truk", "Nogilniy", "Dandleton", "Rortugal", "McScriff", "Bonzalez", "Dugnutt", "Sernandez", "Marx"}
	shuffledFnames := make([]string, len(fnames))
	shuffledLnames := make([]string, len(fnames))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(fnames))
	perm2 := rand.Perm(len(fnames))

	for i, v := range perm {
		shuffledFnames[i] = fnames[v]
	}
	for i, v := range perm2 {
		shuffledLnames[i] = lnames[v]
	}
	creds := new(models.Credentials)
	creds.Credential_type = "shared_secret"
	creds.Secret = "shared_secret"
	for i := range shuffledFnames {
		if db.Model(&models.User{}).Where("ID = ?", uint32(i)).Updates(&models.User{
			Nickname:  shuffledFnames[i],
			Email:     shuffledLnames[i] + "@gmail.com",
			Password:  shuffledFnames[i] + shuffledLnames[i],
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{}, Credentials: *creds}).RowsAffected == 0 {
			db.Create(&models.User{
				Nickname:  shuffledFnames[i],
				Email:     shuffledLnames[i] + "@gmail.com",
				Password:  shuffledFnames[i] + shuffledLnames[i],
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{}, Credentials: *creds})
		}

	}

	return clients

}

//func main() {
//	Store := new(ClientStore)
//	Store.Store = map[string]ClientStruct{}
//
//	handler := newHandler()
//
//	http.HandleFunc("/startup", func(w http.ResponseWriter, r *http.Request) {
//		Startup(w, r, Store)
//	})
//	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
//		Users(w, r, Store)
//	})
//	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
//		Users(w, r, Store)
//	})
//	http.HandleFunc("/", handler.healthchecks)
//	http.HandleFunc("/dogs", handler.getDog)
//	http.HandleFunc("/token", Token)
//	http.HandleFunc("/auth", Authorize)
//	http.HandleFunc("/clientss", Client)
//
//	err := http.ListenAndServe(":8080", nil)
//
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(Store)
//
//}
func Clients(w http.ResponseWriter, r *http.Request, store *ClientStore) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var client ClientStruct
	err = json.Unmarshal(body, &client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	jsonBytes, err := json.Marshal(client)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func hashToString(client int) string {

	h := sha1.New()
	h.Write([]byte(strconv.Itoa(client)))
	return hex.EncodeToString(h.Sum(nil))
}

func (h *Handlers) healthchecks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.get(w, r) //need to update for post
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *Handlers) get(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2) //remove this - solely for concurrency testing

	healthcheck := h.store
	jsonBytes, err := json.Marshal(healthcheck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	fmt.Fprintln(w, "%s", healthcheck)
	w.Write(jsonBytes)
}

func (h *Handlers) getDog(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 2) //remove this - solely for concurrency testing

	healthcheck := "dogs here"
	jsonBytes, err := json.Marshal(healthcheck)
	if err != nil {
		panic(err)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func newHandler() *Handlers {
	return &Handlers{
		store: map[string]Status{
			"healthcheck": Status{
				Status: "healthy",
				Time:   time.Now(),
			},
		},
	}
}
