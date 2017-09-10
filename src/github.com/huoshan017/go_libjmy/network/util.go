package network

func GetConnection(id uint32, conn_mgr *TcpConnectionMgr) *TcpConnection {
	return conn_mgr.GetConn(id)
}

func GetConnectionByRecvInfo(recv_info *RecvInfo) *TcpConnection {
	return recv_info.conn_mgr.GetConn(recv_info.conn_id)
}

func GetConnectionByErrorInfo(err_info *ErrorInfo) *TcpConnection {
	return err_info.conn_mgr.GetConn(err_info.conn_id)
}

func GetConnectionByEventInfo(event_info *EventInfo) *TcpConnection {
	return event_info.conn_mgr.GetConn(event_info.conn_id)
}
