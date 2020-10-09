// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package geo

import (
	"reflect"
	"testing"
)

func TestLatLongToQuadkey(t *testing.T) {
	tests := []struct {
		name   string
		inLat  float64
		inLong float64
		inZ    int
		want   string
	}{
		{
			"normal",
			35.658584,
			139.7454316,
			23,
			"13300211230311333132022",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := LatLongToQuadkey(tt.inLat, tt.inLong, tt.inZ)
			if got != tt.want {
				t.Errorf("want = %v, got = %v", tt.want, got)
			}
		})
	}
}

func TestQuadkeyToLatLong(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantLat  float64
		wantLong float64
	}{
		{
			"normal",
			"13300211230311333132022",
			35.658586409152726,
			139.7454071044922,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotLat, gotLong := QuadkeyToLatLong(tt.in)
			if gotLat != tt.wantLat {
				t.Errorf("want = %v, got = %v", tt.wantLat, gotLat)
			}
			if gotLong != tt.wantLong {
				t.Errorf("want = %v, got = %v", tt.wantLong, gotLong)
			}
		})
	}
}

func TestNeighbors(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			"normal",
			"13300211230311333132022",
			[]string{
				"13300211230311333123131",
				"13300211230311333132020",
				"13300211230311333132021",
				"13300211230311333123133",
				"13300211230311333132022",
				"13300211230311333132023",
				"13300211230311333123311",
				"13300211230311333132200",
				"13300211230311333132201",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Neighbors(tt.in, 1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want = %v, got = %v", tt.want, got)
			}
		})
	}
}
