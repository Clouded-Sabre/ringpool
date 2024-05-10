package main

import (
	"fmt"
	"time"

	"github.com/Clouded-Sabre/ringpool/lib"
)

// ConcreteData represents concrete data type
type Payload struct {
	content []byte
}

func (p *Payload) SetContent(s string, element *lib.Element) {
	p.content = []byte(s)
	element.Length = len(s)
}

// Reset resets the content of the concrete data
func (p *Payload) Reset() {
	fmt.Print("")
}

// PrintContent prints the content of the concrete data
func (p *Payload) PrintContent(length int) {
	fmt.Println("Content:", string(p.content[:length]))
}

func main() {
	lib.Debug = true
	// Define a function that creates a new instance of ConcreteData
	newData := func(length int) lib.DataInterface {
		return &Payload{
			content: make([]byte, length),
		}
	}

	// Create a new RingPool instance
	pool := lib.NewRingPool(10, 100, newData)

	// Get an element from the pool
	element := pool.GetElement()
	fmt.Println("The number of available element is", pool.AvailableChunks())

	// Use the element
	element.Data.(*Payload).SetContent("Hohoho", element)
	element.Data.PrintContent(element.Length)

	// Return the element to the pool
	pool.ReturnElement(element)
	fmt.Println("The number of available element is", pool.AvailableChunks())

	time.Sleep(20 * time.Second)
}
