package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// sync.Condの使い途

	// cond.Waitのソースを読むと, Lock～Wait(), Wait()～Unlock()間はいずれもクリティカルセクション
	// なので次のようにかいても競合しないはず

	l := new(sync.Mutex)
	c := sync.NewCond(l)
	var counter int

	for i := 0; i < 100; i++ {
		go func() {
			c.L.Lock()

			// ここはクリティカルセクション
			// なので他のスレッドと競合しない
			counter++

			// クリティカルセクション中に待機
			// 	 Wait()中にL.Unlock()されるので, 別スレッドがlのロック取れるようになる
			//   Wait()最後で再度Lock()とる
			c.Wait()

			// ここもクリティカルセクション
			// なので他のスレッドと競合しない
			counter++

			c.L.Unlock()
		}()
	}
	for i := 3; i > 0; i-- {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
	c.Broadcast()
	time.Sleep(time.Second * 3)
	fmt.Printf("counter=%d\n", counter)
}
