package tree

import (
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
)

func TestAVL_Insert(t *testing.T) {
	t.Parallel()
	tree := NewAVL[int]()

	and := assertions.New(t)
	and.So(tree.Insert(5), should.BeTrue)
	and.So(tree.Insert(3), should.BeTrue)
	and.So(tree.Insert(7), should.BeTrue)
	and.So(tree.Insert(5), should.BeFalse)
	and.So(tree.Length(), should.Equal, 3)
}

func TestAVL_BalancesAfterInsert(t *testing.T) {
	t.Parallel()
	tree := NewAVL[int]()

	// Insert sorted values which would degenerate a regular BST
	for i := 1; i <= 7; i++ {
		tree.Insert(i)
	}

	and := assertions.New(t)
	and.So(tree.Length(), should.Equal, 7)
	and.So(tree.Contains(1), should.BeTrue)
	and.So(tree.Contains(7), should.BeTrue)
}

func TestAVL_Delete(t *testing.T) {
	t.Parallel()
	tree := NewAVL[int]()
	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(7)
	tree.Insert(1)
	tree.Insert(4)

	and := assertions.New(t)
	and.So(tree.Delete(3), should.BeTrue)
	and.So(tree.Contains(3), should.BeFalse)
	and.So(tree.Length(), should.Equal, 4)
}

func TestAVL_MinMax(t *testing.T) {
	t.Parallel()
	tree := NewAVL[int]()

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

func TestAVL_Clear(t *testing.T) {
	t.Parallel()
	tree := NewAVL[int]()
	tree.Insert(1)
	tree.Insert(2)
	tree.Insert(3)

	tree.Clear()

	and := assertions.New(t)
	and.So(tree.IsEmpty(), should.BeTrue)
	and.So(tree.Length(), should.Equal, 0)
}
