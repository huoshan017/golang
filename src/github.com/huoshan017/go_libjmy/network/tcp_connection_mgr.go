package network

import (
	"container/list"
	"fmt"
	"net"
)

type TcpConnectionMgr struct {
	conn_type    int
	conns        []*TcpConnection // 连接数组
	id_gen       *IdGenerator     // id生成器
	free_indexes *list.List       // 空闲索引队列
	id2idx       map[uint32]int   // id到索引的映射
	max_conn     uint32           // 最大连接数
	curr_conn    uint32           // 当前连接数
	disconn_list chan uint32      // 断连队列
}

func (this *TcpConnectionMgr) Init(conn_type int, max_conn uint32) bool {
	this.id_gen = &IdGenerator{}
	this.id_gen.Init(uint32(DEFAULT_START_ID), max_conn)
	this.conn_type = conn_type
	this.conns = make([]*TcpConnection, max_conn)
	this.free_indexes = list.New()
	this.id2idx = make(map[uint32]int, max_conn)
	for i := 0; i < int(max_conn); i++ {
		this.conns[i] = &TcpConnection{}
		id := this.id_gen.Get()
		if id == 0 {
			return false
		}
		this.conns[i].Init(conn_type, id, this)
		this.free_indexes.PushBack(i)
		this.id2idx[id] = i
	}
	this.max_conn = max_conn
	this.conn_type = conn_type
	this.disconn_list = make(chan uint32, max_conn)
	return true
}

func (this *TcpConnectionMgr) close() {

}

func (this *TcpConnectionMgr) NewConn() *TcpConnection {
	if this.free_indexes.Len() <= 0 {
		fmt.Printf("get new free index failed\n")
		return nil
	}

	e := this.free_indexes.Front()
	index := (e.Value).(int)
	this.free_indexes.Remove(e)

	this.curr_conn += 1
	c := this.conns[index]
	return c
}

func (this *TcpConnectionMgr) NewConn2(conn net.Conn, processor IConnProcessor) *TcpConnection {
	c := this.NewConn()
	c.SetConnProcessor(conn, processor)
	return c
}

func (this *TcpConnectionMgr) NewConn3(processor IConnProcessor) *TcpConnection {
	c := this.NewConn()
	c.processor = processor
	return c
}

func (this *TcpConnectionMgr) GetConn(id uint32) *TcpConnection {
	if !this.id_gen.IsUsed(id) {
		fmt.Printf("connection %v is not used, get failed\n", id)
		return nil
	}
	idx, o := this.id2idx[id]
	if !o {
		return nil
	}
	return this.conns[idx]
}

func (this *TcpConnectionMgr) FreeConn(id uint32) bool {
	if !this.id_gen.IsUsed(id) {
		fmt.Printf("connection %v is not used, free failed\n", id)
		return false
	}
	idx, o := this.id2idx[id]
	if !o {
		return false
	}
	this.free_indexes.PushFront(idx)
	this.conns[idx].Close()
	//this.id_gen.Free(id) // 不需要释放
	if this.curr_conn > 0 {
		this.curr_conn -= 1
	}
	fmt.Printf("free connection %v\n", id)
	return true
}

func (this *TcpConnectionMgr) PushDisconnId(id uint32) {
	this.disconn_list <- id
}

func (this *TcpConnectionMgr) HandleDisconnList() {
	is_break := false
	for {
		if is_break {
			break
		}
		select {
		case id := <-this.disconn_list:
			{
				this.FreeConn(id)
			}
		default:
			{
				is_break = true
				break
			}
		}
	}
	//fmt.Printf("handled disconn list\n")
}
