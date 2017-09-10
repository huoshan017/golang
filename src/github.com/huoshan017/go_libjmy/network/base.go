package network

// 接收信息
type RecvInfo struct {
	conn_id  uint32
	conn_mgr *TcpConnectionMgr
	data     []byte
	data_len uint16
}

func (this *RecvInfo) init(conn_id uint32, conn_mgr *TcpConnectionMgr) {
	this.conn_id = conn_id
	this.conn_mgr = conn_mgr
}

func (this *RecvInfo) GetConnId() uint32 {
	return this.conn_id
}

func (this *RecvInfo) GetConnMgr() *TcpConnectionMgr {
	return this.conn_mgr
}

func (this *RecvInfo) GetData() []byte {
	return this.data
}

func (this *RecvInfo) GetDataLen() uint16 {
	return this.data_len
}

// 错误信息
type ErrorInfo struct {
	conn_id  uint32
	conn_mgr *TcpConnectionMgr
	err      error
}

func (this *ErrorInfo) init(conn_id uint32, conn_mgr *TcpConnectionMgr, err error) {
	this.conn_id = conn_id
	this.conn_mgr = conn_mgr
	this.err = err
}

func (this *ErrorInfo) GetConnId() uint32 {
	return this.conn_id
}

func (this *ErrorInfo) GetConnMgr() *TcpConnectionMgr {
	return this.conn_mgr
}

func (this *ErrorInfo) GetErr() error {
	return this.err
}

// 事件
type EventInfo struct {
	conn_id  uint32
	conn_mgr *TcpConnectionMgr
}

func (this *EventInfo) GetConnId() uint32 {
	return this.conn_id
}

func (this *EventInfo) GetConnMgr() *TcpConnectionMgr {
	return this.conn_mgr
}

func (this *EventInfo) init(conn_id uint32, conn_mgr *TcpConnectionMgr) {
	this.conn_id = conn_id
	this.conn_mgr = conn_mgr
}

// 连接处理器接口
type IConnProcessor interface {
	OnRecv(recv_info *RecvInfo) bool
	OnError(err_info *ErrorInfo)
	OnConnect(event_info *EventInfo)
	OnDisconnect(event_info *EventInfo)
}
