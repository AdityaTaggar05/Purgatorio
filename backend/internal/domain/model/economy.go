package model

type Currency string

const (
	CurrencyPenitence Currency = "penitence"
	CurrencyGrace     Currency = "grace"
)

func (c Currency) String() string {
	return string(c)
}

type ResourceGenerationInfo struct {
	BuildingID      string
	BuildingLevel   int
	CurrentRate     int
	StorageCapacity int
	Currency        Currency
	Metadata        BuildingMetadata
}
