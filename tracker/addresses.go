package tracker

import (
	"encoding/csv"
	"fmt"
	"os"
)

type AddressTracker struct {
	addresses map[string]bool
	file      string
}

func NewAddressTracker(file string) (*AddressTracker, error) {
	a := &AddressTracker{
		addresses: make(map[string]bool),
		file:      file,
	}
	err := a.loadFromFile()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *AddressTracker) AddAddress(address string) error {
	if _, exists := a.addresses[address]; exists {
		return fmt.Errorf("address already exists: %s", address)
	}
	a.addresses[address] = true
	err := a.saveToFile()
	if err != nil {
		return err
	}

	return nil
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
		err := writer.Write([]string{str})
		if err != nil {
			return err
		}
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
