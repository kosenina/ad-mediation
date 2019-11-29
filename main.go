package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/kosenina/ad-mediation/docs"
	"github.com/kosenina/ad-mediation/adding"
	"github.com/kosenina/ad-mediation/listing"
	"github.com/kosenina/ad-mediation/models"
	"github.com/kosenina/ad-mediation/objectcache"
	"github.com/kosenina/ad-mediation/storage"
	"github.com/kosenina/ad-mediation/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "not found"}`))
}

// @title Ad Mediation Swagger API
// @version 1.0
// @description Swagger API for Golang Project Ad Mediation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email luka.kosenina@outlook.com

// @BasePath /api/v1
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

	log.Println("INFO: Starting Server")

	// Prepare repository
	var dbStorage models.Repository

	confStorageType := utils.GetEnv("PERSISTENT_STORAGE", "MongoDB")

	// Initialize service storage
	switch confStorageType {
	case "MongoDB":
		dbStorage = storage.NewMongoDBStorage()
	case "CloudStorage":
		googleProjectID := utils.GetEnv("GOOGLE_PROJECT_ID", "provide in configuration file")
		dbStorage = storage.NewGoogleCloudStorage(googleProjectID)
	}

	// Check if storage is up and running
	pingErr := dbStorage.Ping()
	if pingErr != nil {
		log.Fatal("ERROR: Failed to ping DB.", pingErr)
	} else {
		log.Println("INFO: Ping to DB was successfull.")
	}

	// Initialize cache
	var cache objectcache.ObjectCache
	cache = objectcache.NewInMemoryCache()

	// Create the available services
	lister := listing.NewService(dbStorage, cache)
	adder := adding.NewService(dbStorage, cache)

	// Set the router as the default one shipped with Gin
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup route group for the API
	ginAPI := router.Group("/api/v1")
	ginAPI.GET("/adNetworkList", listing.MakeGetAdNetworkListingEndpoint(lister))
	ginAPI.PUT("/adNetworkList", adding.MakePutAdNetworkListingEndpoint(adder))

	// Start and run the server
	log.Println("INFO: Server is up and running!")
	log.Fatal(router.Run(":8080"))
}
