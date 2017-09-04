package libjmy

import (
	"container/list"
)

const JMY_ID_START = 1

type JmyIdGenerator32 struct {
	max_id       uint32
	curr_id      uint32
	free_id_list *list.List
	used_id_map  map[uint32]uint32
}

func (this *JmyIdGenerator32) init(max_id uint32) {
	this.max_id = max_id
}

func (this *JmyIdGenerator32) get() uint32 {
	if this.free_id_list.Len() > 0 {
		e := this.free_id_list.Front()
		id := (e.Value).(uint32)
		this.free_id_list.Remove(e)
		this.used_id_map[id] = id
		return id
	} else {
		const MAX_ID_32 = 0xffffffff
		if this.max_id == 0 {
			this.max_id = MAX_ID_32
		}
		if this.curr_id >= this.max_id {
			return 0
		}
		this.curr_id += 1
		if _, ok := this.used_id_map[this.curr_id]; ok {
			return 0
		}
		return this.curr_id
	}
}

func (this *JmyIdGenerator32) free(id uint32) bool {
	if _, ok := this.used_id_map[id]; !ok {
		return false
	}
	delete(this.used_id_map, id)
	this.free_id_list.PushBack(id)
	return true
}

type JmyIdGenerator64 struct {
	max_id       uint64
	curr_id      uint64
	free_id_list *list.List
	used_id_map  map[uint64]uint64
}

func (this *JmyIdGenerator64) init(max_id uint64) {
	this.max_id = max_id
}

func (this *JmyIdGenerator64) get() uint64 {
	if this.free_id_list.Len() > 0 {
		e := this.free_id_list.Front()
		id := (e.Value).(uint64)
		this.free_id_list.Remove(e)
		this.used_id_map[id] = id
		return id
	} else {
		const MAX_ID_64 = 0xffffffffffffffff
		if this.max_id == 0 {
			this.max_id = MAX_ID_64
		}
		if this.curr_id >= this.max_id {
			return 0
		}
		this.curr_id += 1
		if _, ok := this.used_id_map[this.curr_id]; ok {
			return 0
		}
		return this.curr_id
	}
}

func (this *JmyIdGenerator64) free(id uint64) bool {
	if _, ok := this.used_id_map[id]; !ok {
		return false
	}
	delete(this.used_id_map, id)
	this.free_id_list.PushBack(id)
	return true
}
