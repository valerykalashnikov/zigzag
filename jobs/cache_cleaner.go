package jobs

import (
        "time"
        "zigzag/zigzag"
      )

func CleanCache(checkForExpirationItemNum int) {
  ticker := time.NewTicker(time.Millisecond * 100)
  for range ticker.C {
    zigzag.DelRandomExpires(checkForExpirationItemNum)
  }
}
