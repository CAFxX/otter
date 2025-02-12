package s3fifo

import (
	"github.com/maypok86/otter/internal/node"
)

type small[K comparable, V any] struct {
	q       *node.Queue[K, V]
	main    *main[K, V]
	ghost   *ghost[K, V]
	cost    uint32
	maxCost uint32
}

func newSmall[K comparable, V any](
	maxCost uint32,
	main *main[K, V],
	ghost *ghost[K, V],
) *small[K, V] {
	return &small[K, V]{
		q:       node.NewQueue[K, V](),
		main:    main,
		ghost:   ghost,
		maxCost: maxCost,
	}
}

func (s *small[K, V]) insert(n *node.Node[K, V]) {
	s.q.Push(n)
	n.MarkSmall()
	s.cost += n.PolicyCost()
}

func (s *small[K, V]) evict(deleted []*node.Node[K, V]) []*node.Node[K, V] {
	if s.cost == 0 {
		return deleted
	}

	n := s.q.Pop()
	s.cost -= n.PolicyCost()
	n.Unmark()
	if n.IsExpired() {
		return append(deleted, n)
	}

	if n.Frequency() > 1 {
		s.main.insert(n)
		for s.main.isFull() {
			deleted = s.main.evict(deleted)
		}
		n.ResetFrequency()
		return deleted
	}

	return s.ghost.insert(deleted, n)
}

func (s *small[K, V]) remove(n *node.Node[K, V]) {
	s.cost -= n.PolicyCost()
	n.Unmark()
	s.q.Remove(n)
}

func (s *small[K, V]) length() int {
	return s.q.Len()
}

func (s *small[K, V]) clear() {
	s.q.Clear()
	s.cost = 0
}
