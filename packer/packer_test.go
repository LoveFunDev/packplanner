package packer

import (
	"example.com/packplanner/utils"
	"github.com/golang-collections/collections/stack"
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func areSameTotalWeight(a, b float64) bool {
	return utils.PrettyFormatFloat(a, 2) == utils.PrettyFormatFloat(b, 2)
}

func areSameItems(a, b *stack.Stack) bool {
	return reflect.DeepEqual(a, b)
}

func areSamePacks(a, b []*Pack) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i].ID != b[i].ID {
			return false
		}
		if a[i].MaxLength != b[i].MaxLength {
			return false
		}
		if !areSameTotalWeight(a[i].TotalWeight, b[i].TotalWeight) {
			return false
		}
		if a[i].ItemCount != b[i].ItemCount {
			return false
		}
		if !areSameItems(a[i].DistinctItems, b[i].DistinctItems) {
			return false
		}
	}

	return true
}

func TestNatural_Basic(t *testing.T) {
	we := NewGomegaWithT(t)

	// Input:
	// NATURAL,40,500.0
	// 1001,6200,30,9.653
	// 2001,7200,50,11.21
	sortOrder := Natural
	maxItems := 40
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 9.653},
		{ID: 2001, Length: 7200, Quantity: 50, Weight: 11.21},
	}

	// Output:
	// Pack Number: 1
	// 1001,6200,30,9.653
	// 2001,7200,10,11.21
	// Pack Length: 7200, Pack Weight: 401.69
	//
	// Pack Number: 2
	// 2001,7200,40,11.21
	// Pack Length: 7200, Pack Weight: 448.4
	p1 := &Pack{
		ID:            1,
		DistinctItems: stack.New(),
		MaxLength:     7200,
		TotalWeight:   401.69,
		ItemCount:     40,
	}
	p1.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 30, Weight: 9.653})
	p1.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 10, Weight: 11.21})

	p2 := &Pack{
		ID:            2,
		DistinctItems: stack.New(),
		MaxLength:     7200,
		TotalWeight:   448.4,
		ItemCount:     40,
	}
	p2.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 40, Weight: 11.21})
	expectedPacks := []*Pack{p1, p2}

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestNatural_ItemSliceEmpty(t *testing.T) {
	we := NewGomegaWithT(t)

	// NATURAL,40,500.0
	sortOrder := Natural
	maxItems := 40
	maxWeight := 500.0
	itemSlice := make([]*Item, 0)

	// Output:
	expectedPacks := make([]*Pack, 0)

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestNatural_AnyItemsWeightExceedMaxWeight(t *testing.T) {
	we := NewGomegaWithT(t)

	// NATURAL,40,500.0
	// 1001,6200,30,501.3
	sortOrder := Natural
	maxItems := 40
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 501.3},
	}

	// Output:
	expectedPacks := make([]*Pack, 0)

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestNatural_MaxItemsPerPackIsZero(t *testing.T) {
	we := NewGomegaWithT(t)

	// NATURAL,0,500.0
	// 1001,6200,30,501.3
	sortOrder := Natural
	maxItems := 0
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 501.3},
	}

	// Output:
	expectedPacks := make([]*Pack, 0)

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestNatural_MaxWeightPerPackIsZero(t *testing.T) {
	we := NewGomegaWithT(t)

	// NATURAL,30,0.0
	// 1001,6200,30,501.3
	sortOrder := Natural
	maxItems := 30
	maxWeight := 0.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 501.3},
	}

	// Output:
	expectedPacks := make([]*Pack, 0)

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestLongest_Basic(t *testing.T) {
	we := NewGomegaWithT(t)

	// LONG_TO_SHORT,40,500.0
	// 1001,6200,30,9.653
	// 2001,7200,50,11.21
	// 3001,4800,50,15.33
	// 4001,9800,10,4.364
	sortOrder := LongToShort
	maxItems := 40
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 9.653},
		{ID: 2001, Length: 7200, Quantity: 50, Weight: 11.21},
		{ID: 3001, Length: 4800, Quantity: 50, Weight: 15.33},
		{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364},
	}

	// Output:
	// Pack Number: 1
	// 4001,9800,10,4.364
	// 2001,7200,30,11.21
	// Pack Length: 9800, Pack Weight: 379.94
	//
	// Pack Number: 2
	// 2001,7200,20,11.21
	// 1001,6200,20,9.653
	// Pack Length: 7200, Pack Weight: 417.26
	//
	// Pack Number: 3
	// 1001,6200,10,9.653
	// 3001,4800,26,15.33
	// Pack Length: 6200, Pack Weight: 495.11
	//
	// Pack Number: 4
	// 3001,4800,24,15.33
	// Pack Length: 4800, Pack Weight: 367.92
	p1 := &Pack{
		ID:            1,
		DistinctItems: stack.New(),
		MaxLength:     9800,
		TotalWeight:   379.94,
		ItemCount:     40,
	}
	p1.DistinctItems.Push(&Item{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364})
	p1.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 30, Weight: 11.21})

	p2 := &Pack{
		ID:            2,
		DistinctItems: stack.New(),
		MaxLength:     7200,
		TotalWeight:   417.26,
		ItemCount:     40,
	}
	p2.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 20, Weight: 11.21})
	p2.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 20, Weight: 9.653})

	p3 := &Pack{
		ID:            3,
		DistinctItems: stack.New(),
		MaxLength:     6200,
		TotalWeight:   495.11,
		ItemCount:     36,
	}
	p3.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 10, Weight: 9.653})
	p3.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 26, Weight: 15.33})

	p4 := &Pack{
		ID:            4,
		DistinctItems: stack.New(),
		MaxLength:     4800,
		TotalWeight:   367.92,
		ItemCount:     24,
	}
	p4.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 24, Weight: 15.33})
	expectedPacks := []*Pack{p1, p2, p3, p4}

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestLongest_ItemsLengthAreSame(t *testing.T) {
	we := NewGomegaWithT(t)

	// LONG_TO_SHORT,40,500.0
	// 1001,6200,30,9.653
	// 2001,9800,5,11.21
	// 3001,4800,50,15.33
	// 4001,9800,10,4.364
	sortOrder := LongToShort
	maxItems := 40
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 9.653},
		{ID: 2001, Length: 9800, Quantity: 5, Weight: 11.21},
		{ID: 3001, Length: 4800, Quantity: 50, Weight: 15.33},
		{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364},
	}

	// Output:
	// Pack Number: 1
	// 2001,9800,5,11.21
	// 4001,9800,10,4.364
	// 1001,6200,25,9.653
	// Pack Length: 9800, Pack Weight: 341.02
	//
	// Pack Number: 2
	// 1001,6200,5,9.653
	// 3001,4800,29,15.33
	// Pack Length: 6200, Pack Weight: 492.83
	//
	// Pack Number: 3
	// 3001,4800,21,15.33
	// Pack Length: 4800, Pack Weight: 321.93
	p1 := &Pack{
		ID:            1,
		DistinctItems: stack.New(),
		MaxLength:     9800,
		TotalWeight:   341.02,
		ItemCount:     40,
	}
	p1.DistinctItems.Push(&Item{ID: 2001, Length: 9800, Quantity: 5, Weight: 11.21})
	p1.DistinctItems.Push(&Item{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364})
	p1.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 25, Weight: 9.653})

	p2 := &Pack{
		ID:            2,
		DistinctItems: stack.New(),
		MaxLength:     6200,
		TotalWeight:   492.83,
		ItemCount:     34,
	}
	p2.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 5, Weight: 9.653})
	p2.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 29, Weight: 15.33})

	p3 := &Pack{
		ID:            3,
		DistinctItems: stack.New(),
		MaxLength:     4800,
		TotalWeight:   321.93,
		ItemCount:     21,
	}
	p3.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 21, Weight: 15.33})
	expectedPacks := []*Pack{p1, p2, p3}

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	//printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}

func TestShortest_Basic(t *testing.T) {
	we := NewGomegaWithT(t)

	// SHORT_TO_LONG,40,500.0
	// 1001,6200,30,9.653
	// 2001,7200,50,11.21
	// 3001,4800,50,15.33
	// 4001,9800,10,4.364
	sortOrder := ShortToLong
	maxItems := 40
	maxWeight := 500.0
	itemSlice := []*Item{
		{ID: 1001, Length: 6200, Quantity: 30, Weight: 9.653},
		{ID: 2001, Length: 7200, Quantity: 50, Weight: 11.21},
		{ID: 3001, Length: 4800, Quantity: 50, Weight: 15.33},
		{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364},
	}

	// Output:
	// Pack Number: 1
	// 3001,4800,32,15.33
	// Pack Length: 4800, Pack Weight: 490.56
	//
	// Pack Number: 2
	// 3001,4800,18,15.33
	// 1001,6200,22,9.653
	// Pack Length: 6200, Pack Weight: 488.31
	//
	// Pack Number: 3
	// 1001,6200,8,9.653
	// 2001,7200,32,11.21
	// Pack Length: 7200, Pack Weight: 435.94
	//
	// Pack Number: 4
	// 2001,7200,18,11.21
	// 4001,9800,10,4.364
	// Pack Length: 9800, Pack Weight: 245.42
	p1 := &Pack{
		ID:            1,
		DistinctItems: stack.New(),
		MaxLength:     4800,
		TotalWeight:   490.56,
		ItemCount:     32,
	}
	p1.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 32, Weight: 15.33})

	p2 := &Pack{
		ID:            2,
		DistinctItems: stack.New(),
		MaxLength:     6200,
		TotalWeight:   488.31,
		ItemCount:     40,
	}
	p2.DistinctItems.Push(&Item{ID: 3001, Length: 4800, Quantity: 18, Weight: 15.33})
	p2.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 22, Weight: 9.653})

	p3 := &Pack{
		ID:            3,
		DistinctItems: stack.New(),
		MaxLength:     7200,
		TotalWeight:   435.94,
		ItemCount:     40,
	}
	p3.DistinctItems.Push(&Item{ID: 1001, Length: 6200, Quantity: 8, Weight: 9.653})
	p3.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 32, Weight: 11.21})

	p4 := &Pack{
		ID:            4,
		DistinctItems: stack.New(),
		MaxLength:     9800,
		TotalWeight:   245.42,
		ItemCount:     28,
	}
	p4.DistinctItems.Push(&Item{ID: 2001, Length: 7200, Quantity: 18, Weight: 11.21})
	p4.DistinctItems.Push(&Item{ID: 4001, Length: 9800, Quantity: 10, Weight: 4.364})
	expectedPacks := []*Pack{p1, p2, p3, p4}

	packer := NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
	packs := packer.packItems()
	printPacks(packs)
	isSame := areSamePacks(packs, expectedPacks)
	we.Expect(isSame).To(BeTrue())
}
