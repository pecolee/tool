package tool

import (
	"github.com/samber/lo"
	"log"
	"reflect"
	"sort"
)

// IsContain 数组是否存在目标值
func IsContain(items interface{}, item interface{}) bool {
	if v, ok := items.([]string); ok {
		if vi, ok := item.(string); ok {
			return lo.Contains(v, vi)
		}
	}
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		log.Println("reflect kind is not slice")
		return false
	}

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 排序去重(字符串数组)
func RemoveDuplicates(arr []string) []string {
	sort.Strings(arr)
	result := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		if i == 0 || arr[i] != arr[i-1] {
			result = append(result, arr[i])
		}
	}
	return result
}
