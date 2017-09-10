package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const DEFAULT_READ_BUFFER_SIZE = 4096

// 连接类型
const (
	TCP_CONNECTION_TYPE_PASSIVE = iota
	TCP_CONNECTION_TYPE_ACTIVE
)

const (
	TCP_CONN_STATE_NOT_CONNECT = iota
	TCP_CONN_STATE_CONNECTED
	TCP_CONN_STATE_DISCONNECTING
	TCP_CONN_STATE_DISCONNECTED
)

type TcpConnection struct {
	net_conn  net.Conn
	id        uint32
	mgr       *TcpConnectionMgr
	processor IConnProcessor

	send_chan  chan []byte
	conn_type  int
	conn_state int
}

func (this *TcpConnection) Init(conn_type int, id uint32, conn_mgr *TcpConnectionMgr) bool {
	this.conn_type = conn_type
	this.conn_state = TCP_CONN_STATE_NOT_CONNECT

	this.id = id
	this.mgr = conn_mgr

	this.send_chan = make(chan []byte, 16)

	fmt.Printf("connection %d inited\n", id)
	return true
}

func (this *TcpConnection) SetConnProcessor(conn net.Conn, processor IConnProcessor) {
	this.net_conn = conn
	this.conn_state = TCP_CONN_STATE_CONNECTED
	this.processor = processor

	var ei EventInfo
	ei.init(this.id, this.mgr)
	this.processor.OnConnect(&ei)
}

func (this *TcpConnection) Reset() {
	this.Close()
	this.send_chan = make(chan []byte, 16)
	fmt.Printf("reset connection %v\n", this.id)
}

func (this *TcpConnection) Close() {
	if this.conn_state == TCP_CONN_STATE_DISCONNECTED {
		return
	}
	this.net_conn.Close()
	_, ok := <-this.send_chan
	if ok {
		close(this.send_chan)
	}
	this.conn_state = TCP_CONN_STATE_DISCONNECTED
	fmt.Printf("close connection %v\n", this.id)
}

func (this *TcpConnection) Connect(addr string) bool {
	if this.conn_state != TCP_CONN_STATE_NOT_CONNECT &&
		this.conn_state != TCP_CONN_STATE_DISCONNECTED {
		return false
	}
	var err error
	this.net_conn, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("connect failed\n")
		return false
	}
	this.conn_state = TCP_CONN_STATE_CONNECTED

	var ei EventInfo
	ei.init(this.id, this.mgr)
	this.processor.OnConnect(&ei)

	fmt.Printf("connected success\n")
	return true
}

func (this *TcpConnection) Start() bool {
	if this.net_conn == nil || this.id == 0 || this.processor == nil {
		return false
	}

	buf := make([]byte, DEFAULT_READ_BUFFER_SIZE)
	var ri RecvInfo
	ri.init(this.id, this.mgr)
	ri.data = buf

	go func() {
		for {
			length, err := this.net_conn.Read(buf)
			if err == nil {
				if length > 0 {
					ri.data_len = uint16(length)
					if !this.processor.OnRecv(&ri) {
						fmt.Printf("OnRecv failed\n")
					}
					fmt.Printf("%s\n", string(buf[:length]))
				} else {
					fmt.Printf("read no data\n")
				}
			} else {
				if err != io.EOF {
					var ei ErrorInfo
					ei.init(this.id, this.mgr, err)
					this.processor.OnError(&ei)
				} else {
					var ei EventInfo
					ei.init(this.id, this.mgr)
					this.processor.OnDisconnect(&ei)
					fmt.Printf("peer(%v) closed connection\n", this.net_conn.RemoteAddr().String())
				}

				this.conn_state = TCP_CONN_STATE_DISCONNECTED

				// 送入断连队列中
				this.mgr.PushDisconnId(this.id)
				break
			}
		}
	}()

	go func() {
		for {
			c := <-this.send_chan
			if err := this.realSend(c); err != nil {
				var ei ErrorInfo
				ei.init(this.id, this.mgr, err)
				this.processor.OnError(&ei)
				break
			}
		}
	}()

	return true
}

func (this *TcpConnection) Send(data []byte) (err error) {
	select {
	case <-time.After(time.Second * 3):
		{
			fmt.Printf("send data timeout")
			err = errors.New("send data timeout")
		}
	case this.send_chan <- data:
		{

		}
	}
	return
}

func (this *TcpConnection) realSend(data []byte) (err error) {
	var length int
	for {
		var w int
		w, err = this.net_conn.Write(data[length:])
		if err != nil {
			fmt.Printf("write data failed, err(%v)", err.Error())
			break
		}
		length += w
		if length >= len(data) {
			break
		}
	}
	return
}
