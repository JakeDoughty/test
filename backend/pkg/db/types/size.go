package types

type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func (size *Size) IsZero() bool {
	return size.Width == 0 && size.Height == 0
}
