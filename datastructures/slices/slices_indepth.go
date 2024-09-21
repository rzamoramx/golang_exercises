package slices

import (
	"fmt"
	"unsafe"
)

func SlicesInDepth() {
	// initial slice
	a := make([]int, 4, 6)
	for i := 0; i < len(a); i++ {
		a[i] = i
	}

	// will be
	// a: [0, 1, 2, 3] len=4 cap=6
	// b: [1, 2] len=2 cap=5
	b := a[1:3]
	printSlice(1, a, b)

	// will be
	// a: [0, 1, 2, 8] len=4 cap=6
	// b: [1, 2, 8] len=3 cap=5
	b = append(b, 8)
	printSlice(2, a, b)

	// will be
	// a: [0, 1, 2, 8] len=4 cap=6
	// b: [1, 2] len=2 cap=2
	b = a[1:3:3]
	printSlice(3, a, b)

	// will be
	// a: [0, 1, 2, 8] len=4 cap=6
	// b: [1, 2, 9] len=3 cap=4
	b = append(b, 9)
	printSlice(4, a, b)

	// that's because the slices share the same underlying array, we can see it as a view of an array
}

func printSlice(printNum int, slices ...[]int) {
	fmt.Printf("Printing #: %d\n", printNum)
	for _, slice := range slices {
		toStr := fmt.Sprintf("%v", slice)
		pointer := unsafe.SliceData(slice)
		//fmt.Printf("Slice: %v, Len: %d, Cap: %d, Data: %p\n", str, len(s), cap(s), p)
		fmt.Printf("%-10slen=%d cap=%d addr+%v\n", toStr, len(slice), cap(slice), pointer)
	}
	fmt.Println()
}
