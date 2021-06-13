package main

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// tom: for Initialize
	"fmt"
	"log"
	"strings"

	// tom: for route handlers
	"encoding/json"
	"net/http"
	"strconv"

	// tom: go get required
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	// migrate postgres
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// tom: initial function is empty, it's filled afterwards
// func (a *App) Initialize(user, password, dbname string) { }

// tom: added "sslmode=disable" to connection string
func (a *App) Initialize(user, password, host, port, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)

	var err error
	fmt.Println("Connect String ", connectionString)
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Connect to Database Error ", err)
	}
	//m, err := migrate.New(
	//	"file://db/migrations", connectionString)
	driver, err := postgres.WithInstance(a.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	m.Steps(2)
	a.Router = mux.NewRouter()

	// tom: this line is added after initializeRoutes is created later on
	a.initializeRoutes()
}

// tom: initial version
// func (a *App) Run(addr string) { }
// improved version
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) getHome(w http.ResponseWriter, r *http.Request) {
	var response string
	if strings.HasPrefix(r.URL.Path, "/api") {
		response = "success"
	} else {
		response = "failure"
	}
	fmt.Println(response + " - path: " + r.URL.Path + "\n")
	respondWithJSON(w, http.StatusOK, "Ok")
}

// tom: these are added later
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Get Lists:", products)
	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	fmt.Println("Get body ", a.DB)

	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Create pruduct has id: ", p.ID, p)
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Update pruduct has id: ", id, p)

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Delete pruduct has id: ", id)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.getHome).Methods("GET")
	a.Router.HandleFunc("/api", a.getHome).Methods("GET")
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}
