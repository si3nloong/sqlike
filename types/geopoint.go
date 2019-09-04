package types

// GeoPoint :
type GeoPoint [2]float64

// SetLatitude :
func (gp *GeoPoint) SetLatitude(lat float64) {
	gp[0] = lat
	return
}

// SetLongitude :
func (gp *GeoPoint) SetLongitude(lng float64) {
	gp[1] = lng
	return
}

// Latitude :
func (gp GeoPoint) Latitude() float64 {
	return gp[0]
}

// Longitude :
func (gp GeoPoint) Longitude() float64 {
	return gp[1]
}
