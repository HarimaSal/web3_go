package main

import (
	"fmt"
)

func main() {
	/* 数字 */
	fmt.Println("------两数之和")
	//task1()
	fmt.Println("------只出现一次的数字")
	//task2()
	fmt.Println("------回文数")
	//task3()
	/*  字符串 */
	fmt.Println("------有效的括号")
	//task4()
	fmt.Println("------最长公共前缀")
	//task5()
	/* 基本值类型 */
	fmt.Println("------大整数加1")
	//task6()
	/* 引用类型 */
	fmt.Println("------删除有序数组中的重复项")
	task7()
	fmt.Println("------合并区间")
	task8()
}

/* 两数之和
 */
func task1() {
	arr := []int{2, 11, 15, 7}
	target := 9
	twoSum := func(nums []int, target int) []int {
		res := map[int]int{}
		for i := 0; i < len(nums); i++ {
			v, exist := res[target-nums[i]]
			if exist {
				return []int{v, i}
			} else {
				res[nums[i]] = i
			}
		}
		return nil
	}
	fmt.Println(twoSum(arr, target))
}

/* 只出现一次的数字
 */
func task2() {
	singleNumber := func(nums []int) int {
		tmpMap := map[int]int{}
		for i := 0; i < len(nums); i++ {
			tmpMap[nums[i]]++
		}
		for _, v := range nums {
			if tmpMap[v] == 1 {
				return v
			}
		}
		return 0
	}
	targetArr := []int{2, 2, 9, 5, 5, 5, 5}
	fmt.Println(singleNumber(targetArr))
	//fmt.Println("----------")
	//fmt.Println("结果：", singleNumber2(targetArr))
}

/* 该异或算法就解题而言是最高效的，但适用性不强。只能在偶数次出现的数字时才可用（除了单个出现数字外）
 */
func singleNumber2(nums []int) int {
	res := 0
	for _, v := range nums {
		res ^= v
	}
	return res
}

/*	回文数
 */
func task3() {
	res := func(x int) bool {
		if x < 0 {
			return false
		}
		if x == 0 {
			return true
		}
		if x%10 == 0 {
			return false
		}

		origNum := x
		reverseNum := 0
		for x != 0 {
			digital := x % 10
			reverseNum = reverseNum*10 + digital
			x = x / 10
		}
		return origNum == reverseNum
	}
	fmt.Println(res(121))
	fmt.Println(res(1212))
}

/*有效的括号*/
func task4() {
	res := func(s string) bool {
		// 创建一个映射，存储括号的对应关系
		pairs := map[rune]rune{
			'(': ')',
			'[': ']',
			'{': '}',
		}

		var stack []rune

		for _, char := range s {
			if closing, isOpen := pairs[char]; isOpen {
				stack = append(stack, closing)
			} else {
				if len(stack) == 0 || stack[len(stack)-1] != char {
					return false
				}
				stack = stack[:len(stack)-1]
			}
		}

		return len(stack) == 0
	}
	fmt.Println(res("()"))
	fmt.Println(res("([)"))
	fmt.Println(res("([)]"))
	fmt.Println(res("([4566]){}"))
	fmt.Println(res("([]{}[])"))
}

/*最长公共前缀*/
func task5() {
	longestCommonPrefix := func(strs []string) string {
		cmpStr := strs[0]
		for i := range cmpStr {
			for _, v := range strs {
				if i == len(v) || cmpStr[i] != v[i] { // 前面的条件是为了防止索引越界，后面的条件是为了截取相同字符
					return cmpStr[:i]
				}
			}
		}
		return cmpStr
	}
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "fool"}))
	fmt.Println(longestCommonPrefix([]string{"ab", "a"}))
}

/* 大整数加1.不可转为整数，然后加1，这样会引发位数越界 */
func task6() {
	addOne := func(digits []int) []int {
		var res []int
		for i := len(digits) - 1; i >= 0; i-- {
			if digits[i] == 9 {
				digits[i] = 0
				res = append(res, digits[i])
			} else {
				digits[i]++
				return digits
			}
		}
		if len(res) != 0 {
			res = append([]int{1}, res...)
		}
		return res
	}
	fmt.Println(addOne([]int{1, 2, 3}))
	fmt.Println(addOne([]int{9, 9, 9}))
	fmt.Println(addOne([]int{7, 9, 6}))
	fmt.Println(addOne([]int{7, 2, 8, 5, 0, 9, 1, 2, 9, 5, 3, 6, 6, 7, 3, 2, 8, 4, 3, 7, 9, 5, 7, 7, 4, 7, 4, 9, 4, 7, 0, 1, 1, 1, 7, 4, 0, 0, 6}))
}

/* 删除有序数组中的重复项 */
func task7() {
	//removeDuplicates := func(nums []int) int {
	//
	//}
	//fmt.Println(removeDuplicates([]int{1, 1, 2}))
	//fmt.Println(removeDuplicates([]int{0,0,1,1,1,2,2,3,3,4}))
}

/* 合并区间 */
func task8() {

}
