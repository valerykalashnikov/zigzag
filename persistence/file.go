package persistence

 import (
  "encoding/gob"
  "os"

  "zigzag/cache"
 )

func SaveToFile(path string, data map[string]*cache.Item) error{
  dataFile, err := os.Create(path)
  defer dataFile.Close()

  if err != nil {
    return err
  }

  dataEncoder := gob.NewEncoder(dataFile)
  dataEncoder.Encode(data)
  return nil
}

func RestoreFromFile(path string) (map[string]*cache.Item, error) {
  var data map[string]*cache.Item

  dataFile, err := os.Open(path)
  defer dataFile.Close()

  if err != nil {
    return nil, err
    os.Exit(1)
  }

  dataDecoder := gob.NewDecoder(dataFile)
  err = dataDecoder.Decode(&data)

  if err != nil {
    return nil, err
  }
  return data, nil

}
