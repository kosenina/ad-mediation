package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kosenina/ad-mediation/adding"
	"github.com/kosenina/ad-mediation/listing"
	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/storage"
	"github.com/kosenina/ad-mediation/utils"
	"gopkg.in/natefinch/lumberjack.v2"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

func main() {
	// Configure Logging
	logFileLocation := utils.GetEnv("LOG_FILE_LOCATION", "")
	if logFileLocation != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFileLocation,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   //days
			Compress:   true, // disabled by default
		})
	}

	log.Println("Starting Server")

	// Prepare repository
	var dbStorage models.Repository

	confStorageType := utils.GetEnv("PERSISTENT_STORAGE", "MongoDB")

	// Initialize service storage
	switch confStorageType {
	case "MongoDB":
		dbStorage, _ = storage.NewMongoDBStorage()
	case "CloudSQL":
		log.Fatal("Not implemented")
	}

	// Check if storage is up and running
	pingErr := dbStorage.Ping()
	if pingErr != nil {
		log.Fatal("Failed to ping DB.", pingErr)
	} else {
		log.Println("Ping to DB was successfull.")
	}

	// Create the available services
	lister := listing.NewService(dbStorage)
	adder := adding.NewService(dbStorage)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/adNetworkList", listing.MakeGetAdNetworkListingEndpoint(lister)).Methods(http.MethodGet)
	api.HandleFunc("/adNetworkList", adding.MakePutAdNetworkListingEndpoint(adder)).Methods(http.MethodPut)
	api.HandleFunc("/adNetworkList", notFound)
	log.Println("Server is up and running!")
	log.Fatal(http.ListenAndServe(":8080", router))
}
