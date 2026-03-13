package ll

import (
	"testing"

	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func BenchmarkLinkedLists_Read(b *testing.B) {
	linkedList := NewLinkedList[int]()
	simpleLinkedList := NewSimpleLinkedList[int]()

	provider := providers.New1(func(int) {})

	numberOfElements := 10_000
	for i := 0; i < numberOfElements/2; i++ {
		linkedList.PushBack(i)
		simpleLinkedList.PushBack(i)
		provider.Add(i)
	}

	for i := 0; i < numberOfElements/2; i += numberOfElements / 10 {
		linkedList.RemoveAt(i)
		simpleLinkedList.RemoveAt(i)
	}

	for i := numberOfElements / 2; i < numberOfElements; i++ {
		linkedList.PushBack(i)
		simpleLinkedList.PushBack(i)
		provider.Add(i)
	}

	linkedList.Crystalize()
	benchy.New(b, options.Medium).
		RegisterBenchmark("linked list", provider.WrapBenchmarkFunc(func(i int) {
			linkedList.Get(i)
		})).
		RegisterBenchmark("simple linked list", provider.WrapBenchmarkFunc(func(i int) {
			simpleLinkedList.Get(i)
		})).
		Run()
}
