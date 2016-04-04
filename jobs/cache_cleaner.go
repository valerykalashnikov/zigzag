package jobs

import (
        "time"
        "zigzag/zigzag"
        "sync"
      )

func CleanCache(wg sync.WaitGroup, checkForExpirationItemNum int) {
  ticker := time.NewTicker(time.Millisecond * 100)
  for range ticker.C {
    wg.Add(1)
    zigzag.DelRandomExpires(checkForExpirationItemNum)
    wg.Done()
  }
}
