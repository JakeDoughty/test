package types

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

func (location *Location) IsZero() bool {
	return location.Latitude == 0 && location.Longitude == 0
}
