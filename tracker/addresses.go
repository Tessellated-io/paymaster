package tracker

import (
	"encoding/csv"
	"os"
)

type AddressTracker struct {
	addresses map[string]bool
	file      string
}

func NewAddressTracker(file string) *AddressTracker {
	a := &AddressTracker{
		addresses: make(map[string]bool),
		file:      file,
	}
	a.loadFromFile()

	return a
}

func (a *AddressTracker) AddAddress(address string) bool {
	if _, exists := a.addresses[address]; exists {
		return false // String already exists
	}
	a.addresses[address] = true
	a.saveToFile()

	return true
}

func (a *AddressTracker) saveToFile() error {
	file, err := os.Create(a.file)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for str := range a.addresses {
		writer.Write([]string{str})
	}
	return nil
}

func (a *AddressTracker) loadFromFile() error {
	file, err := os.Open(a.file)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) > 0 {
			a.addresses[record[0]] = true
		}
	}
	return nil
}
