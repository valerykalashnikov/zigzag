package importers

import (
	"github.com/valerykalashnikov/zigzag"
	"github.com/valerykalashnikov/zigzag/persistence"
	"time"
)

type ClockForImport struct {
	duration int64
}

func (c *ClockForImport) Now() time.Time          { return time.Now() }
func (c *ClockForImport) Duration() time.Duration { return time.Duration(c.duration) }

func FileImport(db *zigzag.DB, path string) error {
	items, err := persistence.RestoreFromFile(path)
	if err != nil {
		return err
	}
	for key, item := range items {
		db.Set(key, item.Object, &ClockForImport{item.ExpireAt})
	}
	return nil
}
