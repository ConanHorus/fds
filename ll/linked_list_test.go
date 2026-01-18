package ll

import (
	"testing"

	"github.com/smarty/assertions"
)

// --- Constructor and Basic State Tests ---

func TestNewLinkedList(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	and.So(list, assertions.ShouldNotBeNil)
	and.So(list.Length(), assertions.ShouldEqual, 0)
	and.So(list.IsEmpty(), assertions.ShouldBeTrue)
}

func TestNewLinkedListFromSlice(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		input          []int
		expectedLength int
		expectedSlice  []int
	}{
		"nil slice":      {input: nil, expectedLength: 0, expectedSlice: []int{}},
		"empty slice":    {input: []int{}, expectedLength: 0, expectedSlice: []int{}},
		"single element": {input: []int{42}, expectedLength: 1, expectedSlice: []int{42}},
		"multiple elements": {
			input:          []int{1, 2, 3, 4, 5},
			expectedLength: 5,
			expectedSlice:  []int{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)

			list := NewLinkedListFromSlice(tc.input)

			and.So(list, assertions.ShouldNotBeNil)
			and.So(list.Length(), assertions.ShouldEqual, tc.expectedLength)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedSlice)
		})
	}
}

func TestLength(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup          func() *LinkedList[int]
		expectedLength int
	}{
		"empty list": {
			setup:          func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedLength: 0,
		},
		"single element": {
			setup: func() *LinkedList[int] {
				list := NewLinkedList[int]()
				list.PushBack(1)
				return list
			},
			expectedLength: 1,
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				list := NewLinkedList[int]()
				for i := 0; i < 10; i++ {
					list.PushBack(i)
				}
				return list
			},
			expectedLength: 10,
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()
			and.So(list.Length(), assertions.ShouldEqual, tc.expectedLength)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup           func() *LinkedList[int]
		expectedIsEmpty bool
	}{
		"new list": {
			setup:           func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedIsEmpty: true,
		},
		"after push": {
			setup: func() *LinkedList[int] {
				list := NewLinkedList[int]()
				list.PushBack(1)
				return list
			},
			expectedIsEmpty: false,
		},
		"after push and pop": {
			setup: func() *LinkedList[int] {
				list := NewLinkedList[int]()
				list.PushBack(1)
				list.PopBack()
				return list
			},
			expectedIsEmpty: true,
		},
		"after clear": {
			setup: func() *LinkedList[int] {
				list := NewLinkedList[int]()
				list.PushBack(1)
				list.PushBack(2)
				list.Clear()
				return list
			},
			expectedIsEmpty: true,
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()
			and.So(list.IsEmpty(), assertions.ShouldEqual, tc.expectedIsEmpty)
		})
	}
}

// --- Access Method Tests ---

func TestGet(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		index         int
		expectedValue int
		expectedOK    bool
	}{
		"empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			index:         0,
			expectedValue: 0,
			expectedOK:    false,
		},
		"negative index": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         -1,
			expectedValue: 0,
			expectedOK:    false,
		},
		"index out of bounds": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         3,
			expectedValue: 0,
			expectedOK:    false,
		},
		"index way out of bounds": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         100,
			expectedValue: 0,
			expectedOK:    false,
		},
		"get first element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{10, 20, 30})
			},
			index:         0,
			expectedValue: 10,
			expectedOK:    true,
		},
		"get middle element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{10, 20, 30})
			},
			index:         1,
			expectedValue: 20,
			expectedOK:    true,
		},
		"get last element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{10, 20, 30})
			},
			index:         2,
			expectedValue: 30,
			expectedOK:    true,
		},
		"single element at index 0": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			index:         0,
			expectedValue: 42,
			expectedOK:    true,
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			value, ok := list.Get(tc.index)

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(value, assertions.ShouldEqual, tc.expectedValue)
		})
	}
}

// --- Push and Pop Tests ---

func TestPushBack(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	list.PushBack(1)
	and.So(list.Length(), assertions.ShouldEqual, 1)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1})

	list.PushBack(2)
	and.So(list.Length(), assertions.ShouldEqual, 2)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1, 2})

	list.PushBack(3)
	and.So(list.Length(), assertions.ShouldEqual, 3)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1, 2, 3})
}

func TestPushFront(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	list.PushFront(1)
	and.So(list.Length(), assertions.ShouldEqual, 1)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1})

	list.PushFront(2)
	and.So(list.Length(), assertions.ShouldEqual, 2)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{2, 1})

	list.PushFront(3)
	and.So(list.Length(), assertions.ShouldEqual, 3)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{3, 2, 1})
}

func TestPopBack(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup             func() *LinkedList[int]
		expectedValue     int
		expectedOK        bool
		expectedRemaining []int
	}{
		"empty list": {
			setup:             func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedValue:     0,
			expectedOK:        false,
			expectedRemaining: []int{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedValue:     42,
			expectedOK:        true,
			expectedRemaining: []int{},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			expectedValue:     3,
			expectedOK:        true,
			expectedRemaining: []int{1, 2},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			value, ok := list.PopBack()

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(value, assertions.ShouldEqual, tc.expectedValue)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedRemaining)
		})
	}
}

func TestPopFront(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup             func() *LinkedList[int]
		expectedValue     int
		expectedOK        bool
		expectedRemaining []int
	}{
		"empty list": {
			setup:             func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedValue:     0,
			expectedOK:        false,
			expectedRemaining: []int{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedValue:     42,
			expectedOK:        true,
			expectedRemaining: []int{},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			expectedValue:     1,
			expectedOK:        true,
			expectedRemaining: []int{2, 3},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			value, ok := list.PopFront()

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(value, assertions.ShouldEqual, tc.expectedValue)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedRemaining)
		})
	}
}

func TestPopBackSequence(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})

	for expected := 5; expected >= 1; expected-- {
		value, ok := list.PopBack()
		and.So(ok, assertions.ShouldBeTrue)
		and.So(value, assertions.ShouldEqual, expected)
	}

	and.So(list.IsEmpty(), assertions.ShouldBeTrue)

	value, ok := list.PopBack()
	and.So(ok, assertions.ShouldBeFalse)
	and.So(value, assertions.ShouldEqual, 0)
}

func TestPopFrontSequence(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})

	for expected := 1; expected <= 5; expected++ {
		value, ok := list.PopFront()
		and.So(ok, assertions.ShouldBeTrue)
		and.So(value, assertions.ShouldEqual, expected)
	}

	and.So(list.IsEmpty(), assertions.ShouldBeTrue)

	value, ok := list.PopFront()
	and.So(ok, assertions.ShouldBeFalse)
	and.So(value, assertions.ShouldEqual, 0)
}

func TestInterleavedPushPop(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	list.PushBack(1)
	list.PushBack(2)
	list.PushFront(0)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 2})

	v, _ := list.PopFront()
	and.So(v, assertions.ShouldEqual, 0)

	v, _ = list.PopBack()
	and.So(v, assertions.ShouldEqual, 2)

	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1})

	list.PushFront(10)
	list.PushBack(20)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{10, 1, 20})
}

// --- Indexed Mutation Tests ---

func TestInsertAt(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		index         int
		value         int
		expectedOK    bool
		expectedSlice []int
	}{
		"insert into empty list at 0": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			index:         0,
			value:         42,
			expectedOK:    true,
			expectedSlice: []int{42},
		},
		"insert into empty list at invalid index": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			index:         1,
			value:         42,
			expectedOK:    false,
			expectedSlice: []int{},
		},
		"insert at beginning": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         0,
			value:         0,
			expectedOK:    true,
			expectedSlice: []int{0, 1, 2, 3},
		},
		"insert at middle": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         1,
			value:         99,
			expectedOK:    true,
			expectedSlice: []int{1, 99, 2, 3},
		},
		"insert at end (append)": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         3,
			value:         4,
			expectedOK:    true,
			expectedSlice: []int{1, 2, 3, 4},
		},
		"insert at negative index": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         -1,
			value:         99,
			expectedOK:    false,
			expectedSlice: []int{1, 2, 3},
		},
		"insert beyond length": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         4,
			value:         99,
			expectedOK:    false,
			expectedSlice: []int{1, 2, 3},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			ok := list.InsertAt(tc.index, tc.value)

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedSlice)
		})
	}
}

func TestRemoveAt(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup             func() *LinkedList[int]
		index             int
		expectedValue     int
		expectedOK        bool
		expectedRemaining []int
	}{
		"remove from empty list": {
			setup:             func() *LinkedList[int] { return NewLinkedList[int]() },
			index:             0,
			expectedValue:     0,
			expectedOK:        false,
			expectedRemaining: []int{},
		},
		"remove at negative index": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:             -1,
			expectedValue:     0,
			expectedOK:        false,
			expectedRemaining: []int{1, 2, 3},
		},
		"remove beyond length": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:             3,
			expectedValue:     0,
			expectedOK:        false,
			expectedRemaining: []int{1, 2, 3},
		},
		"remove first element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:             0,
			expectedValue:     1,
			expectedOK:        true,
			expectedRemaining: []int{2, 3},
		},
		"remove middle element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:             1,
			expectedValue:     2,
			expectedOK:        true,
			expectedRemaining: []int{1, 3},
		},
		"remove last element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:             2,
			expectedValue:     3,
			expectedOK:        true,
			expectedRemaining: []int{1, 2},
		},
		"remove only element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			index:             0,
			expectedValue:     42,
			expectedOK:        true,
			expectedRemaining: []int{},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			value, ok := list.RemoveAt(tc.index)

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(value, assertions.ShouldEqual, tc.expectedValue)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedRemaining)
		})
	}
}

func TestSet(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		index         int
		value         int
		expectedOK    bool
		expectedSlice []int
	}{
		"set at 0 in empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			index:         0,
			value:         42,
			expectedOK:    true,
			expectedSlice: []int{42},
		},
		"set in empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			index:         1,
			value:         42,
			expectedOK:    false,
			expectedSlice: []int{},
		},
		"set at negative index": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         -1,
			value:         99,
			expectedOK:    true,
			expectedSlice: []int{99, 1, 2, 3},
		},
		"set at very negative index": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         -2,
			value:         99,
			expectedOK:    false,
			expectedSlice: []int{1, 2, 3},
		},
		"set just beyond length": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         3,
			value:         99,
			expectedOK:    true,
			expectedSlice: []int{1, 2, 3, 99},
		},
		"set beyond length": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         4,
			value:         99,
			expectedOK:    false,
			expectedSlice: []int{1, 2, 3},
		},
		"set first element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         0,
			value:         10,
			expectedOK:    true,
			expectedSlice: []int{10, 2, 3},
		},
		"set middle element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         1,
			value:         20,
			expectedOK:    true,
			expectedSlice: []int{1, 20, 3},
		},
		"set last element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			index:         2,
			value:         30,
			expectedOK:    true,
			expectedSlice: []int{1, 2, 30},
		},
		"set only element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1})
			},
			index:         0,
			value:         99,
			expectedOK:    true,
			expectedSlice: []int{99},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			ok := list.Set(tc.index, tc.value)

			and.So(ok, assertions.ShouldEqual, tc.expectedOK)
			and.So(list.ToSlice(), assertions.ShouldResemble, tc.expectedSlice)
		})
	}
}

// --- Search Tests ---

func TestContains(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup          func() *LinkedList[int]
		value          int
		expectedResult bool
	}{
		"empty list": {
			setup:          func() *LinkedList[int] { return NewLinkedList[int]() },
			value:          42,
			expectedResult: false,
		},
		"value not present": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:          42,
			expectedResult: false,
		},
		"value at beginning": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:          1,
			expectedResult: true,
		},
		"value in middle": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:          2,
			expectedResult: true,
		},
		"value at end": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:          3,
			expectedResult: true,
		},
		"duplicate values": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 2, 3})
			},
			value:          2,
			expectedResult: true,
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			result := list.Contains(tc.value)

			and.So(result, assertions.ShouldEqual, tc.expectedResult)
		})
	}
}

func TestIndexOf(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		value         int
		expectedIndex int
	}{
		"empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			value:         42,
			expectedIndex: -1,
		},
		"value not present": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:         42,
			expectedIndex: -1,
		},
		"value at beginning": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:         1,
			expectedIndex: 0,
		},
		"value in middle": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:         2,
			expectedIndex: 1,
		},
		"value at end": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			value:         3,
			expectedIndex: 2,
		},
		"duplicate values returns first": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 2, 3})
			},
			value:         2,
			expectedIndex: 1,
		},
		"multiple duplicates": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{5, 5, 5, 5})
			},
			value:         5,
			expectedIndex: 0,
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			index := list.IndexOf(tc.value)

			and.So(index, assertions.ShouldEqual, tc.expectedIndex)
		})
	}
}

// --- Bulk Operation Tests ---

func TestClear(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup func() *LinkedList[int]
	}{
		"clear empty list": {
			setup: func() *LinkedList[int] { return NewLinkedList[int]() },
		},
		"clear single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
		},
		"clear multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})
			},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			list.Clear()

			and.So(list.IsEmpty(), assertions.ShouldBeTrue)
			and.So(list.Length(), assertions.ShouldEqual, 0)
			and.So(list.ToSlice(), assertions.ShouldResemble, []int{})
		})
	}
}

func TestToSlice(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		expectedSlice []int
	}{
		"empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedSlice: []int{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedSlice: []int{42},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})
			},
			expectedSlice: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			slice := list.ToSlice()

			and.So(slice, assertions.ShouldResemble, tc.expectedSlice)
		})
	}
}

func TestToSliceDoesNotShareMemory(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3})
	slice := list.ToSlice()

	// Modify the slice
	slice[0] = 999

	// Original list should be unchanged
	value, _ := list.Get(0)
	and.So(value, assertions.ShouldEqual, 1)
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	original := []int{10, 20, 30, 40, 50}
	list := NewLinkedListFromSlice(original)
	result := list.ToSlice()

	and.So(result, assertions.ShouldResemble, original)
}

// --- Iteration Tests ---

func TestForEach(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup          func() *LinkedList[int]
		expectedValues []int
	}{
		"empty list": {
			setup:          func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedValues: []int{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedValues: []int{42},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})
			},
			expectedValues: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			collected := []int{}
			list.ForEach(func(v int) bool {
				collected = append(collected, v)
				return true
			})

			and.So(collected, assertions.ShouldResemble, tc.expectedValues)
		})
	}
}

func TestForEachIndexed(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		expectedPairs []struct{ index, value int }
	}{
		"empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedPairs: []struct{ index, value int }{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedPairs: []struct{ index, value int }{{0, 42}},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{10, 20, 30})
			},
			expectedPairs: []struct{ index, value int }{{0, 10}, {1, 20}, {2, 30}},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			collected := []struct{ index, value int }{}
			list.ForEachIndexed(func(i int, v int) bool {
				collected = append(collected, struct{ index, value int }{i, v})
				return true
			})

			and.So(collected, assertions.ShouldResemble, tc.expectedPairs)
		})
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup          func() *LinkedList[int]
		expectedValues []int
	}{
		"empty list": {
			setup:          func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedValues: []int{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedValues: []int{42},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})
			},
			expectedValues: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			collected := []int{}
			for v := range list.All() {
				collected = append(collected, v)
			}

			and.So(collected, assertions.ShouldResemble, tc.expectedValues)
		})
	}
}

func TestAllIndexed(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup         func() *LinkedList[int]
		expectedPairs []struct{ index, value int }
	}{
		"empty list": {
			setup:         func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedPairs: []struct{ index, value int }{},
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedPairs: []struct{ index, value int }{{0, 42}},
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{10, 20, 30})
			},
			expectedPairs: []struct{ index, value int }{{0, 10}, {1, 20}, {2, 30}},
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			collected := []struct{ index, value int }{}
			for i, v := range list.AllIndexed() {
				collected = append(collected, struct{ index, value int }{i, v})
			}

			and.So(collected, assertions.ShouldResemble, tc.expectedPairs)
		})
	}
}

func TestAllEarlyBreak(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3, 4, 5})

	collected := []int{}
	for v := range list.All() {
		collected = append(collected, v)
		if v == 3 {
			break
		}
	}

	and.So(collected, assertions.ShouldResemble, []int{1, 2, 3})
}

// --- String Representation Tests ---

func TestString(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		setup          func() *LinkedList[int]
		expectedString string
	}{
		"empty list": {
			setup:          func() *LinkedList[int] { return NewLinkedList[int]() },
			expectedString: "[]",
		},
		"single element": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{42})
			},
			expectedString: "[42]",
		},
		"multiple elements": {
			setup: func() *LinkedList[int] {
				return NewLinkedListFromSlice([]int{1, 2, 3})
			},
			expectedString: "[1, 2, 3]",
		},
	}

	for name, tc := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			list := tc.setup()

			result := list.String()

			and.So(result, assertions.ShouldEqual, tc.expectedString)
		})
	}
}

// --- Edge Cases and Stress Tests ---

func TestLargeList(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	const size = 1000
	list := NewLinkedList[int]()

	// Push all elements
	for i := 0; i < size; i++ {
		list.PushBack(i)
	}

	and.So(list.Length(), assertions.ShouldEqual, size)

	// Verify random access
	for _, idx := range []int{0, 100, 500, 999} {
		v, ok := list.Get(idx)
		and.So(ok, assertions.ShouldBeTrue)
		and.So(v, assertions.ShouldEqual, idx)
	}

	// Pop all elements
	for i := 0; i < size; i++ {
		v, ok := list.PopFront()
		and.So(ok, assertions.ShouldBeTrue)
		and.So(v, assertions.ShouldEqual, i)
	}

	and.So(list.IsEmpty(), assertions.ShouldBeTrue)
}

func TestAlternatingPushFrontBack(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	// Alternating pushes: front, back, front, back...
	list.PushFront(2) // [2]
	list.PushBack(3)  // [2, 3]
	list.PushFront(1) // [1, 2, 3]
	list.PushBack(4)  // [1, 2, 3, 4]
	list.PushFront(0) // [0, 1, 2, 3, 4]

	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 2, 3, 4})
}

func TestRepeatedClearAndReuse(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	for round := 0; round < 3; round++ {
		// Fill
		for i := 0; i < 10; i++ {
			list.PushBack(i + round*10)
		}
		and.So(list.Length(), assertions.ShouldEqual, 10)

		// Clear
		list.Clear()
		and.So(list.IsEmpty(), assertions.ShouldBeTrue)
	}

	// Final use
	list.PushBack(999)
	v, ok := list.Get(0)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 999)
}

func TestInsertAtAllPositions(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	// Start with [0, 2, 4]
	list := NewLinkedListFromSlice([]int{0, 2, 4})

	// Insert 1 at index 1: [0, 1, 2, 4]
	ok := list.InsertAt(1, 1)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 2, 4})

	// Insert 3 at index 3: [0, 1, 2, 3, 4]
	ok = list.InsertAt(3, 3)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 2, 3, 4})

	// Insert 5 at end: [0, 1, 2, 3, 4, 5]
	ok = list.InsertAt(5, 5)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 2, 3, 4, 5})

	// Insert -1 at beginning: [-1, 0, 1, 2, 3, 4, 5]
	ok = list.InsertAt(0, -1)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{-1, 0, 1, 2, 3, 4, 5})
}

func TestRemoveAtAllPositions(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	// Start with [0, 1, 2, 3, 4]
	list := NewLinkedListFromSlice([]int{0, 1, 2, 3, 4})

	// Remove middle
	v, ok := list.RemoveAt(2)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 2)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 3, 4})

	// Remove last
	v, ok = list.RemoveAt(3)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 4)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{0, 1, 3})

	// Remove first
	v, ok = list.RemoveAt(0)
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 0)
	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1, 3})
}

func TestMixedOperations(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedList[int]()

	// Complex sequence of operations
	list.PushBack(1)
	list.PushBack(2)
	list.PushFront(0)
	// [0, 1, 2]

	list.InsertAt(2, 99)
	// [0, 1, 99, 2]

	v, _ := list.PopFront()
	and.So(v, assertions.ShouldEqual, 0)
	// [1, 99, 2]

	list.Set(1, 100)
	// [1, 100, 2]

	and.So(list.Contains(100), assertions.ShouldBeTrue)
	and.So(list.IndexOf(100), assertions.ShouldEqual, 1)

	v, _ = list.RemoveAt(1)
	and.So(v, assertions.ShouldEqual, 100)
	// [1, 2]

	list.PushBack(3)
	// [1, 2, 3]

	and.So(list.ToSlice(), assertions.ShouldResemble, []int{1, 2, 3})
	and.So(list.Length(), assertions.ShouldEqual, 3)
}

func TestPopUntilEmpty(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3})

	// Pop from both ends alternately
	v, ok := list.PopFront()
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 1)

	v, ok = list.PopBack()
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 3)

	v, ok = list.PopFront()
	and.So(ok, assertions.ShouldBeTrue)
	and.So(v, assertions.ShouldEqual, 2)

	// Now empty
	and.So(list.IsEmpty(), assertions.ShouldBeTrue)

	_, ok = list.PopFront()
	and.So(ok, assertions.ShouldBeFalse)

	_, ok = list.PopBack()
	and.So(ok, assertions.ShouldBeFalse)
}

func TestContainsAfterModification(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{1, 2, 3})

	and.So(list.Contains(2), assertions.ShouldBeTrue)

	list.RemoveAt(1)

	and.So(list.Contains(2), assertions.ShouldBeFalse)
	and.So(list.Contains(1), assertions.ShouldBeTrue)
	and.So(list.Contains(3), assertions.ShouldBeTrue)
}

func TestIndexOfAfterModification(t *testing.T) {
	t.Parallel()
	and := assertions.New(t)

	list := NewLinkedListFromSlice([]int{10, 20, 30, 40})

	and.So(list.IndexOf(30), assertions.ShouldEqual, 2)

	list.RemoveAt(0)
	// [20, 30, 40]

	and.So(list.IndexOf(30), assertions.ShouldEqual, 1)
	and.So(list.IndexOf(10), assertions.ShouldEqual, -1)
}
