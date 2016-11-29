package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"time"

	"github.com/valerykalashnikov/zigzag"
	"github.com/valerykalashnikov/zigzag/importers"
	"github.com/valerykalashnikov/zigzag/jobs"
)

var lightning = `
         zzzzzz/
        zzzzzz/
       zzzzzz/
      zzzzzzzzzzzzzzz
     zzzzzzzzzzzzzz
          /zzzzzz
         /zzzzz
        /zzzz
       /zzz
      /zz
     /z
`

func runBackgroundJobs(db *zigzag.DB) sync.WaitGroup {
	var wg sync.WaitGroup

	backupFilePath := os.Getenv("ZIGZAG_BACKUP_FILE")

	backupInterval := os.Getenv("ZIGZAG_BACKUP_INTERVAL")

	if backupInterval != "" {
		period, err := strconv.Atoi(backupInterval)

		if err != nil {
			fmt.Println(err)
		}

		go jobs.SaveToFile(wg, db, backupFilePath, period)
	}

	go jobs.CleanCache(wg, db, 20)

	if db.CheckRole() == "slave" {
		fmt.Println("* Running replication service...")

		go jobs.StartReplicationService(wg, db)
	}

	return wg
}

func handleInterruptSignal(wg sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Waiting for the workers until they have finished their job...")
		wg.Wait()
		os.Exit(1)
	}()
}

func ImportCache(db *zigzag.DB, path string) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(" - Nothing to import")
		return
	}

	importers.FileImport(db, path)

	fmt.Println(" - Cache successfully imported")
}

func main() {

	port, authToken := os.Getenv("ZIGZAG_PORT"), os.Getenv("ZIGZAG_AUTH")
	if port == "" {
		port = "8082"
	}
	fmt.Print(lightning)

	engineType := os.Getenv("ZIGZAG_ENGINE_TYPE")
	if engineType == "" {
		engineType = "cache"
	}

	role := os.Getenv("ZIGZAG_ROLE")
	if role == "" {
		role = "master"
	}

	db, err := zigzag.New(engineType, role)
	if err != nil {
		panic(err)
	}
	fmt.Println(" - Engine type:", engineType)
	fmt.Println(" - Role:", role)

	backupFilePath := os.Getenv("ZIGZAG_BACKUP_FILE")
	if backupFilePath != "" && os.Getenv("ZIGZAG_BACKUP_INTERVAL") != "" {
		fmt.Println("* Importing cache from ", backupFilePath)
		ImportCache(db, backupFilePath)
	}

	repPort := os.Getenv("ZIGZAG_REPLICATION_PORT")
	if repPort == "" {
		repPort = ":8084"
		zigzag.SetReplicationPort(db, repPort)
	}

	fmt.Println("* Running background jobs...")
	wg := runBackgroundJobs(db)

	handleInterruptSignal(wg)

	router := NewRouter(authToken, db)

	srv := http.TimeoutHandler(router, 3*time.Second, "Timed out")

	fmt.Println("* Listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, srv))
}
