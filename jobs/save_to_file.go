package jobs

import (
	"sync"
	"time"

	"github.com/valerykalashnikov/zigzag/persistence"
	"github.com/valerykalashnikov/zigzag/zigzag"
)

func SaveToFile(wg sync.WaitGroup, db *zigzag.DB, path string, period int) {
	ticker := time.NewTicker(time.Minute * time.Duration(period))
	items := db.Items()
	for range ticker.C {
		wg.Add(1)
		defer wg.Done()
		persistence.SaveToFile(path, items)
	}
}
