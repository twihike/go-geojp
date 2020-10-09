// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package geo

import (
	"math"
	"reflect"
	"testing"
)

func TestPoint_Distance(t *testing.T) {
	p := Point{X: 1, Y: 1}
	tests := []struct {
		name string
		in   Point
		want float64
	}{
		{"normal", Point{X: 2, Y: 2}, math.Sqrt(2)},
		{"normal", Point{X: 4, Y: 4}, math.Sqrt(9 + 9)},
		{"normal", Point{X: 8, Y: 4}, math.Sqrt(49 + 9)},
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

func TestPoint_Nearest(t *testing.T) {
	p := Point{X: 2, Y: 2}
	tests := []struct {
		name     string
		in       []Point
		wantIdx  int
		wantDist float64
	}{
		{"normal", []Point{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 3, Y: 3}}, 1, 1},
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
