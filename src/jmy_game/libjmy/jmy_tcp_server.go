package libjmy

import (
	"fmt"
	"net"
)

type JmyTcpServer struct {
	listener net.Listener
	conn_mgr JmyTcpConnectionMgr
}

func (this *JmyTcpServer) init(max_conn uint32) bool {
	this.conn_mgr.init(max_conn)
	return true
}

func (this *JmyTcpServer) listen(addr string) bool {
	var err error
	this.listener, err = net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen addr:%v failed, err:%v", addr, err)
		return false
	}
	return true
}

func (this *JmyTcpServer) accept() bool {
	if this.listener == nil {
		return false
	}
	conn, err := this.listener.Accept()
	if err != nil {
		fmt.Printf("accept new conn failed")
		return false
	}
	c := this.conn_mgr.getFree()
	if c == nil {
		conn.Close()
		fmt.Printf("get free tcp connection failed")
		return false
	}
	c.create_from(conn)
	if !c.start() {
		fmt.Printf("tcp connection start failed")
		return false
	}
	return true
}
