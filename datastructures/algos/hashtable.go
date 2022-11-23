package algos

import (
	"fmt"
	"math"
)

const arrayLenght = 1000

type HashTable struct {
	Data [arrayLenght]*LinkedList
}

type listData struct {
	key   string
	value any
}

func NewHashTable() *HashTable {
	return &HashTable{
		[arrayLenght]*LinkedList{},
	}
}

func (h *HashTable) Set(k string, v any) *HashTable {
	index := index(hash(k))

	if h.Data[index] == nil {
		h.Data[index] = &LinkedList{}
		h.Data[index].Push(listData{k, v})
	} else {
		node := h.Data[index].Head
		for {
			if node != nil {
				d := node.Data.(listData)
				if d.key == k {
					d.value = v
					break
				}
			} else {
				h.Data[index].Push(listData{k, v})
				break
			}
			node = node.Next
		}
	}
	return h
}

func (h *HashTable) Get(k string) (result any, ok bool) {
	index := index(hash(k))
	linkedList := h.Data[index]

	if linkedList == nil {
		return "", false
	}
	node := linkedList.Head
	for {
		if node != nil {
			d := node.Data.(listData)
			if d.key == k {
				return d.value, true
			}
		} else {
			return "", false
		}
		node = node.Next
	}
}

func hash(key string) int {
	hash := 0
	for pos, char := range key {
		fmt.Println(int(char))
		fmt.Println(len(key))
		fmt.Println(pos)
		fmt.Println("--------")
		hash += int(char) * int(math.Pow(31, float64(len(key)-pos+1)))
	}

	return hash
}

func index(hash int) int {
	return hash % arrayLenght
}
