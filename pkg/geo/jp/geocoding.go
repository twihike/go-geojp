// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

// Package jp handles Japanese geographic.
package jp

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/twihike/go-geojp/pkg/geo"
)

// AddressPosition is a Japanese address and its position.
type AddressPosition struct {
	PrefCode     string
	PrefName     string
	PrefKanaName string
	PrefRomaName string
	CityCode     string
	CityName     string
	CityKanaName string
	CityRomaName string
	AreaCode     string
	AreaName     string
	Latitude     float64
	Longitude    float64
	quadkey      string
}

// AddressPositions is a slice of AddressPosition.
type AddressPositions []AddressPosition

// IndexedAPs is a map of AddressPosition keyed by the quadkey.
type IndexedAPs map[string]AddressPositions

// NearbyAP is a nearby AddressPosition.
type NearbyAP struct {
	AddressPosition
	Distance float64
}

// ReadAPsFromFile reads AddressPositions from a file.
func ReadAPsFromFile(path string) (AddressPositions, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	aps := make(AddressPositions, 0, len(records))
	for i, record := range records {
		// Skip header.
		if i == 0 {
			continue
		}

		lat, err := strconv.ParseFloat(record[10], 64)
		if err != nil {
			return nil, err
		}
		long, err := strconv.ParseFloat(record[11], 64)
		if err != nil {
			return nil, err
		}
		quadkey := geo.LatLongToQuadkey(lat, long, 23)

		ap := AddressPosition{
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
			record[7],
			record[8],
			record[9],
			lat,
			long,
			quadkey,
		}
		aps = append(aps, ap)
	}
	return aps, nil
}

// FindByAreaName returns address positions containing the specified name.
func (aps AddressPositions) FindByAreaName(n string, base geo.LatLong) []NearbyAP {
	var unordered AddressPositions
	for _, ap := range aps {
		if strings.Contains(ap.AreaName, n) {
			unordered = append(unordered, ap)
		}
	}
	result := make([]NearbyAP, len(unordered))
	for i, ap := range unordered {
		d := base.Distance(geo.LatLong{
			Latitude:  ap.Latitude,
			Longitude: ap.Longitude,
		})
		result[i] = NearbyAP{ap, d}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})
	return result
}

// Nearest returns an address position closest to the specified position.
func (aps AddressPositions) Nearest(p geo.LatLong) NearbyAP {
	points := make([]geo.LatLong, len(aps))
	for i, ap := range aps {
		po := geo.LatLong{Latitude: ap.Latitude, Longitude: ap.Longitude}
		points[i] = po
	}
	i, d := p.Nearest(points)
	return NearbyAP{aps[i], d}
}

// Near returns address positions close to the specified position.
func (aps AddressPositions) Near(p geo.LatLong, zoom int) []NearbyAP {
	quadkey := geo.LatLongToQuadkey(p.Latitude, p.Longitude, zoom)
	quadkeys := geo.Neighbors(quadkey, 1)
	var near []NearbyAP
	for _, ap := range aps {
		q := geo.LatLongToQuadkey(ap.Latitude, ap.Longitude, zoom)
		for _, qk := range quadkeys {
			if q == qk {
				t := geo.LatLong{Latitude: ap.Latitude, Longitude: ap.Longitude}
				a := NearbyAP{ap, p.Distance(t)}
				near = append(near, a)
			}
		}
	}
	sort.Slice(near, func(i, j int) bool {
		return near[i].Distance < near[j].Distance
	})
	return near
}

// CreateIndexedAPs creates IndexedAPs from the specified data.
func CreateIndexedAPs(aps AddressPositions) IndexedAPs {
	const (
		maxZoomLevel = 23
		minZoomLevel = 4
	)
	index := IndexedAPs{}
	for zoom := minZoomLevel; zoom <= maxZoomLevel; zoom++ {
		createIndexedAPsByZoomLevel(index, aps, zoom)
	}
	return index
}

func createIndexedAPsByZoomLevel(idx IndexedAPs, aps AddressPositions, zoom int) {
	for _, ap := range aps {
		key := ap.quadkey[0:zoom]
		if s, ok := idx[key]; !ok {
			idx[key] = AddressPositions{ap}
		} else {
			idx[key] = append(s, ap)
		}
	}
}

// Nearest returns an address position closest to the specified position.
func (idx IndexedAPs) Nearest(p geo.LatLong) NearbyAP {
	const (
		maxZoomLevel = 23
		minZoomLevel = 4
		minHits      = 10
	)

	for zoom := maxZoomLevel; zoom >= minZoomLevel; zoom-- {
		quadkey := geo.LatLongToQuadkey(p.Latitude, p.Longitude, zoom)
		quadkeys := geo.Neighbors(quadkey, 1)
		var aps []AddressPosition
		var points []geo.LatLong
		for _, q := range quadkeys {
			filteredAPs, ok := idx[q[0:zoom]]
			if !ok {
				continue
			}
			for _, ap := range filteredAPs {
				p := geo.LatLong{
					Latitude:  ap.Latitude,
					Longitude: ap.Longitude,
				}
				aps = append(aps, ap)
				points = append(points, p)
			}
		}

		if len(aps) > minHits {
			i, d := p.Nearest(points)
			nearest := NearbyAP{aps[i], d}
			return nearest
		}
	}

	var allPoints []geo.LatLong
	var allAPs AddressPositions
	for _, aps := range idx {
		for _, ap := range aps {
			po := geo.LatLong{Latitude: ap.Latitude, Longitude: ap.Longitude}
			allPoints = append(allPoints, po)
			allAPs = append(allAPs, ap)
		}
	}
	minIdx, minDist := p.Nearest(allPoints)
	return NearbyAP{allAPs[minIdx], minDist}
}

// Near returns address positions close to the specified position.
func (idx IndexedAPs) Near(p geo.LatLong, zoom int) []NearbyAP {
	quadkey := geo.LatLongToQuadkey(p.Latitude, p.Longitude, zoom)
	quadkeys := geo.Neighbors(quadkey, 1)
	var aps []NearbyAP
	for _, q := range quadkeys {
		filterdAPs, ok := idx[q[0:zoom]]
		if !ok {
			continue
		}
		for _, ap := range filterdAPs {
			t := geo.LatLong{
				Latitude:  ap.Latitude,
				Longitude: ap.Longitude,
			}
			d := p.Distance(t)
			newAP := NearbyAP{ap, d}
			aps = append(aps, newAP)
		}
	}
	sort.Slice(aps, func(i, j int) bool {
		return aps[i].Distance < aps[j].Distance
	})
	return aps
}
