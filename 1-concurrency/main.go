package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Нужно создать приложение с 2-мя горутинами, где:
// Первая создаёт slice из 10 случайных элементов от 0 до 100 и передаёт их по одному во вторую горутину.
// Вторая получает числа от 1-й и возводит в квадрат передавая результат в main.
// В main дожидаемся всех 10 чисел, которые были возведены в квадрат и выводим их в консоль.

func main() {
	numbCount := 10
	squareCh := make(chan int, numbCount)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		numbSlice := make([]int, numbCount)
		for i := range numbCount {
			numbSlice[i] = getRandomNumber()
		}
		for _, v := range numbSlice {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				squareCh <- getSquare(n)
			}(v)
		}
	}()

	go func() {
		wg.Wait()
		close(squareCh)
	}()

	for r := range squareCh {
		output(r)
	}
}

// func main() {
// 	numbCount := 10
// 	numbCh := make(chan int, numbCount)
// 	squareCh := make(chan int, numbCount)

// 	go func() {
// 		for range numbCount {
// 			numbCh <- getRandomNumber()
// 		}
// 		close(numbCh)
// 	}()

// 	go func() {
// 		for v := range numbCh {
// 			squareCh <- getSquare(v)
// 		}
// 		close(squareCh)
// 	}()

// 	for r := range squareCh {
// 		output(r)
// 	}
// }

func getRandomNumber() int {
	return rand.Intn(101)
}

func getSquare(n int) int {
	return n * n
}

func output(n int) {
	fmt.Println(n)
}
