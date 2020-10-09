// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package jp

import (
	"testing"

	"github.com/twihike/go-geojp/pkg/geo"
)

func TestAPs_Nearest(t *testing.T) {
	aps, err := ReadAPsFromFile("../../../testdata/japanese-addresses.csv")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		in           geo.LatLong
		wantAreaCode string
		wantDistance float64
	}{
		{
			"normal",
			geo.LatLong{Latitude: 35.658584, Longitude: 139.7454316},
			"131030002003",
			220.37123693585445,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := aps.Nearest(tt.in)
			if got.AreaCode != tt.wantAreaCode {
				t.Errorf("want = %v, got = %v", tt.wantAreaCode, got.AreaCode)
			}
			if got.Distance != tt.wantDistance {
				t.Errorf("want = %v, got = %v", tt.wantDistance, got.Distance)
			}
		})
	}
}

func TestAPs_Near(t *testing.T) {
	aps, err := ReadAPsFromFile("../../../testdata/japanese-addresses.csv")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		in           geo.LatLong
		wantAreaCode string
		wantDistance float64
	}{
		{
			"normal",
			geo.LatLong{Latitude: 35.658584, Longitude: 139.7454316},
			"131030002003",
			220.37123693585445,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := aps.Near(tt.in, 18)
			if got[0].AreaCode != tt.wantAreaCode {
				t.Errorf("want = %v, got = %v", tt.wantAreaCode, got[0].AreaCode)
			}
			if got[0].Distance != tt.wantDistance {
				t.Errorf("want = %v, got = %v", tt.wantDistance, got[0].Distance)
			}
		})
	}
}

func TestIndexedAPs_Nearest(t *testing.T) {
	aps, err := ReadAPsFromFile("../../../testdata/japanese-addresses.csv")
	iaps := CreateIndexedAPs(aps)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		in           geo.LatLong
		wantAreaCode string
		wantDistance float64
	}{
		{
			"normal",
			geo.LatLong{Latitude: 35.658584, Longitude: 139.7454316},
			"131030002003",
			220.37123693585445,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := iaps.Nearest(tt.in)
			if got.AreaCode != tt.wantAreaCode {
				t.Errorf("want = %v, got = %v", tt.wantAreaCode, got.AreaCode)
			}
			if got.Distance != tt.wantDistance {
				t.Errorf("want = %v, got = %v", tt.wantDistance, got.Distance)
			}
		})
	}
}

func TestIndexedAPs_Near(t *testing.T) {
	aps, err := ReadAPsFromFile("../../../testdata/japanese-addresses.csv")
	iaps := CreateIndexedAPs(aps)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		in           geo.LatLong
		wantAreaCode string
		wantDistance float64
	}{
		{
			"normal",
			geo.LatLong{Latitude: 35.658584, Longitude: 139.7454316},
			"131030002003",
			220.37123693585445,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := iaps.Near(tt.in, 18)
			if got[0].AreaCode != tt.wantAreaCode {
				t.Errorf("want = %v, got = %v", tt.wantAreaCode, got[0].AreaCode)
			}
			if got[0].Distance != tt.wantDistance {
				t.Errorf("want = %v, got = %v", tt.wantDistance, got[0].Distance)
			}
		})
	}
}
