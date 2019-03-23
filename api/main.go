package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/isogram/clean-golang/api/handler"
	localConfig "github.com/isogram/clean-golang/config"
	"github.com/isogram/clean-golang/pkg/bookmark"
	"github.com/isogram/clean-golang/pkg/user"

	config "github.com/joho/godotenv"
	"github.com/juju/mgosession"
	mgo "gopkg.in/mgo.v2"

	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
)

const (
	//DefaultPort default http port
	DefaultPort = 8080
)

func main() {
	err := config.Load(".env")
	if err != nil {
		fmt.Println(".env is not loaded properly")
		os.Exit(2)
	}

	session, err := mgo.Dial(os.Getenv("MONGO_DB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	cPool, _ := strconv.Atoi(os.Getenv("MONGO_DB_CONNECTION_POOL"))
	mPool := mgosession.NewPool(nil, session, cPool)
	defer mPool.Close()

	// initiate database and other connections
	readDB := localConfig.ReadPostgresDB()
	writeDB := localConfig.WritePostgresDB()

	//handlers
	e := echo.New()
	e.Pre(mid.RemoveTrailingSlash())
	e.Use(mid.Recover(), mid.CORS(), mid.Logger())

	//bookmark
	bookmarkRepo := bookmark.NewMongoRepository(mPool, os.Getenv("MONGO_DB_DATABASE"))
	bookmarkService := bookmark.NewService(bookmarkRepo)

	bookmarkHandler := handler.NewHTTPBookmarkHandler(bookmarkService)
	bookmarkGroup := e.Group("/v1/bookmark")
	bookmarkHandler.MakeBookmarkHandlers(bookmarkGroup)

	//user
	userRepo := user.NewRepoPostgres(readDB, writeDB)
	userService := user.NewService(userRepo)

	userHandler := handler.NewHTTPUserHandler(userService)
	userGroup := e.Group("/v1/user")
	userHandler.MakeUserHandlers(userGroup)

	// set REST port
	var port uint16
	if portEnv, ok := os.LookupEnv("API_PORT"); ok {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			port = DefaultPort
		} else {
			port = uint16(portInt)
		}
	} else {
		port = DefaultPort
	}

	listenerPort := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(listenerPort))
}
