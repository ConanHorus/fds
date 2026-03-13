package sorted

import (
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func TestBinarySearch_Table(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		slice         []int
		target        int
		expectedIndex int
		expectedOK    bool
	}{
		"nil slice":   {slice: nil, target: 5, expectedIndex: 0, expectedOK: false},
		"empty slice": {slice: []int{}, target: 5, expectedIndex: 0, expectedOK: false},

		"single found":        {slice: []int{5}, target: 5, expectedIndex: 0, expectedOK: true},
		"single missing low":  {slice: []int{5}, target: 4, expectedIndex: 0, expectedOK: false},
		"single missing high": {slice: []int{5}, target: 6, expectedIndex: 1, expectedOK: false},

		// ascending order examples
		"found in middle": {slice: []int{2, 4, 6, 8, 10}, target: 6, expectedIndex: 2, expectedOK: true},
		"found at start":  {slice: []int{2, 4, 6, 8, 10}, target: 2, expectedIndex: 0, expectedOK: true},
		"found at end":    {slice: []int{2, 4, 6, 8, 10}, target: 10, expectedIndex: 4, expectedOK: true},

		"missing between":          {slice: []int{2, 4, 6, 8, 10}, target: 7, expectedIndex: 3, expectedOK: false},
		"missing larger than all":  {slice: []int{2, 4, 6, 8, 10}, target: 11, expectedIndex: 5, expectedOK: false},
		"missing smaller than all": {slice: []int{2, 4, 6, 8, 10}, target: 1, expectedIndex: 0, expectedOK: false},

		// value runs (duplicates) — must return FIRST occurrence
		"run at end":       {slice: []int{1, 3, 5, 7, 7, 7}, target: 7, expectedIndex: 3, expectedOK: true},
		"run in middle":    {slice: []int{2, 5, 7, 7, 7, 9}, target: 7, expectedIndex: 2, expectedOK: true},
		"run at beginning": {slice: []int{5, 5, 5, 8, 9}, target: 5, expectedIndex: 0, expectedOK: true},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			index, ok := BinarySearch(testCase.slice, testCase.target)
			and.So(index, assertions.ShouldEqual, testCase.expectedIndex)
			and.So(ok, assertions.ShouldEqual, testCase.expectedOK)
		})
	}
}

func TestGallopingSearch_Table(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		slice         []int
		target        int
		expectedIndex int
		expectedOK    bool
	}{
		"nil slice":   {slice: nil, target: 5, expectedIndex: 0, expectedOK: false},
		"empty slice": {slice: []int{}, target: 5, expectedIndex: 0, expectedOK: false},

		"single found":        {slice: []int{5}, target: 5, expectedIndex: 0, expectedOK: true},
		"single missing low":  {slice: []int{5}, target: 4, expectedIndex: 0, expectedOK: false},
		"single missing high": {slice: []int{5}, target: 6, expectedIndex: 1, expectedOK: false},

		// ascending order examples
		"found in middle": {slice: []int{2, 4, 6, 8, 10}, target: 6, expectedIndex: 2, expectedOK: true},
		"found at start":  {slice: []int{2, 4, 6, 8, 10}, target: 2, expectedIndex: 0, expectedOK: true},
		"found at end":    {slice: []int{2, 4, 6, 8, 10}, target: 10, expectedIndex: 4, expectedOK: true},

		"missing between":          {slice: []int{2, 4, 6, 8, 10}, target: 7, expectedIndex: 3, expectedOK: false},
		"missing larger than all":  {slice: []int{2, 4, 6, 8, 10}, target: 11, expectedIndex: 5, expectedOK: false},
		"missing smaller than all": {slice: []int{2, 4, 6, 8, 10}, target: 1, expectedIndex: 0, expectedOK: false},

		// value runs (duplicates) — must return FIRST occurrence
		"run at end":       {slice: []int{1, 3, 5, 7, 7, 7}, target: 7, expectedIndex: 3, expectedOK: true},
		"run in middle":    {slice: []int{2, 5, 7, 7, 7, 9}, target: 7, expectedIndex: 2, expectedOK: true},
		"run at beginning": {slice: []int{5, 5, 5, 8, 9}, target: 5, expectedIndex: 0, expectedOK: true},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			index, ok := GallopingSearch(testCase.slice, testCase.target)
			and.So(index, assertions.ShouldEqual, testCase.expectedIndex)
			and.So(ok, assertions.ShouldEqual, testCase.expectedOK)
		})
	}
}

func TestInsert_Table(t *testing.T) {
	t.Parallel()

	testTable := map[string]struct {
		slice           []int
		value           int
		allowDuplicates bool
		expected        []int
	}{
		// Basic insertion cases
		"insert into nil slice":   {slice: nil, value: 5, allowDuplicates: true, expected: []int{5}},
		"insert into empty slice": {slice: []int{}, value: 5, allowDuplicates: true, expected: []int{5}},

		// Single element cases
		"insert before single element":                      {slice: []int{5}, value: 3, allowDuplicates: true, expected: []int{3, 5}},
		"insert after single element":                       {slice: []int{5}, value: 7, allowDuplicates: true, expected: []int{5, 7}},
		"insert equal to single element with duplicates":    {slice: []int{5}, value: 5, allowDuplicates: true, expected: []int{5, 5}},
		"insert equal to single element without duplicates": {slice: []int{5}, value: 5, allowDuplicates: false, expected: []int{5}},

		// Multiple element insertion cases
		"insert at beginning": {slice: []int{2, 4, 6, 8}, value: 1, allowDuplicates: true, expected: []int{1, 2, 4, 6, 8}},
		"insert at end":       {slice: []int{2, 4, 6, 8}, value: 9, allowDuplicates: true, expected: []int{2, 4, 6, 8, 9}},
		"insert in middle":    {slice: []int{2, 4, 6, 8}, value: 5, allowDuplicates: true, expected: []int{2, 4, 5, 6, 8}},
		"insert between":      {slice: []int{2, 4, 8, 10}, value: 6, allowDuplicates: true, expected: []int{2, 4, 6, 8, 10}},

		// Duplicate handling
		"insert duplicate at beginning with duplicates allowed": {slice: []int{2, 4, 6, 8}, value: 2, allowDuplicates: true, expected: []int{2, 2, 4, 6, 8}},
		"insert duplicate at beginning without duplicates":      {slice: []int{2, 4, 6, 8}, value: 2, allowDuplicates: false, expected: []int{2, 4, 6, 8}},
		"insert duplicate in middle with duplicates allowed":    {slice: []int{2, 4, 6, 8}, value: 6, allowDuplicates: true, expected: []int{2, 4, 6, 6, 8}},
		"insert duplicate in middle without duplicates":         {slice: []int{2, 4, 6, 8}, value: 6, allowDuplicates: false, expected: []int{2, 4, 6, 8}},
		"insert duplicate at end with duplicates allowed":       {slice: []int{2, 4, 6, 8}, value: 8, allowDuplicates: true, expected: []int{2, 4, 6, 8, 8}},
		"insert duplicate at end without duplicates":            {slice: []int{2, 4, 6, 8}, value: 8, allowDuplicates: false, expected: []int{2, 4, 6, 8}},

		// Complex duplicate scenarios
		"insert into existing duplicate run with duplicates":    {slice: []int{2, 5, 5, 5, 8}, value: 5, allowDuplicates: true, expected: []int{2, 5, 5, 5, 5, 8}},
		"insert into existing duplicate run without duplicates": {slice: []int{2, 5, 5, 5, 8}, value: 5, allowDuplicates: false, expected: []int{2, 5, 5, 5, 8}},
		"insert before duplicate run":                           {slice: []int{2, 5, 5, 5, 8}, value: 4, allowDuplicates: true, expected: []int{2, 4, 5, 5, 5, 8}},
		"insert after duplicate run":                            {slice: []int{2, 5, 5, 5, 8}, value: 6, allowDuplicates: true, expected: []int{2, 5, 5, 5, 6, 8}},

		// Edge cases with larger slices
		"insert into large slice beginning": {slice: []int{1, 3, 5, 7, 9, 11, 13, 15}, value: 0, allowDuplicates: true, expected: []int{0, 1, 3, 5, 7, 9, 11, 13, 15}},
		"insert into large slice end":       {slice: []int{1, 3, 5, 7, 9, 11, 13, 15}, value: 16, allowDuplicates: true, expected: []int{1, 3, 5, 7, 9, 11, 13, 15, 16}},
		"insert into large slice middle":    {slice: []int{1, 3, 5, 7, 9, 11, 13, 15}, value: 8, allowDuplicates: true, expected: []int{1, 3, 5, 7, 8, 9, 11, 13, 15}},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)

			// Make a copy of the input slice to avoid modifying the test case
			inputCopy := make([]int, len(testCase.slice))
			copy(inputCopy, testCase.slice)

			result := Insert(inputCopy, testCase.value, testCase.allowDuplicates)

			and.So(result, assertions.ShouldResemble, testCase.expected)

			// Verify the result is still sorted
			for i := 1; i < len(result); i++ {
				and.So(result[i], assertions.ShouldBeGreaterThanOrEqualTo, result[i-1])
			}
		})
	}
}

func BenchmarkBinarySearch(b *testing.B) {
	slice := make([]int, 1<<16)
	for i := range slice {
		slice[i] = i * 2
	}

	provider := providers.New1(func(int) {}).
		Add(0).
		Add(1234).
		Add(32768).
		Add(65534).
		Add(70000)

	benchy.New(b, options.Medium).
		RegisterBenchmark("BinarySearch", provider.WrapBenchmarkFunc(func(x int) {
			BinarySearch(slice, x)
		})).
		Run()
}

func BenchmarkGallopingSearch(b *testing.B) {
	slice := make([]int, 1<<24)
	for i := range slice {
		slice[i] = i * 2
	}

	provider := providers.New1(func(int) {}).
		Add(0).
		Add(1234).
		Add(32768).
		Add(65534).
		Add(70000).
		Add(((1 << 24) * 2) - 6)

	benchy.New(b, options.Medium).
		RegisterBenchmark("GallopingSearch", provider.WrapBenchmarkFunc(func(x int) {
			GallopingSearch(slice, x)
		})).
		Run()
}
