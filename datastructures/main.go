package main

import (
	"datastructures/algos"
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println("-----------------------Linked Lists---------------------------")
	linkedList := algos.LinkedList{}
	for i := 0; i < 10; i++ {
		linkedList.Push(rand.Intn(1000))
	}
	printLinkedList(&linkedList)

	fmt.Println("--------------------------------------------------")

	fmt.Printf("index: %d\n", linkedList.GetIdx(318))

	fmt.Println("--------------------------------------------------")

	linkedList.Delete(81)
	printLinkedList(&linkedList)

	fmt.Println("--------------------------------------------------")

	fmt.Printf("data: %d\n", linkedList.Get(2))

	fmt.Println("----------------------Hash Tables----------------------------")

	hashTable := algos.NewHashTable()
	hashTable.Set("Rod", "Oh babe Oh babe!")
	hashTable.Set("Ivancito", "Hello world from hash tables!")
	hashTable.Set("foo", "bar")
	hashTable.Set("Cod", "bla bla")

	result, ok := hashTable.Get("Rod")
	if !ok {
		fmt.Println("ERROR cannot get item from hash table")
	} else {
		fmt.Printf("result: %+v \n", result)
	}

	printHashTable(hashTable)
}

func printLinkedList(ll *algos.LinkedList) {
	p := ll.Head
	for p != nil {
		fmt.Printf(" -> %+v \n", p)
		p = p.Next
	}
}

func printHashTable(ht *algos.HashTable) {
	for idx, ll := range ht.Data {
		if ll != nil {
			fmt.Printf("at index: %d has linked list:\n", idx)
			printLinkedList(ll)
		}
	}
}
