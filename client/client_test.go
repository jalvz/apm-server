package client

import (
	"sync"
	"testing"
	"fmt"
	"time"
)

func TestPauseResume (t *testing.T) {
	var wgTest sync.WaitGroup

	var wg sync.WaitGroup
	go func() {
		wgTest.Add(1)
		time.Sleep(time.Second * 2)
		fmt.Println("Pause 2 secs")
		wg.Add(1)
		time.Sleep(time.Second * 2)
		fmt.Println("Resume")
		wg.Done()

		time.Sleep(time.Second * 2)
		fmt.Println("Pause 4 secs")
		wg.Add(1)
		time.Sleep(time.Second * 4)
		fmt.Println("Resume")
		wg.Done()

		wgTest.Done()
	}()

	c := make(chan int, 20)
	go func () {
		wgTest.Add(1)

		for x := 0; x < 12; x++{
			time.Sleep(time.Second)
			c <- x
		}

		wgTest.Done()
	}()

	go func() {
		for i := range c {
			wg.Wait()
			fmt.Println("got ", i)
		}
	}()

	time.Sleep(time.Second )
	wgTest.Wait()

}

func TestClosure(t *testing.T) {
	fn := HealthStatusTracker()

	go func() {
		up := map[string]bool{"node1": true, "node2": true, "node3": false}
		down := map[string]bool{"node1": true, "node2": false, "node3": false}

		fn(up)
		time.Sleep(2 * time.Second)

		fn(down)
		time.Sleep(2 * time.Second)
		fn(down)
		time.Sleep(6 * time.Second)

		fn(up)
	}()

	time.Sleep(time.Second)
	fmt.Println("waiting")
	backPressure.Wait()
	fmt.Println("all good!")

	time.Sleep(5 * time.Second)
	fmt.Println("waiting again")
	backPressure.Wait()
	fmt.Println("really all good!")

}
