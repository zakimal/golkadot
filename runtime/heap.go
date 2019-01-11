package runtime

import (
	"errors"
	"fmt"
	"math"
)

// Heap ...
type Heap struct {
	memory     *Memory
	wasmMemory *WasmMemory
}

// NewHeap ...
func NewHeap() *Heap {
	memory := createMemory(nil, 0)

	return &Heap{
		memory: memory,
	}
}

// WasResized ...
func (h *Heap) WasResized() bool {
	return h.memory.IsResized
}

// Allocate ...
func (h *Heap) Allocate(size int64) Pointer {
	if size == 0 {
		return 0
	}

	ptr := h.memory.Offset
	offset := ptr + size

	if offset < h.memory.Size {
		h.memory.Offset = offset
		h.memory.Allocated[ptr] = size

		return Pointer(ptr)
	}

	return h.FreeAlloc(size)
}

// Deallocate ...
func (h *Heap) Deallocate(ptr Pointer) (int64, error) {
	size, ok := h.memory.Allocated[int64(ptr)]
	if !ok {
		return 0, errors.New("Calling free() on unallocated memory")
	}

	delete(h.memory.Allocated, int64(ptr))

	h.memory.Deallocated[int64(ptr)] = size

	return size, nil
}

// Dup ...
func (h *Heap) Dup(ptr Pointer, length int64) []uint8 {
	return h.memory.Buffer[ptr : int64(ptr)+length]
}

// Fill ...
func (h *Heap) Fill(ptr Pointer, value uint8, length int64) []uint8 {
	for i := int64(ptr); i < int64(ptr)+length; i++ {
		h.memory.Buffer[i] = value
	}

	return h.memory.Buffer
}

// FreeAlloc ...
func (h *Heap) FreeAlloc(size int64) Pointer {
	ptr := h.FindContaining(size)

	if ptr == -1 {
		fmt.Printf("allocate(%d) failed, consider increasing the base runtime memory size\n", size)

		return h.GrowAlloc(size)
	}

	// TODO: being wasteful here so need to un-free the requested size instead of everything (long-term fragmentation and loss)
	delete(h.memory.Deallocated, int64(ptr))
	h.memory.Allocated[int64(ptr)] = int64(size)

	return ptr
}

// GrowAlloc ...
func (h *Heap) GrowAlloc(size int64) Pointer {
	// NOTE: grow memory by 4 times the requested amount (rounded up)
	if h.GrowMemory(1 + int64(math.Ceil(float64(4)*float64(size)/float64(PageSize)))) {
		return h.Allocate(size)
	}

	return 0
}

// Get ...
func (h *Heap) Get(ptr Pointer, length int64) []uint8 {
	return h.memory.Buffer[ptr : int64(ptr)+length]
}

// GetU32 ...
func (h *Heap) GetU32(ptr Pointer) []uint8 {
	return h.memory.Buffer[ptr:]
}

// Set ...
func (h *Heap) Set(ptr Pointer, data []uint8) Pointer {
	copy(h.memory.Buffer[ptr:], data)

	return ptr
}

// SetU32 ...
func (h *Heap) SetU32(ptr Pointer, value []uint8) Pointer {
	copy(h.memory.Buffer[ptr:], value)

	return ptr
}

// SetWASMMemory ...
func (h *Heap) SetWASMMemory(wasmMemory *WasmMemory, pageOffset int64) {
	if pageOffset == -1 {
		pageOffset = 4
	}

	offset := pageOffset * int64(PageSize)

	h.wasmMemory = wasmMemory
	h.memory = createMemory(wasmMemory.Buffer, offset)
}

// Size ...
func (h *Heap) Size() int64 {
	return h.memory.Size
}

// Used ...
func (h *Heap) Used() SizeUsed {
	return SizeUsed{
		Allocated:   h.CalculateSize(h.memory.Allocated),
		Deallocated: h.CalculateSize(h.memory.Deallocated),
	}
}

// CalculateSize ...
func (h *Heap) CalculateSize(buffer MemoryBuffer) int64 {
	var total int64
	for _, size := range buffer {
		total += size
	}

	return total
}

// GrowMemory ...
func (h *Heap) GrowMemory(pages int64) bool {
	if h.wasmMemory == nil {
		return false
	}

	// TODO
	//h.wasmMemory.grow(pages)
	h.memory.Size = int64(len(h.wasmMemory.Buffer))
	h.memory.Buffer = h.wasmMemory.Buffer
	h.memory.IsResized = true

	return true
}

// FindContaining ...
func (h *Heap) FindContaining(size int64) Pointer {
	var ptr int64 = -1

	for offset, size := range h.memory.Deallocated {
		if h.memory.Deallocated[offset] > size {
			continue
		}

		if size < ptr || ptr == -1 {
			ptr = offset
		}
	}

	return Pointer(ptr)
}

func createMemory(buffer []uint8, offset int64) *Memory {
	if buffer == nil {
		buffer = []uint8{}
	}

	if offset == -1 {
		offset = 256 * 1024
	}

	size := int64(len(buffer))
	// NOTE clear memory, it could be provided from a previous run
	for i := offset; i < size; i++ {
		buffer[i] = 0
	}

	return &Memory{
		Allocated:   nil,
		Deallocated: nil,
		IsResized:   false,
		Offset:      offset, // aligned with Rust (should have offset)
		Size:        size,
		Buffer:      buffer,
	}
}