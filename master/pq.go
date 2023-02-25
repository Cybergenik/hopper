package master

import (
	"container/heap"
	"fmt"
    c "github.com/Cybergenik/hopper/common"
)

// An Item is something we manage in a priority queue.
type Item struct {
    value    c.HashID
	priority float32
	index int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
// Max heap
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority > pq[j].priority }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := PriorityQueue{}
	heap.Init(&pq)
	// Insert a new item and then modify its priority.
	item1 := &Item{
		value:    321321,
		priority: 3,
	}
	item2 := &Item{
		value:    22121,
		priority: 10,
	}
	item3 := &Item{
		value:    2199918,
		priority: 7,
	}
	heap.Push(&pq, item1)
	heap.Push(&pq, item2)
	heap.Push(&pq, item3)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%d \n", item.priority, item.value)
	}
}
