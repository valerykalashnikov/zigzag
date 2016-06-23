package jobs

import (
	"encoding/gob"
	"net"

	"github.com/valerykalashnikov/zigzag/cache"
	"github.com/valerykalashnikov/zigzag/zigzag"
)

func StartReplicationService(slave *zigzag.DB, port string) error {
	var handleConnection = func(conn net.Conn) {
		dec := gob.NewDecoder(conn)
		importedCache := &cache.CacheImport{}
		dec.Decode(importedCache)
		slave.Set(importedCache.Key, importedCache.Value, importedCache.M)
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn)
	}
	return nil
}
