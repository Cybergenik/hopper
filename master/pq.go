package master

import (
	"container/heap"
	"fmt"
	c "github.com/Cybergenik/hopper/common"
)

const MAX = 5_000

// An Item is something we manage in a priority queue.
type PQItem struct {
	Seed     []byte
	Energy   float64
	Id       c.FTaskID
	priority float64
	index    int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*PQItem

func (pq PriorityQueue) Len() int { return len(pq) }

// Max heap
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].priority > pq[j].priority }

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// Changed to keep a fixed size PQ
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*PQItem)
	if n >= MAX {
		smallestIndex := 0
		smallest := 10.0
		for _, item := range *pq {
			if item.priority < smallest {
				smallest = item.priority
				smallestIndex = item.index
			}
		}
		if item.priority > smallest {
			sItem := (*pq)[smallestIndex]
			sItem.priority = item.priority
			sItem.Energy = item.Energy
			sItem.Seed = item.Seed
			(*pq)[smallestIndex] = sItem
			heap.Fix(pq, smallestIndex)
		}
	} else {
		item.index = n
		*pq = append(*pq, item)
	}
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
	item1 := &PQItem{
		Id:       321321,
		priority: 3,
	}
	item2 := &PQItem{
		Id:       22121,
		priority: 10,
	}
	item3 := &PQItem{
		Id:       2199918,
		priority: 7,
	}
	heap.Push(&pq, item1)
	heap.Push(&pq, item2)
	heap.Push(&pq, item3)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PQItem)
		fmt.Printf("%.2f:%d \n", item.priority, item.Id)
	}
}
