package main

import (
	"flag"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	// "os"
	// "os/signal"
	// "syscall"
	"database/sql"
	_ "github.com/lib/pq"

	repo "github.com/interviewr/interviewr-server/repository/postgres"
	usecases "github.com/interviewr/interviewr-server/usecases"
	httpDeliver "github.com/interviewr/interviewr-server/transport/http"
	cfg "github.com/interviewr/interviewr-server/config/env"
)

var config cfg.Config

func init() {
	config = cfg.NewViperConfig()

	if config.GetBool("debug") {
		fmt.Println("Service run on debug mode")
	}
}

func main() {
	dbHost := config.GetString("database.host")
	dbPort := config.GetString("database.port")
	dbUser := config.GetString("database.user")
	dbPass := config.GetString("database.pass")
	dbName := config.GetString("database.name")
	connectionStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	dbConnection, err := sql.Open("postgres", connectionStr)
	if err != nil && config.GetBool("debug") {
		fmt.Println(err)
	}
	defer dbConnection.Close()

	err = dbConnection.Ping()
	if err != nil {
		panic(err)
	}

	var (
		httpAddr = flag.String("http.addr", ":8090", "HTTP listen address")
	)

	flag.Parse()

	r := mux.NewRouter()

	orgRepo := repo.NewOrganizationRepository(dbConnection)
	orgUsecase := usecases.NewOrganizationUsecase(orgRepo)
	httpDeliver.NewOrganizationHttpHandler(r, orgUsecase)

	// errs := make(chan error)
	// go func() {
	// 	c := make(chan os.Signal)
	// 	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	// 	errs <- fmt.Errorf("%s", <-c)
	// }()

	// go func() {
	// 	errs <- http.ListenAndServe(*httpAddr, r)
	// }()

	http.ListenAndServe(*httpAddr, r)
}