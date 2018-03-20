package main

import (
	"database/sql"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Config struct handling config file access
type Config struct {
	Port           string `json:"port"`
	RouterEndpoint string `json:"routerendpoint"`
	Database       struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Dbname   string `json:"dbname"`
	} `json:"database"`
}

//App struct handling router & db
type App struct {
	Router   *mux.Router
	DB       *sql.DB
	myConfig Config
}

//InitializeRoutes responsible for route handling
func (a *App) initializeRoutes() {
	RouterEndpointInt := fmt.Sprintf("%s/%s", a.myConfig.RouterEndpoint, "{id:[0-9]+}")
	a.Router.HandleFunc(a.myConfig.RouterEndpoint, a.getDomains).Methods("GET")
	a.Router.HandleFunc(RouterEndpointInt, a.getDomain).Methods("GET")
	a.Router.HandleFunc(a.myConfig.RouterEndpoint, a.createDomain).Methods("POST")
	a.Router.HandleFunc(RouterEndpointInt, a.updateDomain).Methods("PUT")
	a.Router.HandleFunc(RouterEndpointInt, a.deleteDomain).Methods("DELETE")
}

//Initialize db & routing
func (a *App) Initialize() {
	var err error
	a.myConfig, err = LoadConfiguration("config.json")
	if err != nil {
		log.Fatal(err)
	}
	connectionString := fmt.Sprintf("%s:%s@/%s", a.myConfig.Database.Username, a.myConfig.Database.Password, a.myConfig.Database.Dbname)

	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

//Run router
func (a *App) Run() {
	log.Fatal(http.ListenAndServe(a.myConfig.Port, a.Router))
}

func main() {
	a := App{}
	log.Info("Opening DB...")
	a.Initialize()
	log.Info("Starting Server...")
	a.Run()
}
