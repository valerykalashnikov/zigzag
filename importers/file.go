package importers

import (
          "github.com/valerykalashnikov/zigzag/cache"
          "github.com/valerykalashnikov/zigzag/zigzag"
          "github.com/valerykalashnikov/zigzag/persistence"
        )


type FileImport struct{
  path string
}

func (f *FileImport) Import() (map[string]*cache.Item, error) {
  items, err := persistence.RestoreFromFile(f.path)
  return items, err

}


func ImportCacheFromFile(path string) error {
  err := zigzag.ImportCache(&FileImport{path: path})
  return err
}
