package cache

import (
        "errors"
        "fmt"
      )


type DataStoreFactory func() DataStore

var datastoreFactories = make(map[string]DataStoreFactory)

func Register(name string, factory DataStoreFactory) {
    if factory == nil {
      error := errors.New(fmt.Sprintf("Datastore factory %s does not exist.", name))
      panic(error)
    }
    _, registered := datastoreFactories[name]
    if registered {
      error := errors.New(fmt.Sprintf("Datastore factory %s already registered.", name))
      panic(error)
    }

    datastoreFactories[name] = factory
}

func CreateCache(arguments ...string) (DataStore, error) {
  var engineName string = "cache"
  if len(arguments) > 0 {
    engineName = arguments[0]
  }
  cacheFactory, ok := datastoreFactories[engineName]
  if !ok {
    // Factory has not been registered.
    // Make a list of all available datastore factories for logging.
    availableDatastores := make([]string, len(datastoreFactories))
    for k, _ := range datastoreFactories {
        availableDatastores = append(availableDatastores, k)
    }
    return nil, errors.New(fmt.Sprintf("Invalid Datastore name. Must be one of: %v", availableDatastores))
  }
  datastore := cacheFactory()
  return datastore, nil
}

func init() {
  Register("cache", NewCache)
  Register("sharded", NewShardedCache)
}
