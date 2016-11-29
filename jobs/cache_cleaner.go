package jobs

import (
	"sync"
	"time"

	"github.com/valerykalashnikov/zigzag"
)

func CleanCache(wg sync.WaitGroup, db *zigzag.DB, checkForExpirationItemNum int) {
	ticker := time.NewTicker(time.Millisecond * 100)
	for range ticker.C {
		wg.Add(1)
		db.DelRandomExpires(checkForExpirationItemNum)
		wg.Done()
	}
}
