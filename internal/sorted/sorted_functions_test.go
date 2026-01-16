package sorted

import (
	"fmt"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func TestBinarySearchInt_Table(t *testing.T) {
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
			if name == "single missing high" {
				fmt.Println("debug")
			}

			and := assertions.New(t)
			index, ok := BinarySearchInt(testCase.slice, testCase.target)
			and.So(index, assertions.ShouldEqual, testCase.expectedIndex)
			and.So(ok, assertions.ShouldEqual, testCase.expectedOK)
		})
	}
}

func TestGallopingSearchInt_Table(t *testing.T) {
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
			index, ok := GallopingSearchInt(testCase.slice, testCase.target)
			and.So(index, assertions.ShouldEqual, testCase.expectedIndex)
			and.So(ok, assertions.ShouldEqual, testCase.expectedOK)
		})
	}
}

func BenchmarkBinarySearchInt(b *testing.B) {
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
		RegisterBenchmark("BinarySearchInt", provider.WrapBenchmarkFunc(func(x int) {
			BinarySearchInt(slice, x)
		})).
		Run()
}

func BenchmarkGallopingSearchInt(b *testing.B) {
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
		RegisterBenchmark("GallopingSearchInt", provider.WrapBenchmarkFunc(func(x int) {
			GallopingSearchInt(slice, x)
		})).
		Run()
}
