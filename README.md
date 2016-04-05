# ZigZag
Yet another in-memory data-structure storage

ZigZag is an in-memory key:value store/cache similar to Redis that is suitable for applications running on a single machine(replication and sharding support is planning in the near future).

##You can use it as

1.a thread-safe ```map[string]interface{}``` with expiration times without transmiting and serialization data over network
~~~go
import "zigzag/zigaz"

// Firstly you have to implement the structure to store ttl values
// For example to store ttl values in minutes
type Clock struct{
  ex int64
}

func (c *Clock) Now() time.Time { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.ex) * time.Minute }

// and ivoke the methods provided to you

// to set value
moment := &Clock{}
zigzag.Set(key, value, moment)

// to get value
value, found := zigzag.Get(key)

// to update value

zigzag.Upd(key, value)

//to delete value
zigzag.Del(key)
~~~

Important: ZigZag provides the function getting the n random items from the storage and checking it for expired items.
If expired items were more than 25% it will run again.
Using running the passive expiration is up to you.
~~~go
  zigzag.DelRandomExpires(n)
~~~

2.a data-structure storage with an JSON api.
  Supported features:
  * JSON API
  * Active and passive values expiration
  * Authentication
  * Syncronization with disk

## How to install:
* Land the code to the ```${GOPATH}/src```
* ```cd to http-server```
* ```go build -o zigzag-server```
* ```./zigzag-server ```

## Supported options by setting proper env variables:

* Setting up syncronization with disk:
~~~bash
ZIGZAG_BACKUP_FILE=path_to_file
# for example  ZIGZAG_BACKUP_FILE=/var/lib/zigzag/storage.zz
ZIGZAG_BACKUP_INTERVAL=interval_in_minutes
# for example  ZIGZAG_BACKUP_FILE=2
~~~

2. Authentication
~~~bash
  ZIGZAG_AUTH=password
~~~
3. Set port:
~~~bash
  ZIGZAG_PORT=3000 # default 8082
~~~


## How to use:

###Set
Set key to hold the value. If key already holds a value, it is overwritten, regardless of its type. Any previous time to live associated with the key overwrites.

* Without expiration time

``` curl -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"Todo"}' http://localhost:8080/set/your_key```

* With expiration time(in minutes)

``` curl -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"Todo"}' http://localhost:8080/set/your_key?ex=1 ```

###Get
Get the value of key. If the key does not exist 404 is returned

```curl -X GET -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/get/your_key```

###Delete
Remove the specified keys. A key is ignored if it does not exist

```curl -X DELETE -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/delete/your_key```

###Update
Update the specified key. Time to live will not be overwritten.

```curl -X PUT -v -H "Content-Type: application/json" -H "Authorization: Token password" -d '{"name":"Todo"}' http://localhost:8082/update/key```

###Keys
Returns all keys matching pattern
For example, all keys matching ^[a-z]* pattern (dont' forget about escaping)

```curl -X GET -v -H "Content-Type: application/json" -H "Authorization: Token password" http://localhost:8082/keys/%5E\[a-z\]\*```





