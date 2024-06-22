package lib

import (
	"fmt"
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
	if pos >= 0 && pos < len(e.footprints) {
		e.footprints[pos] = e.footprints[pos] + "✓"
	}
}

// AddToChannel assign a channel string to StayAtChannel of the chunk
func (e *Element) AddChannel(channelString string) int {
	return e.AddFootPrint("(" + channelString + ")")
}

// PrintCallStack prints the call stack of the chunk
func (e *Element) PrintCallStack() {
	fmt.Print("Call Stack:")
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
