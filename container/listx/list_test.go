package listx

import (
	"testing"
)

func TestList(t *testing.T) {
	l := New[int]()

	// 测试 PushFront 和 Len
	l.PushFront(1)
	if l.Len() != 1 {
		t.Errorf("Len() = %d; want 1", l.Len())
	}

	// 测试 Front
	if l.Front().Value != 1 {
		t.Errorf("Front().Value = %d; want 1", l.Front().Value)
	}

	// 测试 PushBack
	l.PushBack(2)
	if l.Len() != 2 {
		t.Errorf("Len() = %d; want 2", l.Len())
	}

	// 测试 Back
	if l.Back().Value != 2 {
		t.Errorf("Back().Value = %d; want 2", l.Back().Value)
	}

	// 测试 Remove
	removed := l.Remove(l.Front())
	if removed != 1 {
		t.Errorf("Remove() = %d; want 1", removed)
	}
	if l.Len() != 1 {
		t.Errorf("Len() = %d; want 1", l.Len())
	}

	// 测试 InsertBefore
	e := l.Front()
	newE := l.InsertBefore(3, e)
	if newE.Value != 3 {
		t.Errorf("InsertBefore() = %d; want 3", newE.Value)
	}
	if l.Len() != 2 {
		t.Errorf("Len() = %d; want 2", l.Len())
	}

	// 测试 InsertAfter
	newE = l.InsertAfter(4, e)
	if newE.Value != 4 {
		t.Errorf("InsertAfter() = %d; want 4", newE.Value)
	}
	if l.Len() != 3 {
		t.Errorf("Len() = %d; want 3", l.Len())
	}

	// 测试 MoveToFront
	l.MoveToFront(e)
	if l.Front().Value != 2 {
		t.Errorf("MoveToFront() Front().Value = %d; want 2", l.Front().Value)
	}

	// 测试 MoveToBack
	l.MoveToBack(e)
	if l.Back().Value != 2 {
		t.Errorf("MoveToBack() Back().Value = %d; want 2", l.Back().Value)
	}

	// 测试 MoveBefore
	newE = l.Front()
	l.MoveBefore(e, newE)
	if l.Front().Value != 2 {
		t.Errorf("MoveBefore() Front().Value = %d; want 2", l.Front().Value)
	}

	// 测试 MoveAfter
	newE = l.Back()
	l.MoveAfter(e, newE)
	if l.Back().Value != 2 {
		t.Errorf("MoveAfter() Back().Value = %d; want 2", l.Back().Value)
	}

	// 测试 PushBackList
	other := New[int]()
	other.PushFront(5)
	other.PushFront(6)
	l.PushBackList(other)
	if l.Len() != 5 {
		t.Errorf("PushBackList() Len() = %d; want 5", l.Len())
	}

	// 测试 PushFrontList
	another := New[int]()
	another.PushFront(7)
	another.PushFront(8)
	l.PushFrontList(another)
	if l.Len() != 7 {
		t.Errorf("PushFrontList() Len() = %d; want 7", l.Len())
	}
}


// 定义一个自定义结构体用于测试
type Person struct {
	Name string
	Age  int
}

func TestListWithDifferentTypes(t *testing.T) {
	// 测试 string 类型
	stringList := New[string]()
	stringList.PushFront("hello")
	stringList.PushBack("world")

	if stringList.Len() != 2 {
		t.Errorf("Len() = %d; want 2", stringList.Len())
	}

	if stringList.Front().Value != "hello" {
		t.Errorf("Front().Value = %s; want hello", stringList.Front().Value)
	}

	if stringList.Back().Value != "world" {
		t.Errorf("Back().Value = %s; want world", stringList.Back().Value)
	}

	// 测试自定义结构体类型
	personList := New[Person]()
	personList.PushFront(Person{Name: "Alice", Age: 30})
	personList.PushBack(Person{Name: "Bob", Age: 25})

	if personList.Len() != 2 {
		t.Errorf("Len() = %d; want 2", personList.Len())
	}

	if personList.Front().Value.Name != "Alice" {
		t.Errorf("Front().Value.Name = %s; want Alice", personList.Front().Value.Name)
	}

	if personList.Back().Value.Age != 25 {
		t.Errorf("Back().Value.Age = %d; want 25", personList.Back().Value.Age)
	}
}