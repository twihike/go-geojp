// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package geo

import (
	"math"
)

// LatLong represents a position on the earth.
type LatLong struct {
	Latitude  float64
	Longitude float64
}

// Distance returns the distance to the specified point.
func (p *LatLong) Distance(target LatLong) float64 {
	const r = 6371000
	rad := math.Pi / 180
	lat1 := p.Latitude * rad
	lat2 := target.Latitude * rad
	sinDLat := math.Sin((target.Latitude - p.Latitude) * rad / 2)
	sinDLong := math.Sin((target.Longitude - p.Longitude) * rad / 2)
	a := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLong*sinDLong
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return r * c
}

// Nearest returns the closest point among the specified points.
func (p *LatLong) Nearest(targets []LatLong) (int, float64) {
	var minIndex int
	minDistance := math.MaxFloat64
	for i, target := range targets {
		d := p.Distance(target)
		if d < minDistance {
			minDistance = d
			minIndex = i
		}
	}
	return minIndex, minDistance
}
