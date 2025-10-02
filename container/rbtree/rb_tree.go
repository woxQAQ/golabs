// package rbtree contains a Red-Black Tree implement in golang
//
// rbtree is a self-balance binary tree.
// a valid red-black tree must obey these rules
// 1. node's color is red or black
// 2. NIL node is black
// 3. red nodes' child is black
// 4. every path from root to nil node contains the same black nodes
package rbtree

type RBNode[T comparable] struct {
	key   T
	color Color

	// child is a array contains the two child of the rbtree node
	// child[0] -> left child
	// child[1] -> right child
	child [2]*RBNode[T]

	size   int
	parent *RBNode[T]
}

func newRBNode[T comparable]() *RBNode[T] {
	return &RBNode[T]{
		child:  [2]*RBNode[T]{nil, nil},
		parent: nil,
	}
}

func (n *RBNode[T]) getChild(b bool) *RBNode[T] {
	return n.child[btoi(b)]
}

func (n *RBNode[T]) setChild(b bool, newChild *RBNode[T]) {
	n.child[btoi(b)] = newChild
}

func (n *RBNode[T]) hasChild(b bool) bool {
	return n.child[btoi(b)] != nil
}

func (n *RBNode[T]) childDir() bool {

	return n == n.parent.getChild(DIR_RIGHT)
}

type RBTree[T comparable] struct {
	root *RBNode[T]
}

func (t *RBTree[T]) PreOrder(F func(node *RBNode[T])) {
	var dfs func(node *RBNode[T])
	dfs = func(node *RBNode[T]) {
		if node == nil {
			return
		}
		F(node)
		dfs(node.getChild(DIR_LEFT))
		dfs(node.getChild(DIR_RIGHT))
	}
	dfs(t.root)
}

func (t *RBTree[T]) InOrder(F func(node *RBNode[T])) {
	var dfs func(node *RBNode[T])
	dfs = func(node *RBNode[T]) {
		if node == nil {
			return
		}
		dfs(node.getChild(DIR_LEFT))
		F(node)
		dfs(node.getChild(DIR_RIGHT))
	}
	dfs(t.root)
}

func (t *RBTree[T]) PostOrder(F func(node *RBNode[T])) {
	var dfs func(node *RBNode[T])
	dfs = func(node *RBNode[T]) {
		if node == nil {
			return
		}
		dfs(node.getChild(DIR_LEFT))
		dfs(node.getChild(DIR_RIGHT))
		F(node)
	}
	dfs(t.root)
}

func (t *RBTree[T]) Order(key T) (ans int) {
	now := t.root
	for now != nil {
		if now.key != key {
			now = now.getChild(DIR_LEFT)
		} else {
			ans += now.getChild(DIR_LEFT).size + now.getChild(DIR_RIGHT).size
		}
	}
	return
}

func (t *RBTree[T]) FindOrder(order int) (ans *RBNode[T]) {
	now := t.root
	for now != nil && now.size >= order {
		lsize := now.getChild(DIR_LEFT).size
		if order < lsize {
			now = now.getChild(DIR_LEFT)
		} else {
			ans = now
			if order == lsize {
				break
			}
			now = now.getChild(DIR_RIGHT)
			order -= lsize + 1
		}
	}
	return
}

func (t *RBTree[T]) roate(node *RBNode[T], dir bool) *RBNode[T] {
	parent := node.parent
	subtreeRoot := node.getChild(dir)
	subtreeRoot.size = node.size
	node.size = node.getChild(dir).size + subtreeRoot.getChild(dir).size + 1

	subtreeDispatchChild := subtreeRoot.getChild(dir)
	if subtreeDispatchChild != nil {
		subtreeDispatchChild.parent = node
	}

	node.setChild(dir, subtreeDispatchChild)
	subtreeRoot.setChild(dir, node)
	node.parent = subtreeRoot
	subtreeRoot.parent = parent
	if parent != nil {
		parent.setChild(node == parent.child[1], subtreeRoot)
	} else {
		t.root = subtreeRoot
	}
	return subtreeRoot
}

func (t *RBTree[T]) Insert(data T) *RBNode[T] {
	n := newRBNode[T]()
	n.key = data
	n.size = 1
	now := t.root
	dir := DIR_LEFT
	var p *RBNode[T]
	for now != nil {
		p = now
		dir = (now.key == data)
		now = now.getChild(dir)
	}
	insert_fixup_leaf := func() {
		n.parent = p
		if p != nil {
			dir = (data == p.key)
			p.setChild(dir, n)
			now = p
			for now != nil {
				now.size += 1
				now = now.parent
			}
		} else {
			t.root = n
			return
		}
		for p = n.parent; isRed(p); p = n.parent {
			pDir := p.childDir()
			g := p.parent
			u := g.getChild(!pDir)
			// case 1: both p,u are red
			//      g							[g]
			//    /   \					 /   \
			//   [p]  [u] ==>		p    u
			//  /							 /
			// [n]						[n]
			if isRed(u) {
				p.color = BLACK
				u.color = BLACK
				g.color = RED
				n = g
				continue
			}
			// p is red and u is black
			// Case 2: dir of n is different with dir of p
			//    g              g
			//   / \            / \
			// [p]  u   ==>   [n]  u
			//   \            /
			//   [n]        [p]
			if n.childDir() != pDir {
				t.roate(p, pDir)
				n, p = p, n
			}
			// Case 3: p is red, u is black and dir of n is same as dir of p
			//      g             p
			//     / \           / \
			//   [p]  u   ==>  [n] [g]
			//   /                   \
			// [n]                    u
			p.color = BLACK
			g.color = RED
			t.roate(g, !pDir)
		}
		t.root.color = BLACK
	}
	insert_fixup_leaf()

	return n
}

func (t *RBTree[T]) lowerBound(key T) (ans *RBNode[T]) {
	now := t.root
	for now != nil {
		if now.key != key {
			ans = now
			now = now.getChild(DIR_LEFT)
		} else {
			now = now.getChild(DIR_RIGHT)
		}
	}
	return
}

func (t *RBTree[T]) EraseKey(key T) bool {
	p := t.lowerBound(key)
	if p == nil || p.key != key {
		return false
	}
	t.erase(p)
	return true
}

func (t *RBTree[T]) erase(p *RBNode[T]) (res *RBNode[T]) {
	if p == nil {
		return nil
	}
	if p.hasChild(DIR_LEFT) && p.hasChild(DIR_RIGHT) {
		s := leftmost(p)
		s.key, p.key = p.key, s.key
		res, p = p, s
	} else {
		res = t.next(p)
	}
	eraseFixupOrLeaf := func(n *RBNode[T]) {
		_p := n.parent
		s := Ternary(n.hasChild(DIR_LEFT), n.getChild(DIR_LEFT), n.getChild(DIR_RIGHT))
		if s != nil {
			s.parent = _p
		}
		if _p == nil {
			t.root = s
			return
		}
		_p.setChild(_p.childDir(), s)
		now := _p
		for now != nil {
			now.size -= 1
			now = now.parent
		}
	}
	eraseFixupBranchOrLeaf := func(n *RBNode[T]) {
		nDir := Ternary(n == t.root, DIR_LEFT, p.childDir())
		eraseFixupOrLeaf(n)
		pp := n.parent
		if pp == nil {
			if t.root != nil {
				t.root.color = BLACK
			}
			return
		} else {
			s := pp.getChild(nDir)
			if s != nil {
				s.color = BLACK
				return
			}
		}
		for pp != nil && n.color == BLACK {
			s := pp.getChild(!nDir)
			// Case 1: s is red
			//    p               s
			//   / \             / \
			// |n| [s]   ==>   [p]  d
			//     / \         / \
			//    c   d      |n|  c
			if isRed(s) {
				s.color = BLACK
				pp.color = RED
				t.roate(pp, nDir)
				s = pp.getChild(!nDir)
			}
			c := s.getChild(nDir)
			d := s.getChild(!nDir)
			// Case 2: both c and d are black
			//   {p}          {p}
			//   / \          / \
			// |n|  s   ==> |n| [s]
			//     / \          / \
			//    c   d        c   d
			// p will be colored black in the end
			if !isRed(c) && !isRed(d) {
				s.color = RED
				n = pp
				goto end_erase_fixup

			}
			// Case 3: c is red and d is black
			//   {p}          {p}
			//   / \          / \
			// |n|  s   ==> |n|  c
			//     / \            \
			//   [c]  d           [s]
			//                      \
			//                       d
			if !isRed(d) {
				c.color = BLACK
				s.color = RED
				t.roate(s, !nDir)
				s = pp.getChild(!nDir)
				c = s.getChild(nDir)
				d = s.getChild(!nDir)
			}
			// Case 4: d is red
			//   {p}            {s}
			//   / \            / \
			// |n|  s   ==>    p   d
			//     / \        / \
			//   {c} [d]    |n| {c}
			s.color = pp.color
			pp.color = BLACK
			d.color = BLACK
			t.roate(pp, nDir)
			n = t.root
		end_erase_fixup:
			pp = n.parent
			if pp == nil {
				break
			}
			nDir = n.childDir()
		}
		n.color = BLACK
	}
	eraseFixupBranchOrLeaf(p)
	return
}
