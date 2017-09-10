package network

import (
	"fmt"
	"net"
	"time"
)

const (
	TCP_SERVER_STATE_NOT_START = iota
	TCP_SERVER_STATE_LISTENING
	TCP_SERVER_STATE_CLOSING
	TCP_SERVER_STATE_CLOSED
)

type TcpServer struct {
	listener  *net.TCPListener
	conn_mgr  *TcpConnectionMgr
	processor IConnProcessor
	state     int
	conn_chan chan net.Conn
}

func CreateTcpServer() *TcpServer {
	server := &TcpServer{}
	server.state = TCP_SERVER_STATE_NOT_START
	return server
}

func (this *TcpServer) Init(max_conn uint32, processor IConnProcessor) bool {
	if this.state != TCP_SERVER_STATE_NOT_START && this.state != TCP_SERVER_STATE_CLOSED {
		return false
	}
	this.conn_chan = make(chan net.Conn, max_conn)
	this.conn_mgr = &TcpConnectionMgr{}
	this.processor = processor
	return this.conn_mgr.Init(TCP_CONNECTION_TYPE_PASSIVE, max_conn)
}

func (this *TcpServer) Clear() {

}

func (this *TcpServer) Close() {
	if this.state != TCP_SERVER_STATE_LISTENING {
		return
	}
	this.state = TCP_SERVER_STATE_CLOSING
}

func (this *TcpServer) listening_routine() {
	var conn net.Conn
	var err error
	for {
		conn, err = this.listener.Accept()
		if err != nil {
			fmt.Printf("listen failed, err(%v)", err.Error())
			this.state = TCP_SERVER_STATE_CLOSING
			break
		}
		fmt.Printf("get new accept conn\n")
		this.conn_chan <- conn
	}
}

func (this *TcpServer) get_listened_conn() net.Conn {
	var c net.Conn
	select {
	case c = <-this.conn_chan:
		{

		}
	default:
		break
	}
	return c
}

func (this *TcpServer) StartListen(addr string) bool {
	a, e := net.ResolveTCPAddr("tcp4", addr)
	if e != nil {
		fmt.Printf("%v\n", e.Error())
		return false
	}

	this.listener, e = net.ListenTCP("tcp", a)
	if e != nil {
		fmt.Printf("%v\n", e.Error())
		return false
	}
	defer this.listener.Close()

	this.state = TCP_SERVER_STATE_LISTENING

	go this.listening_routine()
	fmt.Printf("Start listening %v\n", addr)

	for {
		// 被关闭
		if this.state == TCP_SERVER_STATE_CLOSING {
			break
		}

		c := this.get_listened_conn()
		if c != nil {

			// 生成新连接
			conn := this.conn_mgr.NewConn2(c, this.processor)
			if conn == nil {
				fmt.Printf("get new connection failed\n")
				continue
			}
			fmt.Printf("new connection with id %v\n", conn.id)

			// 开始
			if !conn.Start() {
				fmt.Printf("connection %v start failed\n", conn.id)
				break
			}
		}

		// 处理断开的连接
		this.conn_mgr.HandleDisconnList()

		time.Sleep(time.Millisecond)
	}

	this.state = TCP_SERVER_STATE_CLOSED
	return true
}
