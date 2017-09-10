package network

import (
	"container/list"
	"fmt"
)

var (
	DEFAULT_START_ID = 1
)

type IdGenerator struct {
	start_id  uint32
	free_list *list.List
	used_map  map[uint32]uint32
}

func (this *IdGenerator) Init(start_id uint32, max_count uint32) bool {
	this.start_id = start_id
	this.free_list = list.New()
	for i := start_id; i < max_count+start_id; i++ {
		this.free_list.PushBack(i)
	}
	fmt.Printf("free list len is %v\n", this.free_list.Len())
	this.used_map = make(map[uint32]uint32, max_count)
	return true
}

func (this *IdGenerator) IsUsed(id uint32) bool {
	if _, ok := this.used_map[id]; !ok {
		return false
	}
	return true
}

func (this *IdGenerator) GetStartId() uint32 {
	return this.start_id
}

func (this *IdGenerator) Clear() {
	for k, _ := range this.used_map {
		delete(this.used_map, k)
	}
	for {
		e := this.free_list.Back()
		if e == nil {
			break
		}
		this.free_list.Remove(e)
	}
}

func (this *IdGenerator) Reset() {

}

func (this *IdGenerator) Get() uint32 {
	var id uint32
	if this.free_list.Len() <= 0 {
		return 0
	}
	e := this.free_list.Front()
	this.free_list.Remove(e)
	id = (e.Value).(uint32)
	this.used_map[id] = id
	fmt.Printf("get new id %v\n", id)
	return id
}

func (this *IdGenerator) Free(id uint32) bool {
	if _, ok := this.used_map[id]; !ok {
		return false
	}
	delete(this.used_map, id)
	this.free_list.PushFront(id)
	fmt.Printf("free id %v\n", id)
	return true
}
