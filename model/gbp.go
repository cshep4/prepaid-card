package model

type GBP int64

func ToGBP(f float64) GBP {
	return GBP((f * 100) + 0.5)
}

func (g GBP) ToPounds() float64 {
	x := float64(g)
	x = x / 100
	return x
}
