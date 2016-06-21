package jobs

import (
	"net"
	"net/rpc"

	"github.com/valerykalashnikov/zigzag/remote"
)

func StartReplicationService(port string) error {
	replicate := new(remote.Replicate)
	rpc.Register(replicate)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go rpc.ServeConn(conn)
	}
	return nil
}
