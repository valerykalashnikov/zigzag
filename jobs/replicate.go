package jobs

import (
	"encoding/gob"
	"net"
	"sync"

	"github.com/valerykalashnikov/zigzag/cache"
	"github.com/valerykalashnikov/zigzag/structures"
	"github.com/valerykalashnikov/zigzag/zigzag"
)

func StartReplicationService(wg sync.WaitGroup, slave *zigzag.DB) {
	var handleConnection = func(conn net.Conn) {
		gob.Register(map[string]interface{}{})
		gob.RegisterName("*structures.Clock", &structures.Clock{})

		dec := gob.NewDecoder(conn)
		importedCache := &cache.CacheImport{}
		dec.Decode(importedCache)
		slave.Set(importedCache.Key, importedCache.Value, importedCache.M)
	}

	listener, err := net.Listen("tcp", slave.GetReplicationPort())
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
