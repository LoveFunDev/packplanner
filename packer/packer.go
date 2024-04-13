package packer

import (
	"example.com/packplanner/utils"
	"fmt"
	"sort"
)

// Item represents an individual item.
type Item struct {
	ID       int
	Length   int
	Quantity int
	Weight   float64
}

// Pack represents a group of items.
type Pack struct {
	ID            int
	DistinctItems []*Item
	MaxLength     int
	TotalWeight   float64
	ItemCount     int
}

// PackSortOrder defines the sorting criteria for packs.
type PackSortOrder int

const (
	Natural PackSortOrder = iota
	ShortToLong
	LongToShort
	InvalidOrder
)

// Packer manages pack resources and operations
type Packer struct {
	SortOrder PackSortOrder
	MaxItems  int
	MaxWeight float64
	ItemSlice []*Item
}

// NewPacker creates a packer to handle pack operation
func NewPacker(sortOrder PackSortOrder, maxItems int, maxWeight float64, itemSlice []*Item) *Packer {
	return &Packer{
		SortOrder: sortOrder,
		MaxItems:  maxItems,
		MaxWeight: maxWeight,
		ItemSlice: itemSlice,
	}
}

// Process sorts items into packs, then output to stdout
func (packer *Packer) Process() {
	packs := packer.packItems()
	printPacks(packs)
}

func (packer *Packer) packItems() []*Pack {
	packer.sortItemSliceByOrder()

	packs := make([]*Pack, 0)
	currentPack := &Pack{ID: 1, MaxLength: 0, TotalWeight: 0.0, ItemCount: 0}

	for _, item := range packer.ItemSlice {
		// If an item's weight exceeds the pack's maximum weight, cannot add it to any pack.
		if item.Weight > packer.MaxWeight {
			continue
		}

		for i := 0; i < item.Quantity; i++ {
			if packer.checkIfExceedLimits(currentPack, item) {
				packs = append(packs, currentPack)
				currentPack = &Pack{ID: len(packs) + 1, MaxLength: 0, TotalWeight: 0.0, ItemCount: 0}
			}

			addItemToPack(currentPack, item)
		}
	}

	if currentPack.ItemCount > 0 { // Add the last pack
		packs = append(packs, currentPack)
	}

	return packs
}

func (packer *Packer) sortItemSliceByOrder() {
	switch packer.SortOrder {
	case ShortToLong:
		sort.Slice(packer.ItemSlice, func(i, j int) bool {
			return packer.ItemSlice[i].Length < packer.ItemSlice[j].Length
		})
	case LongToShort:
		sort.Slice(packer.ItemSlice, func(i, j int) bool {
			return packer.ItemSlice[i].Length > packer.ItemSlice[j].Length
		})
	}
}

func (packer *Packer) checkIfExceedLimits(currentPack *Pack, item *Item) bool {
	return (currentPack.ItemCount+1) > packer.MaxItems || (currentPack.TotalWeight+item.Weight) > packer.MaxWeight
}

func addItemToPack(currentPack *Pack, item *Item) {
	currentPack.TotalWeight += item.Weight
	currentPack.ItemCount++
	if item.Length > currentPack.MaxLength {
		currentPack.MaxLength = item.Length
	}

	length := len(currentPack.DistinctItems)
	if length > 0 && currentPack.DistinctItems[length-1].ID == item.ID {
		currentPack.DistinctItems[length-1].Quantity++
	} else {
		newItem := &Item{ID: item.ID, Length: item.Length, Quantity: 1, Weight: item.Weight}
		currentPack.DistinctItems = append(currentPack.DistinctItems, newItem)
	}
}

func printPacks(packs []*Pack) {
	for _, pack := range packs {
		fmt.Printf("Pack Number: %d\n", pack.ID)
		for _, item := range pack.DistinctItems {
			fmt.Printf("%d,%d,%d,%s\n", item.ID, item.Length, item.Quantity, utils.PrettyFormatFloat(item.Weight, -1))
		}
		fmt.Printf("Pack Length: %d, Pack Weight: %s\n\n", pack.MaxLength, utils.PrettyFormatFloat(pack.TotalWeight, 2))
	}
}
