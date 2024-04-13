package main

import (
	"bufio"
	"errors"
	"example.com/packplanner/packer"
	"os"
	"strconv"
	"strings"
)

func extractSortOrder(sortOrderStr string) (packer.PackSortOrder, error) {
	trimmedStr := strings.TrimSpace(sortOrderStr)
	switch trimmedStr {
	case "NATURAL":
		return packer.Natural, nil
	case "SHORT_TO_LONG":
		return packer.ShortToLong, nil
	case "LONG_TO_SHORT":
		return packer.LongToShort, nil
	default:
		msg := "invalid sort order:" + trimmedStr + ". Supported values: NATURAL, SHORT_TO_LONG, LONG_TO_SHORT"
		return packer.InvalidOrder, errors.New(msg)
	}
}

func parseStartLine(lineFields []string) (packer.PackSortOrder, int, float64, error) {
	sortOrder := packer.InvalidOrder
	maxItems := 0
	maxWeight := 0.0

	sortOrder, err := extractSortOrder(lineFields[0])
	if err != nil {
		return sortOrder, maxItems, maxWeight, err
	}

	maxItems, err = strconv.Atoi(strings.TrimSpace(lineFields[1]))
	if err != nil {
		return sortOrder, maxItems, maxWeight, err
	}

	maxWeight, err = strconv.ParseFloat(strings.TrimSpace(lineFields[2]), 64)
	return sortOrder, maxItems, maxWeight, err
}

func parseItem(lineFields []string) (*packer.Item, error) {
	itemID, err := strconv.Atoi(strings.TrimSpace(lineFields[0]))
	if err != nil {
		return &packer.Item{}, err
	}

	itemLength, err := strconv.Atoi(strings.TrimSpace(lineFields[1]))
	if err != nil {
		return &packer.Item{}, err
	}

	itemQuantity, err := strconv.Atoi(strings.TrimSpace(lineFields[2]))
	if err != nil {
		return &packer.Item{}, err
	}

	itemWeight, err := strconv.ParseFloat(strings.TrimSpace(lineFields[3]), 64)
	if err != nil {
		return &packer.Item{}, err
	}

	return &packer.Item{
		ID:       itemID,
		Length:   itemLength,
		Quantity: itemQuantity,
		Weight:   itemWeight,
	}, nil
}

func main() {
	itemSlice := make([]*packer.Item, 0)
	sortOrder := packer.InvalidOrder
	maxItems := 0
	maxWeight := 0.0
	startLineParsed := false
	var err error

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		numOfFields := len(fields)
		switch numOfFields {
		case 3: // case start
			itemSlice = make([]*packer.Item, 0)
			sortOrder, maxItems, maxWeight, err = parseStartLine(fields)
			if err != nil { // parse failed, skip the testcase
				//fmt.Println("Start line parse failed. err=", err)
				continue
			}
			startLineParsed = true
		case 4: // item line
			if startLineParsed {
				item, err := parseItem(fields)
				if err != nil { // parse failed, skip the testcase
					//fmt.Println("Item line parse failed. err=", err)
					startLineParsed = false
					continue
				}
				itemSlice = append(itemSlice, item)
			}
		case 1: // empty line is case end
			if strings.TrimSpace(fields[0]) == "" && startLineParsed {
				p := packer.NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
				p.Process()
				startLineParsed = false
			}
		}
	}

	if scanner.Err() == nil && startLineParsed { // EOF received, run the last one
		p := packer.NewPacker(sortOrder, maxItems, maxWeight, itemSlice)
		p.Process()
	}
}
