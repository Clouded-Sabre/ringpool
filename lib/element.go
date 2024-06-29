package lib

import (
	"fmt"
	"log"
	"time"
)

// Payload represents a single chunk of payload
type DataInterface interface {
	Reset()
	PrintContent()
}

type Element struct {
	Data DataInterface
	//Length         int // length of the real data
	index          int
	LastAllocation time.Time
	footprints     []string
	isAllocated    bool
}

type NewData func(params ...interface{}) DataInterface

// NewChunk creates a new chunk with the given length
func NewElement(index int, newData NewData, params ...interface{}) *Element {
	return &Element{
		Data:  newData(params...),
		index: index,
	}
}

// AddCallStack adds a function string to the call stack of the chunk
func (e *Element) AddFootPrint(funcString string) int {
	e.footprints = append(e.footprints, funcString)

	return len(e.footprints) - 1
}

// PopCallStack removes the last function string from the call stack of the chunk
func (e *Element) TickFootPrint(pos int) {
	if e.isAllocated {
		if pos >= 0 && pos < len(e.footprints) {
			e.footprints[pos] = e.footprints[pos] + "✓"
		}
	} else {
		log.Println("Element cannot be ticked because it is already been returned to pool.")
	}
}

// AddChannel assign a channel string to footprints of the chunk
func (e *Element) AddChannel(channelString string) int {
	return e.AddFootPrint("(" + channelString + ")")
}

// AddChannel assign a channel string to footprints of the chunk
func (e *Element) TickChannel() error {
	if e.isAllocated {
		if len(e.footprints) > 0 {
			pos := len(e.footprints) - 1 // the last one
			e.footprints[pos] = e.footprints[pos] + "✓"
			return nil
		} else {
			return fmt.Errorf("tickChannel: there is no channel to be ticked")
		}
	} else {
		return fmt.Errorf("element cannot be ticked because it is already been returned to pool")
	}
}

// PrintCallStack prints the call stack of the chunk
func (e *Element) PrintCallStack() {
	fmt.Print("Footprints:")
	for _, call := range e.footprints {
		fmt.Printf(" -> %s", call)
	}
	fmt.Println()
	e.Data.PrintContent()
}

// Reset resets necessary fields after chunk is returned
func (e *Element) Reset() {
	e.LastAllocation = time.Time{}
	e.footprints = nil
	e.isAllocated = false
	e.Data.Reset()
}

// IsTimedOut checks if an allocated chunk is been held for too long
func (e *Element) IsTimedOut(timeout time.Duration) bool {
	return time.Since(e.LastAllocation) > timeout
}

// GetAgeSeconds returns how long has it been since the chunk was last allocated
func (e *Element) GetAgeSeconds() float64 {
	return time.Since(e.LastAllocation).Seconds()
}
