package libjmy

import (
	"container/list"
	"fmt"
	_ "log"
	"net"
)

const (
	CONN_STATE_NOT_CONNECT = iota
	CONN_STATE_CONNECTING
	CONN_STATE_CONNECTED
	CONN_STATE_DISCONNECTING
	CONN_STATE_DISCONNECTED
)

const (
	CONN_TYPE_ACTIVE = iota
	CONN_TYPE_PASSIVE
)

type JmyTcpConnection struct {
	conn      net.Conn
	conn_id   uint32
	state     int
	conn_type int
	recv_list *list.List
}

func (this *JmyTcpConnection) init(conn_type int) {
	this.conn_type = conn_type
	this.state = CONN_STATE_NOT_CONNECT
}

func (this *JmyTcpConnection) close() {

}

func (this *JmyTcpConnection) force_close() {

}

func (this *JmyTcpConnection) create_from(conn net.Conn) bool {
	this.conn = conn
	this.state = CONN_STATE_CONNECTED
	return true
}

func (this *JmyTcpConnection) connect(addr string) bool {
	if this.conn == nil {
		return false
	}
	var err error
	this.conn, err = net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial addr(%v) failed", addr)
		return false
	}
	this.state = CONN_STATE_CONNECTED
	return true
}

func (this *JmyTcpConnection) start() bool {
	if this.conn == nil {
		return false
	}
	if this.state != CONN_STATE_CONNECTED {
		return false
	}
	go this.read_loop()
	go this.write_loop()
	return true
}

func (this *JmyTcpConnection) handle_packet_type(pack_type int) int {
	if pack_type == JMY_PACKET_DISCONNECT {

	} else if pack_type == JMY_PACKET_DISCONNECT_ACK {

	} else if pack_type == JMY_PACKET_HEARTBEAT {

	} else if pack_type == JMY_PACKET_USER_DATA {
		return 1
	} else if pack_type == JMY_PACKET_USER_ID_DATA {
		return 1
	}
	return 0
}

func (this *JmyTcpConnection) read_loop() {
	if this.conn == nil {
		return
	}
	this.recv_list = list.New()
	const head_len = int(JMY_PACKET_LEN_HEAD + JMY_PACKET_LEN_TYPE)
	var tmp_head = [head_len]byte{}
	for {
		// read packet head
		read_len, err := this.conn.Read(tmp_head[:head_len])
		if err != nil {
			this.force_close()
			fmt.Printf("read packet head error: %v", err)
			return
		}
		if read_len != head_len {
			this.force_close()
			fmt.Printf("read packet head len failed, read_len(%d) head_len(%d)", read_len, head_len)
			return
		}

		res, pack_type, pack_len := jmy_net_proto_get_packet_len_type(tmp_head[:head_len])
		if !res {
			this.force_close()
			fmt.Printf("get packet len type failed")
			return
		}

		if this.handle_packet_type(pack_type) == 1 {
			pack_buf := make([]byte, pack_len)
			// read packet
			read_len, err = this.conn.Read(pack_buf)
			if err != nil {
				this.force_close()
				fmt.Printf("receive packet failed, error: %v", err)
				return
			}
			this.recv_list.PushBack(pack_buf)
		}
	}
}

func (this *JmyTcpConnection) write_loop() {

}

type JmyTcpConnectionMgr struct {
	conn_array []*JmyTcpConnection
	max_conn   uint32
	free_list  *list.List
	used_map   map[uint32]uint32
}

func (this *JmyTcpConnectionMgr) init(max_conn uint32) {
	this.free_list = list.New()
	this.used_map = make(map[uint32]uint32)
	this.conn_array = make([]*JmyTcpConnection, max_conn)
	for i := uint32(0); i < max_conn; i++ {
		this.conn_array[i] = &JmyTcpConnection{}
		this.conn_array[i].conn_id = i + 1
		this.free_list.PushBack(i + 1)
	}
	this.max_conn = max_conn
}

func (this *JmyTcpConnectionMgr) getFree() *JmyTcpConnection {
	if this.free_list.Len() == 0 {
		fmt.Printf("get free connection failed, used out")
		return nil
	}
	e := this.free_list.Front()
	id := (e.Value).(uint32)
	this.free_list.Remove(e)
	this.used_map[id] = id
	return this.conn_array[id-1]
}

func (this *JmyTcpConnectionMgr) free(conn *JmyTcpConnection) bool {
	if conn == nil {
		return false
	}
	if _, ok := this.used_map[conn.conn_id]; !ok {
		return false
	}
	this.free_list.PushFront(conn.conn_id)
	delete(this.used_map, conn.conn_id)
	return true
}
