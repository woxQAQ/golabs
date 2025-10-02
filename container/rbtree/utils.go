package rbtree

type Color uint8

const (
	BLACK Color = iota
	RED
)

const (
	DIR_LEFT  = false
	DIR_RIGHT = true
)

func btoi(b bool) int {
	return Ternary(b, 1, 0)
}

func isRed[T comparable](node *RBNode[T]) bool {
	if node == nil {
		return false
	}
	return Ternary(node == nil, false, node.color == RED)
}

func most[T comparable](p *RBNode[T], dir bool) *RBNode[T] {
	if p == nil {
		return nil
	}
	q := p
	for q.getChild(dir) != nil {
		q = q.getChild(dir)
	}
	return q
}

func leftmost[T comparable](p *RBNode[T]) *RBNode[T] {
	return most(p, DIR_LEFT)
}

func rightmost[T comparable](p *RBNode[T]) *RBNode[T] {
	return most(p, DIR_RIGHT)
}

func (t *RBTree[T]) neighbour(p *RBNode[T], dir bool) *RBNode[T] {
	if p == nil {
		return nil
	}
	if p.getChild(dir) != nil {
		return most(p.getChild(dir), dir)
	}
	if p == t.root {
		return nil
	}
	for p != nil && p.parent != nil && p.childDir() == dir {
		p = p.parent
	}
	return Ternary(p == nil, nil, p.parent)
}

func (t *RBTree[T]) prev(p *RBNode[T]) *RBNode[T] {
	return t.neighbour(p, DIR_LEFT)
}

func (t *RBTree[T]) next(p *RBNode[T]) *RBNode[T] {
	return t.neighbour(p, DIR_RIGHT)
}
