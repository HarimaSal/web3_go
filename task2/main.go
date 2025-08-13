package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	/* 【任务1】
	题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，
	然后在主函数中调用该函数并输出修改后的值。
	*/
	fmt.Println("任务1：整数指针---")
	task1()
	/* 【任务2】
	编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	*/
	fmt.Println("任务2：协程，并发---")
	task2()
	/* 【任务3】
	定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，
	实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
	*/
	fmt.Println("任务3：面向对象---")
	task3()
	/* 【任务4】
	编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数
	，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	*/
	fmt.Println("任务4：Channel---")
	task4()
	/* 【任务5】
	编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
	启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	*/
	fmt.Println("任务5：锁机制---")
	task5()
}

func task1() {
	addFunc := func(intPtr *int) {
		*intPtr += 10
	}
	var a = 10
	addFunc(&a)
	fmt.Println(a)
}

func task2() {
	//go printOdd()
	go func() {
		fmt.Println("打印奇数：")
		for i := 1; i <= 10; i += 2 {
			fmt.Println(i)
		}
	}()
	//go printEven()
	go func() {
		fmt.Println("打印偶数：")
		for i := 2; i <= 10; i += 2 {
			fmt.Println(i)
		}
	}()
	// 等待协程执行
	time.Sleep(1 * time.Second)
}

func printEven() {
	fmt.Println("打印偶数：")
	for i := 2; i <= 10; i += 2 {
		fmt.Println(i)
	}
}

func printOdd() {
	fmt.Println("打印奇数：")
	for i := 1; i <= 10; i += 2 {
		fmt.Println(i)
	}
}

func task3() {
	r := Rectangle{width: 5, height: 10}
	c := Circle{radius: 5}
	fmt.Println("方形：")
	fmt.Println("面积为：", r.Area())
	fmt.Println("周长为：", r.Perimeter())
	fmt.Println("圆形：")
	fmt.Println("面积为：", c.Area())
	fmt.Println("周长为：", c.Perimeter())
}

type Shape interface {
	Area() float64      // 面积
	Perimeter() float64 // 周长
}

type Rectangle struct {
	width  float64
	height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

func task4() {
	var ch = make(chan int) // 初始化
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	go func() {
		for {
			if i, ok := <-ch; ok {
				fmt.Println(i)
			} else {
				break
			}
		}
	}()

	time.Sleep(1 * time.Second)
}

var counter int
var mu sync.Mutex

/* 共享计数器 */
func increment() {
	mu.Lock()
	defer mu.Unlock()
	counter++
}

func task5() {
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				increment()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println(counter)
}
