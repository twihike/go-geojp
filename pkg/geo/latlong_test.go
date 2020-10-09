// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package geo

import (
	"reflect"
	"testing"
)

func TestLatLong_Distance(t *testing.T) {
	p := LatLong{Latitude: 1, Longitude: 1}
	tests := []struct {
		name string
		in   LatLong
		want float64
	}{
		{"normal", LatLong{Latitude: 2, Longitude: 2}, 157225.4320380729},
		{"normal", LatLong{Latitude: 4, Longitude: 4}, 471508.55302713305},
		{"normal", LatLong{Latitude: 8, Longitude: 4}, 846347.8879953761},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := p.Distance(tt.in); got != tt.want {
				t.Errorf("want = %v, got = %v", tt.want, got)
			}
		})
	}
}

func TestLatLong_Nearest(t *testing.T) {
	p := LatLong{Latitude: 2, Longitude: 2}
	tests := []struct {
		name     string
		in       []LatLong
		wantIdx  int
		wantDist float64
	}{
		{
			"normal",
			[]LatLong{
				{Latitude: 1, Longitude: 1},
				{Latitude: 1, Longitude: 2},
				{Latitude: 3, Longitude: 3},
			},
			1,
			111194.92664455874,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotIdx, gotDist := p.Nearest(tt.in)
			if !reflect.DeepEqual(gotIdx, tt.wantIdx) {
				t.Errorf("want = %v, got = %v", tt.wantIdx, gotIdx)
			}
			if !reflect.DeepEqual(gotDist, tt.wantDist) {
				t.Errorf("want = %v, got = %v", tt.wantDist, gotDist)
			}
		})
	}
}
