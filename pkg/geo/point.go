// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package geo handles geographic.
package geo

import (
	"math"
)

// Point represents a position in space.
type Point struct {
	X float64
	Y float64
}

// Distance returns the distance to the specified point.
func (p *Point) Distance(target Point) float64 {
	xs := math.Pow(target.X-p.X, 2)
	ys := math.Pow(target.Y-p.Y, 2)
	d := math.Sqrt(xs + ys)
	return d
}

// Nearest returns the closest point among the specified points.
func (p *Point) Nearest(targets []Point) (int, float64) {
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
