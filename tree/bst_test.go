package tree

import (
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

func TestBST_Insert(t *testing.T) {
	t.Parallel()
	tree := NewBST[int]()

	and := assertions.New(t)
	and.So(tree.Insert(5), should.BeTrue)
	and.So(tree.Insert(3), should.BeTrue)
	and.So(tree.Insert(7), should.BeTrue)
	and.So(tree.Insert(5), should.BeFalse)
	and.So(tree.Length(), should.Equal, 3)
}

func TestBST_Contains(t *testing.T) {
	t.Parallel()
	tree := NewBST[int]()
	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(7)

	and := assertions.New(t)
	and.So(tree.Contains(5), should.BeTrue)
	and.So(tree.Contains(3), should.BeTrue)
	and.So(tree.Contains(99), should.BeFalse)
}

func TestBST_Delete(t *testing.T) {
	t.Parallel()
	tree := NewBST[int]()
	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(7)
	tree.Insert(1)
	tree.Insert(4)

	and := assertions.New(t)
	and.So(tree.Delete(3), should.BeTrue)
	and.So(tree.Contains(3), should.BeFalse)
	and.So(tree.Length(), should.Equal, 4)
	and.So(tree.Delete(99), should.BeFalse)
}

func TestBST_MinMax(t *testing.T) {
	t.Parallel()
	tree := NewBST[int]()

	and := assertions.New(t)

	_, ok := tree.Min()
	and.So(ok, should.BeFalse)

	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(7)

	min, ok := tree.Min()
	and.So(ok, should.BeTrue)
	and.So(min, should.Equal, 3)

	maximum, ok := tree.Max()
	and.So(ok, should.BeTrue)
	and.So(maximum, should.Equal, 7)
}

func TestBST_Clear(t *testing.T) {
	t.Parallel()
	tree := NewBST[int]()
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	tree.Clear()

	and := assertions.New(t)
	and.So(tree.IsEmpty(), should.BeTrue)
	and.So(tree.Length(), should.Equal, 0)
}
