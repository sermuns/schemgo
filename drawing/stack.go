package drawing

type Stack struct {
	items []Point
}

func (s *Stack) Push(item Point) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() Point {
	itemNum := len(s.items)
	if itemNum == 0 {
		panic("Tried to pop from an empty stack!")
	}
	item := s.items[itemNum-1]
	s.items = s.items[:itemNum-1]
	return item
}
