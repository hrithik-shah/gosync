package storage

import "errors"

type Garage struct {
	// Add fields for the Garage struct here
}

func NewGarage() *Garage {
	return &Garage{
		// Initialize fields for the Garage struct here
	}
}

func (g *Garage) StoreFile(data []byte) (id string, err error) {
	// Implement the logic to store data in the garage
	return "Not implemented", errors.New("Not implemented")
}

func (g *Garage) RetrieveFile(id string) (data []byte, err error) {
	// Implement the logic to retrieve data from the garage
	return nil, errors.New("Not implemented")
}

func (g *Garage) DeleteFile(id string) (err error) {
	// Implement the logic to delete data from the garage
	return errors.New("Not implemented")
}
