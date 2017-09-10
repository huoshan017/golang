package network

type TcpClient struct {
	net_conn *TcpConnection
}

func (this *TcpClient) Init(conn_id uint32, conn_mgr *TcpConnectionMgr) bool {
	return this.net_conn.Init(TCP_CONNECTION_TYPE_ACTIVE, conn_id, conn_mgr)
}

func (this *TcpClient) Connect(addr string) bool {
	return this.net_conn.Connect(addr)
}

func (this *TcpClient) Start() bool {
	return this.net_conn.Start()
}

func (this *TcpClient) Send(data []byte) error {
	return this.net_conn.Send(data)
}
