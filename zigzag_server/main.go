package main

import (
  "log"
  "net/http"
  "os"
  "sync"
  "os/signal"
  "syscall"
  "fmt"
  "strconv"

  "github.com/valerykalashnikov/zigzag/jobs"
  "github.com/valerykalashnikov/zigzag/importers"
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

func runBackgroundJobs() sync.WaitGroup {
  var wg sync.WaitGroup

  backupFilePath := os.Getenv("ZIGZAG_BACKUP_FILE")

  backupInterval := os.Getenv("ZIGZAG_BACKUP_INTERVAL")

  if backupInterval != "" {
    period, err := strconv.Atoi(backupInterval)

    if err != nil { fmt.Println(err) }

    go jobs.SaveToFile(wg, backupFilePath, period)
  }

  go jobs.CleanCache(wg, 20)

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

func ImportCache(path string) {
  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println(" - Nothing to import")
    return
  }
  err := importers.ImportCacheFromFile(path)
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(" - Cache successfully imported")
}

func main() {

  port, authToken := os.Getenv("ZIGZAG_PORT"), os.Getenv("ZIGZAG_AUTH")
  if port == "" { port = "8082" }
  fmt.Print(lightning)

  backupFilePath := os.Getenv("ZIGZAG_BACKUP_FILE")
  if backupFilePath != "" && os.Getenv("ZIGZAG_BACKUP_INTERVAL") != "" {
    fmt.Println("* Importing cache from ", backupFilePath)
    ImportCache(backupFilePath)
  }

  fmt.Println("* Running background jobs...")
  wg := runBackgroundJobs()

  handleInterruptSignal(wg)

  router := NewRouter(authToken)
  fmt.Println("* Listening on http://localhost:" + port)
  log.Fatal(http.ListenAndServe(":" + port, router))
}


