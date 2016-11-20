package main

//"fmt"
//"log"
//"time"

func partition(arr []int, left, right int) (pivot int) {
	if left >= right {
		return
	}
	i := left
	j := right
	pivot = left
	for i < j {
		if arr[i] > arr[pivot] && arr[j] < arr[pivot] {
			arr[i], arr[j] = arr[j], arr[i]
		}
		for i <= right && arr[i] <= arr[pivot] {
			i += 1
		}
		for j >= left && arr[j] >= arr[pivot] {
			j -= 1
		}
	}
	if i-1 >= 0 {
		arr[pivot], arr[i-1] = arr[i-1], arr[pivot]
		pivot = i - 1
	}
	//fmt.Printf("i = %v  pivot = %v\n", i, pivot)
	return
}

func partition2(arr []int, left, right int) (pivot int) {
	if left >= right {
		return
	}
	pivot_value := arr[right] // 最后一个存放基准值
	pivot = left              // 基准位置初始化成第一个索引
	for i := left; i < right; i++ {
		if pivot_value > arr[i] {
			if pivot != i {
				arr[i], arr[pivot] = arr[pivot], arr[i]
			}
			pivot += 1
		}
	}
	arr[pivot], arr[right] = arr[right], arr[pivot]
	return
}

func Qsort(t int, arr []int, left, right int) {
	if left >= right {
		return
	}

	pivot := 0
	if t == 0 {
		pivot = partition(arr, left, right)
	} else {
		pivot = partition2(arr, left, right)
	}
	Qsort(t, arr, left, pivot-1)
	Qsort(t, arr, pivot+1, right)
}
