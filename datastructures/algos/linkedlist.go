package algos

type Node struct {
	Index int
	Data  any
	Next  *Node
}

type LinkedList struct {
	Head *Node
}

func (ll *LinkedList) Push(d any) {
	newNode := &Node{Data: d}

	if ll.Head == nil {
		ll.Head = newNode
	} else {
		node := ll.Head
		for node.Next != nil {
			node = node.Next
		}
		newNode.Index = node.Index + 1
		node.Next = newNode
	}
}

func (ll *LinkedList) GetIdx(data any) int {
	node := ll.Head
	for node.Next != nil {
		if node.Data == data {
			return node.Index
		}
		node = node.Next
	}

	return -1
}

func (ll *LinkedList) Get(idx int) any {
	node := ll.Head
	for node.Next != nil {
		if node.Index == idx {
			return node.Data
		}
		node = node.Next
	}

	return nil
}

func (ll *LinkedList) Delete(data any) {
	node := ll.Head
	for node.Next != nil {
		if node.Data == data {
			next := node.Next
			node.Index = next.Index - 1
			node.Data = next.Data
			node.Next = next.Next
		}
		node = node.Next
	}
}
