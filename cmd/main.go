package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/isogram/clean-golang/pkg/bookmark"
	"github.com/isogram/clean-golang/pkg/entity"

	config "github.com/joho/godotenv"
	"github.com/juju/mgosession"
	mgo "gopkg.in/mgo.v2"
)

func handleParams() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("Invalid query")
	}
	return os.Args[1], nil
}

func main() {
	err := config.Load(".env")
	if err != nil {
		fmt.Println(".env is not loaded properly")
		os.Exit(2)
	}

	query, err := handleParams()
	if err != nil {
		log.Fatal(err.Error())
	}

	session, err := mgo.Dial(os.Getenv("MONGO_DB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer session.Close()

	cPool, _ := strconv.Atoi(os.Getenv("MONGO_DB_CONNECTION_POOL"))
	mPool := mgosession.NewPool(nil, session, cPool)
	defer mPool.Close()

	bookmarkRepo := bookmark.NewMongoRepository(mPool, os.Getenv("MONGO_DB_DATABASE"))
	bookmarkService := bookmark.NewService(bookmarkRepo)
	all, err := bookmarkService.Search(query)
	if err != nil {
		log.Fatal(err)
	}
	if len(all) == 0 {
		log.Fatal(entity.ErrNotFound.Error())
	}
	for _, j := range all {
		fmt.Printf("%s %s %v \n", j.Name, j.Link, j.Tags)
	}
}
