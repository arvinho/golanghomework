package main

import "errors"

//删除指定下标元素，支持泛型
func DeleteSliceByIndex[T any](slice []T, idx int) ([]T, error) {
	if len(slice) <= 0 {
		return nil, errors.New("切片长度出错")
	} else if idx < 0 || idx >= len(slice) {
		return nil, errors.New("下标出错")
	}
	temp := 0
	for i, val := range slice {
		if i != idx {
			slice[temp] = val
			temp++
		}
	}
	return slice[:temp], nil
}

//删除指定下标元素，支持泛型并缩容
func DeleteSlice[T any](slice []T, idx int) ([]T, error) {
	if len(slice) <= 0 {
		return nil, errors.New("切片长度出错")
	} else if idx < 0 || idx >= len(slice) {
		return nil, errors.New("下标出错")
	}
	temp := slice[:0]
	for i, val := range slice {
		if i != idx {
			temp = append(temp, val)
		}
	}
	temp1 := make([]T, len(slice)-1, cap(slice)-1)
	copy(temp1, temp)
	return temp1, nil
}
