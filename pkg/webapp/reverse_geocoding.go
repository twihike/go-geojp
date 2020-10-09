// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package webapp

import (
	"encoding/json"
	"net/http"

	"github.com/twihike/go-geojp/pkg/geo"
	"github.com/twihike/go-structconv/structconv"
)

type reverseGeocodingInput struct {
	Latitude  float64 `strmap:"latitude,required"`
	Longitude float64 `strmap:"longitude,required"`
	Zoom      int     `strmap:"zoom"`
}

type reverseGeocodingOutput struct {
	PrefName  string  `json:"pref_name"`
	CityName  string  `json:"city_name"`
	AreaName  string  `json:"area_name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Distance  float64 `json:"distance"`
}

func reverseGeocoding(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	strMap := map[string]string{}
	for k, v := range r.Form {
		if len(v) > 0 {
			strMap[k] = v[0]
		}
	}

	var in reverseGeocodingInput
	if err := structconv.DecodeStringMap(strMap, &in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	target := geo.LatLong{Latitude: in.Latitude, Longitude: in.Longitude}
	var body interface{}
	if in.Zoom > 0 {
		filteredAPs := iaps.Near(target, in.Zoom)
		b := []reverseGeocodingOutput{}
		for _, ap := range filteredAPs {
			b = append(b, reverseGeocodingOutput{
				PrefName:  ap.PrefName,
				CityName:  ap.CityName,
				AreaName:  ap.AreaName,
				Latitude:  ap.Latitude,
				Longitude: ap.Longitude,
				Distance:  ap.Distance,
			})
		}
		body = b
	} else {
		filteredAPs := iaps.Nearest(target)
		b := reverseGeocodingOutput{
			PrefName:  filteredAPs.PrefName,
			CityName:  filteredAPs.CityName,
			AreaName:  filteredAPs.AreaName,
			Latitude:  filteredAPs.Latitude,
			Longitude: filteredAPs.Longitude,
			Distance:  filteredAPs.Distance,
		}
		body = b
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
