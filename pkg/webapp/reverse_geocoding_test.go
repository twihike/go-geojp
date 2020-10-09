// Copyright (c) 2020 twihike. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package webapp

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/twihike/go-geojp/pkg/geo/jp"
)

func TestReverseGeocoding(t *testing.T) {
	aps, err := jp.ReadAPsFromFile("../../testdata/japanese-addresses.csv")
	if err != nil {
		t.Fatal(err)
	}
	iaps = jp.CreateIndexedAPs(aps)

	target := "http://example.com/api/reverse-geocoding"
	body := url.Values{}
	body.Set("latitude", "35.658584")
	body.Set("longitude", "139.7454316")
	bodyReader := strings.NewReader(body.Encode())
	req := httptest.NewRequest(http.MethodPost, target, bodyReader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	got := httptest.NewRecorder()
	reverseGeocoding(got, req)

	want := `{"pref_name":"東京都","city_name":"港区","area_name":"芝公園三丁目","latitude":35.659943,"longitude":139.747207,"distance":220.37123693585445}
`
	if got.Code != http.StatusOK {
		t.Errorf("want = %v, got = %v", http.StatusOK, got.Code)
	}
	if got := got.Body.String(); got != want {
		t.Errorf("\nwant = %v\ngot  = %v", want, got)
	}
}
