package main

import (
	"fmt"
	"time"

	"github.com/Clouded-Sabre/ringpool/lib"
)

// ConcreteData represents concrete data type
type Payload struct {
	content []byte
	length  int
}

func (p *Payload) SetContent(s string) {
	p.content = []byte(s)
	p.length = len(s)
}

// Reset resets the content of the concrete data
func (p *Payload) Reset() {
	p.length = 0
}

// PrintContent prints the content of the concrete data
func (p *Payload) PrintContent() {
	fmt.Println("Content:", string(p.content[:p.length]))
}

func main() {
	lib.Debug = true
	// Define a function that creates a new instance of ConcreteData
	newData := func(params ...interface{}) lib.DataInterface {
		if len(params) != 1 {
			// Handle error: invalid number of parameters
			return nil
		}

		// Extract bufferLength from params
		bufferLength, ok := params[0].(int)
		if !ok {
			// Handle error: invalid type for bufferLength
			return nil
		}

		return &Payload{
			content: make([]byte, bufferLength),
		}
	}

	// Create a new RingPool instance
	pool := lib.NewRingPool(10, newData, 100)

	// Get an element from the pool
	element := pool.GetElement()
	fmt.Println("The number of available element is", pool.AvailableChunks())

	// Use the element
	element.Data.(*Payload).SetContent("Hohoho")
	element.Data.PrintContent()

	// Return the element to the pool
	pool.ReturnElement(element)
	fmt.Println("The number of available element is", pool.AvailableChunks())

	time.Sleep(20 * time.Second)
}
