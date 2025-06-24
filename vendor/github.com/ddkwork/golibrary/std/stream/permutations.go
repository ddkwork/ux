package stream

import "encoding/binary"

// Permute 递归回溯法实现全排列 https://www.cnblogs.com/xwxz/p/14812448.html
func Permute[T comparable](data []T) [][]T {
	var res [][]T
	var track []T
	used := make([]bool, len(data))

	var backtrack func()
	backtrack = func() {
		if len(track) == len(data) {
			// 添加当前track的副本
			temp := make([]T, len(track))
			copy(temp, track)
			res = append(res, temp)
			return
		}
		for i := range data {
			if used[i] {
				continue // 当前元素已经使用，跳过
			}
			track = append(track, data[i])
			used[i] = true               // 标记当前元素为已使用
			backtrack()                  // 递归调用
			track = track[:len(track)-1] // 移除最后一个元素，回溯
			used[i] = false              // 恢复标记
		}
	}

	backtrack()
	return res
}

func PermuteToUint32Slice(slice [][]byte) []uint32 {
	var uint32s []uint32
	for _, bytes := range slice { // permutations
		// 只处理长度为4的字节数组
		if len(bytes) == 4 {
			uint32s = append(uint32s, binary.LittleEndian.Uint32(bytes))
		}
	}
	return uint32s
}
