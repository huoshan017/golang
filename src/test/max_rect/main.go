package main

import (
	"fmt"
)

type stack struct {
	top      int32 // top小于等于0时表示没有数据
	capacity int32
	data     []int32
	inited   bool
}

func new_stack(size int32) (s *stack) {
	if size <= 0 {
		return
	}
	s = &stack{}
	s.top = 0
	s.capacity = size
	s.data = make([]int32, size)
	s.inited = true
	return
}

func (s *stack) is_inited() bool {
	return s.inited
}

func (s *stack) is_empty() bool {
	if s.top > 0 {
		return false
	}
	return true
}

func (s *stack) is_full() bool {
	if s.top == 0 {
		return false
	}
	return true
}

func (s *stack) get_top() (bool, int32) {
	if s.top <= 0 {
		return false, 0
	}
	return true, s.data[s.top]
}

func (s *stack) get_bottom() (bool, int32) {
	if s.top <= 0 {
		return false, 0
	}
	return true, s.data[0]
}

func (s *stack) get_data(index int32) int32 {
	if index < 0 || index >= s.top {
		return -1
	}

	return s.data[index]
}

func (s *stack) push(d int32) bool {
	if s.top >= s.capacity {
		return false
	}

	s.data[s.top] = d
	s.top += 1
	return true
}

func (s *stack) pop() (bool, int32) {
	if s.top <= 0 {
		return false, 0
	}

	d := s.data[s.top-1]
	s.top -= 1
	return true, d
}

func main() {
	size := int32(2)
	s := new_stack(size)
	if s == nil {
		fmt.Printf("create new stack with size %d failed\n", size)
		return
	}
	s.push(1)
	s.push(2)
	for {
		b, a := s.pop()
		if !b {
			fmt.Printf("pop failed\n")
			break
		}
		fmt.Printf("pop %v\n", a)
	}
}
