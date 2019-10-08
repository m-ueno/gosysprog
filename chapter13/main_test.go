package main_test

import (
	"fmt"
	"sync"
	"time"
)

func ExampleWithoutMutex() {
	// Mutexとらないと競合する

	var n int

	for i := 0; i < 1000; i++ {
		go func() {
			n++
		}()
	}
	time.Sleep(time.Second)
	fmt.Println(n)
	// Output: 918
}

func ExampleMutex() {
	// Mutexとると競合しない

	var n int
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		go func() {
			mu.Lock()
			n++
			mu.Unlock()
		}()
	}
	time.Sleep(time.Second * 3)
	fmt.Println(n)
	// Output: 1000
}

// 13.2 Go言語の並列処理のための道具
// goroutineの起動よりもループ変数が回るのが早いので
// println()が実行される頃にはtaskは"taskC"になってしまう
func ExampleGoroutinePassingByCaptureNG() {
	tasks := []string{"taskA", "taskB", "taskC"}
	for _, task := range tasks {
		go func() { fmt.Println(task) }()
	}

	time.Sleep(time.Second) // goroutine終了待ち

	// Output:
	// taskC
	// taskC
	// taskC
}

func ExampleGoroutinePassingArgsWait() {
	// 上の改良 goroutineの引数わたしにする
	tasks := []string{"taskA", "taskB", "taskC"}
	for i, task := range tasks {
		go func(n int, task string) {
			time.Sleep(time.Duration(n * 1000)) // 0, 1000, 2000 milliseconds
			fmt.Println(task)
		}(i, task)
	}

	time.Sleep(time.Second * 3) // goroutine終了待ち

	// Output:
	// taskA
	// taskB
	// taskC
}
func ExampleConditionCritical() {
	// cond.Waitのソースを読むと, Lock～Wait(), Wait()～Unlock()間はいずれもクリティカルセクション
	// なので次のようにかいても競合しない

	l := new(sync.Mutex)
	c := sync.NewCond(l)
	N := 10

	var counter int

	for i := 0; i < N; i++ {
		i := i
		go func() {
			c.L.Lock()
			defer c.L.Unlock()
			fmt.Printf("id=%d\n", i)
			// ここはクリティカルセクション
			counter++
			c.Wait()

			// ここもクリティカルセクション
			counter++
		}()
	}
	time.Sleep(time.Second)
	c.Broadcast()
	time.Sleep(time.Second * 3)
	// Output: 20
}

// 13.7.4 sync.Condの練習
// スレッドセーフなFIFOキューをchannelを使わずに実装してみる
type Q struct {
	elements []int
	capacity int
	l        *sync.Mutex
	c        *sync.Cond
}

func NewQ(cap int) *Q {
	l := new(sync.Mutex)
	c := sync.NewCond(l)
	return &Q{
		[]int{},
		cap,
		l,
		c,
	}
}

func (q *Q) Queue(elem int) {
	q.l.Lock()
	defer q.l.Unlock()

	for len(q.elements) == q.capacity {
		fmt.Println("waiting nofull")
		q.c.Wait()
	}
	q.elements = append(q.elements, elem)
	q.c.Signal()
}

func (q *Q) Dequeue() int {
	q.l.Lock()
	defer q.l.Unlock()

	for len(q.elements) == 0 {
		fmt.Println("waiting noempty")
		q.c.Wait()
	}
	a := q.elements
	x, a := a[len(a)-1], a[:len(a)-1]
	q.elements = a // ?
	q.c.Signal()

	return x
}

func ExampleQ() {
	q := NewQ(5)
	for i := 0; i < 10; i++ {
		i := i
		go q.Queue(i)
		go func() {
			time.Sleep(time.Millisecond * 100)
			fmt.Println(q.Dequeue())
		}()
	}
	time.Sleep(time.Second)
	// Output:
}
