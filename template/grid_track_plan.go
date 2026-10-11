package template

// gridTrackPlan is one track mid-resolution: the size it is heading for, the
// floor and ceiling it may not leave, and the fraction and the stretch that
// still owe it space.
type gridTrackPlan struct {
	size    int
	min     int
	max     int
	flex    float64
	stretch bool
}
