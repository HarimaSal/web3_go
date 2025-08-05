package main

import "fmt"

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
		splitStr := []rune(s)
		if len(splitStr)%2 != 0 {
			return false
		}
		stack := map[rune]rune{
			')': '(',
			']': '[',
			'}': '{',
		}
		k := 0
		for _, char := range splitStr {
			if char == '(' || char == '[' || char == '{' {
				k++
				continue
			}
			if k > 0 && splitStr[k-1] == stack[char] { // ([]{})
				k--
			} else {
				return false
			}
		}
		return true
	}
	fmt.Println(res("()"))
	fmt.Println(res("([)"))
	fmt.Println(res("([)]"))
	fmt.Println(res("([]){}"))
	fmt.Println(res("([]{}[])"))
}

/*最长公共前缀*/
func task5() {

}
