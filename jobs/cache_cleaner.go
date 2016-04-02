package jobs

import (
        "time"
        "zigzag/zigzag"
      )

const checkForExpirationItemNum = 20

func CleanCache() {
  ticker := time.NewTicker(time.Millisecond * 100)
  for range ticker.C {
    zigzag.DelRandomExpires(checkForExpirationItemNum)
  }
}
