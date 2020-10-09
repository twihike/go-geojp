// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package geo

import (
	"math"
	"strings"
)

const (
	// MinLatitude is minimum latitude. Avoiding singularity at the poles.
	MinLatitude = -85.05112878
	// MaxLatitude is maximum latitude. Avoiding singularity at the poles.
	MaxLatitude = 85.05112878
	// MinLongitude is minimum longitude.
	MinLongitude = -180
	// MaxLongitude is maximum longitude.
	MaxLongitude = 180
)

// LatLongToQuadkey converts latitude and longitude to a quadkey.
func LatLongToQuadkey(lat, long float64, zoom int) string {
	pixelX, pixelY := LatLongToPixel(lat, long, zoom)
	tileX, tileY := PixelToTile(pixelX, pixelY)
	return TileToQuadkey(tileX, tileY, zoom)
}

// QuadkeyToLatLong converts a quadkey to latitude and longitude.
func QuadkeyToLatLong(quadkey string) (lat, long float64) {
	tileX, tileY, zoom := QuadkeyToTile(quadkey)
	pixelX, pixelY := TileToPixel(tileX, tileY)
	lat, long = PixelToLatLong(pixelX, pixelY, zoom)
	return lat, long
}

// Neighbors returns quadkeys that are adjacent to the specified quadkey.
func Neighbors(quadkey string, tile int) []string {
	tileX, tileY, zoom := QuadkeyToTile(quadkey)
	var ds [][]int64
	for dy := -tile; dy <= tile; dy++ {
		for dx := -tile; dx <= tile; dx++ {
			ds = append(ds, []int64{int64(dx), int64(dy)})
		}
	}
	min := int64(0)
	max := int64(math.Pow(4, float64(zoom)) - 1)

	neighbors := make([]string, 0, len(ds))
	for _, d := range ds {
		x := int64(tileX) + d[0]
		y := int64(tileY) + d[1]
		if x < min || x > max || y < min || y > max {
			continue
		}
		q := TileToQuadkey(uint(x), uint(y), zoom)
		neighbors = append(neighbors, q)
	}
	return neighbors
}

func clip(n, minValue, maxValue float64) float64 {
	return math.Min(math.Max(n, minValue), maxValue)
}

// MapSize returns the width and height of the map at the specified zoom level
// in pixels.
func MapSize(zoom uint) uint {
	return 256 << zoom
}

// GroundResolution returns the ground resolution at the specified latitude and
// zoom level.
func GroundResolution(lat float64, zoom int) float64 {
	const earthRadius = 6378137
	lat = clip(lat, MinLatitude, MaxLatitude)
	return math.Cos(lat*math.Pi/180) * 2 * math.Pi * earthRadius / float64(MapSize(uint(zoom)))
}

// LatLongToPixel converts latitude and longitude to pixel coordinates.
func LatLongToPixel(lat, long float64, zoom int) (pixelX, pixelY uint) {
	la := clip(lat, MinLatitude, MaxLatitude)
	lo := clip(long, MinLongitude, MaxLongitude)

	x := (lo + 180) / 360
	sinLatitude := math.Sin(la * math.Pi / 180)
	y := 0.5 - math.Log((1+sinLatitude)/(1-sinLatitude))/(4*math.Pi)

	mapSize := float64(MapSize(uint(zoom)))
	pixelX = uint(clip(x*mapSize+0.5, 0, mapSize-1))
	pixelY = uint(clip(y*mapSize+0.5, 0, mapSize-1))
	return pixelX, pixelY
}

// PixelToLatLong converts pixel coordinates to latitude and longitude.
func PixelToLatLong(pixelX, pixelY uint, zoom int) (lat, long float64) {
	mapSize := float64(MapSize(uint(zoom)))
	x := clip(float64(pixelX), 0, mapSize-1)/mapSize - 0.5
	y := 0.5 - clip(float64(pixelY), 0, mapSize-1)/mapSize

	lat = 90 - 360*math.Atan(math.Exp(-y*2*math.Pi))/math.Pi
	long = 360 * x
	return lat, long
}

// PixelToTile converts pixel coordinates to tile coordinates.
// One tile is 256x256 pixels.
func PixelToTile(pixelX, pixelY uint) (tileX, tileY uint) {
	tileX = pixelX / 256
	tileY = pixelY / 256
	return tileX, tileY
}

// TileToPixel converts tile coordinates to pixel coordinates.
// One tile is 256x256 pixels.
func TileToPixel(tileX, tileY uint) (pixelX, pixelY uint) {
	pixelX = tileX * 256
	pixelY = tileY * 256
	return pixelX, pixelY
}

// TileToQuadkey converts tile coordinates to a quadkey.
func TileToQuadkey(tileX, tileY uint, zoom int) string {
	var quadkey strings.Builder
	quadkey.Grow(zoom)
	for i := zoom; i > 0; i-- {
		digit := '0'
		mask := uint(1) << uint(i-1)
		if (tileX & mask) != 0 {
			digit++
		}
		if (tileY & mask) != 0 {
			digit += 2
		}
		quadkey.WriteRune(digit)
	}
	return quadkey.String()
}

// QuadkeyToTile converts a quadkey to tile coordinates.
func QuadkeyToTile(quadkey string) (tileX, tileY uint, zoom int) {
	tileX = 0
	tileY = 0
	zoom = len(quadkey)
	for i := zoom; i > 0; i-- {
		mask := uint(1) << uint(i-1)
		switch quadkey[zoom-i] {
		case '0':
		case '1':
			tileX |= mask
		case '2':
			tileY |= mask
		case '3':
			tileX |= mask
			tileY |= mask
		default:
			return 0, 0, 0
		}
	}
	return tileX, tileY, zoom
}
