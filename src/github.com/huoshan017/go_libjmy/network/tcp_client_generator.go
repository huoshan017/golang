package network

type TcpClientGenerator struct {
	conn_mgr *TcpConnectionMgr
	clients  map[*TcpClient]*TcpClient
}

func (this *TcpClientGenerator) Init(max_count uint32) bool {
	this.clients = make(map[*TcpClient]*TcpClient, max_count)
	this.conn_mgr = &TcpConnectionMgr{}
	return this.conn_mgr.Init(TCP_CONNECTION_TYPE_ACTIVE, max_count)
}

func (this *TcpClientGenerator) NewClient(processor IConnProcessor) *TcpClient {
	conn := this.conn_mgr.NewConn3(processor)
	if conn == nil {
		return nil
	}
	c := &TcpClient{}
	this.clients[c] = c
	c.net_conn = conn
	return c
}

func (this *TcpClientGenerator) FreeClient(c *TcpClient) bool {
	if _, o := this.clients[c]; !o {
		return false
	}
	delete(this.clients, c)
	this.conn_mgr.FreeConn(c.net_conn.id)
	return true
}
