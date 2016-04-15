# ZigZag
Yet another in-memory data-structure storage

ZigZag is a golang implementation of in-memory key:value store/cache similar to Redis that is suitable for applications running on a single machine(replication and sharding support is planning in the near future).

##Supports

ZigZag supports two types of storage:

1. Simple storage - requires less resources but is locked at each operation in concurrent mode.
2. Sharded cache - requires more resources but allows to reduce amount of locks in concurrent mode.

##You can use it as

1\. A thread-safe ```map[string]interface{}``` with expiration times without transmiting and serialization data over network
~~~go
import "github.com/valerykalashnikov/zigzag/zigzag"

// Firstly you have to implement the structure to store ttl values
// For example to store ttl values in minutes
type Clock struct{
  ex int64
}

func (c *Clock) Now() time.Time { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

//initialize cache with the type that you want

db, err := zigzag.New("cache")

//or

db, err := zigzag.New("sharded")

// and ivoke the methods provided to you

// to set value
moment := &Clock{}
db.Set(key, value, moment)

// to get value
value, found := db.Get(key)

// to update value

db.Upd(key, value)

//to delete value
db.Del(key)

//to get all the keys which is stored in the cache matching pattern
pattern := "^[a-z]+[[0-9]+]$"
keys := db.Keys(pattern)
~~~

*Important*: ZigZag provides the function getting the n random items from the storage and checking it for expired items.
If expired items were more than 25% it will run again.
Using running the passive expiration is up to you.
~~~go
  n := 5
  db.DelRandomExpires(n)
~~~

2\. A data-structure storage with an JSON api.
  Supported features:
  * JSON API
  * Active and passive values expiration
  * Authentication
  * Syncronization with disk

## How to install:
* ```go get -u github.com/valerykalashnikov/zigzag/zigzag_server```
* ```zigzag_server```


## Supported options by setting proper env variables:

* Setting proper cache type

Simple cash is used by default. If you want to use sharded cache you have to set ```ZIGZAG_ENGINE_TYPE``` env value
~~~bash
ZIGZAG_ENGINE_TYPE=sharded
~~~

* Setting up syncronization with disk:
~~~bash
ZIGZAG_BACKUP_FILE=path_to_file
# for example  ZIGZAG_BACKUP_FILE=/var/lib/zigzag/storage.zz
ZIGZAG_BACKUP_INTERVAL=interval_in_minutes
# for example  ZIGZAG_BACKUP_FILE=2
~~~

* Authentication
~~~bash
  ZIGZAG_AUTH=password
~~~
* Set port:
~~~bash
  ZIGZAG_PORT=3000 # default 8082
~~~


## How to use:

###Set
Set key to hold the value. If key already holds a value, it is overwritten, regardless of its type. Any previous time to live associated with the key overwrites.

* Without expiration time

``` curl -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"Todo"}' http://localhost:8082/set/your_key```

* With expiration time(in minutes)

``` curl -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"Todo"}' http://localhost:8082/set/your_key?ex=1 ```

###Get
Get the value of key. If the key does not exist 404 is returned

```curl -X GET -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/get/your_key```


###Update
Update the specified key. Time to live will not be overwritten.

```curl -X PUT -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"New todo"}' http://localhost:8082/update/your_key```

###Keys
Returns all keys matching pattern
For example, all keys matching ^[a-z]* pattern (dont' forget about escaping)

```curl -X GET -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/keys/%5E\[a-z\]\*```

###Delete
Remove the specified keys. A key is ignored if it does not exist

```curl -X DELETE -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/delete/your_key```





