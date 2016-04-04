package jobs

import (
        "time"
        "zigzag/persistence"
        "zigzag/zigzag"
        "sync"
      )

func SaveToFile(wg sync.WaitGroup, path string ,period int) {
  ticker := time.NewTicker(time.Minute * time.Duration(period))
  cache := zigzag.GetCache()
  for range ticker.C {
    wg.Add(1)
    persistence.SaveToFile(path, cache.Items)
    wg.Done()
  }
}
