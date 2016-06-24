package jobs

import (
	"encoding/gob"
	"net"
	"sync"

	"github.com/valerykalashnikov/zigzag/cache"
	"github.com/valerykalashnikov/zigzag/zigzag"
)

func StartReplicationService(wg sync.WaitGroup, slave *zigzag.DB, port string) {
	var handleConnection = func(conn net.Conn) {
		dec := gob.NewDecoder(conn)
		importedCache := &cache.CacheImport{}
		dec.Decode(importedCache)
		slave.Set(importedCache.Key, importedCache.Value, importedCache.M)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	for {
		wg.Add(1)
		defer wg.Done()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		handleConnection(conn)
	}
}
