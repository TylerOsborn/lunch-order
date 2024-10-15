package utils

import (
	"encoding/csv"
	"fmt"
	"strings"
)

var ErrIncorrectCSVFormat = fmt.Errorf("incorrect CSV format")

func ParseCSV(input string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(input))

	r.FieldsPerRecord = 2

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %v", err)
	}

	return records, nil
}
