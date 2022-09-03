package pudzianmechaniczny

// Pm - Pudzian mechaniczny (Pudzianpower), W - Watt, HP - horsepower
const (
	gravity       = 9.81
	hp            = 745.7 // 1HP = 745.7W
	pudzian       = 1105
	pudzianHeight = 1.85
)

func PMToW(n float64) float64 {
	return n * pudzian * gravity * pudzianHeight // /1s is implicit
}

func PMToHP(n float64) float64 {
	return PMToW(n) / hp
}

func WToPM(n float64) float64 {
	return n / PMToW(1)
}

func HPToPM(n float64) float64 {
	return n / PMToHP(1)
}
