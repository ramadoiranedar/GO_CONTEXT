package go_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)

	// mine, err := myTst(12, 12)
	// fmt.Println(mine, err)
}

// func myTst(satu interface{}, dua interface{}) (interface{}, interface{}) {
// 	return true, false
// }

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextC, "g", "G")

	fmt.Println("contextA :", contextA)
	fmt.Println("contextB :", contextB)
	fmt.Println("contextC :", contextC)
	fmt.Println("contextD :", contextD)
	fmt.Println("contextE :", contextE)
	fmt.Println("contextF :", contextF)
	fmt.Println("contextG :", contextG)

	// IN THIS PARENT CHILD GET VALUE CONTEXT BEHAVIOR
	// GET VALUE FROM BOTTOM TO TOP, CANNOT GET VALUE FROM TOP TO BOTTOM
	fmt.Println("contextA Value :", contextA.Value("b"))

	fmt.Println("contextF Value :", contextF.Value("e"))
	fmt.Println("contextE Value :", contextE.Value("d"))
	fmt.Println("contextD Value :", contextD.Value("c"))
	fmt.Println("contextC Value :", contextC.Value("b"))
	fmt.Println("contextB Value :", contextB.Value("a"))

	fmt.Println("contextF Value :", contextF.Value("e"))
	fmt.Println("contextF Value :", contextF.Value("f"))
	fmt.Println("contextF Value :", contextF.Value("c"))
}

func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
	ctxParent := context.Background()
	ctx, cancel := context.WithCancel(ctxParent)

	destination := CreateCounter(ctx)
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel() // mengirim signal Cancel to CONTEXT
	time.Sleep(5 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	ctxParent := context.Background()
	ctx, cancel := context.WithTimeout(ctxParent, 5*time.Second)
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	ctxParent := context.Background()
	ctx, cancel := context.WithDeadline(ctxParent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
