# go-geojp

[![ci status](https://github.com/twihike/go-geojp/workflows/ci/badge.svg)](https://github.com/twihike/go-geojp/actions) [![license](https://img.shields.io/github/license/twihike/go-geojp)](LICENSE)

An API server for Japanese geographics.

## Installation

```shell
go get -u github.com/twihike/go-geojp/cmd/...
```

## Usage

Start the server.

```shell
# Download Japanese addresses.
curl -LO https://raw.githubusercontent.com/geolonia/japanese-addresses/master/data/latest.csv
# Set file path.
export ADDR_POS_PATH=latest.csv
# Run server
./geojp
```

Reverse geocoding.

```shell
curl -sS \
  -X POST localhost:8080/api/reverse-geocoding \
  -d 'latitude=35.658584' \
  -d 'longitude=139.7454316' \
| jq .
```

Output:

```json
{
  "pref_name": "東京都",
  "city_name": "港区",
  "area_name": "芝公園三丁目",
  "latitude": 35.659943,
  "longitude": 139.747207,
  "distance": 220.37123693585445
}
```

Geocoding.

```shell
curl -sS \
  -X POST localhost:8080/api/geocoding \
  -d 'area_name=芝公園' \
| jq .
```

Output:

```js
[
  {
    "pref_name": "広島県",
    "city_name": "広島市西区",
    "area_name": "大芝公園",
    "latitude": 34.417138,
    "longitude": 132.460336
  },
  {
    "pref_name": "東京都",
    "city_name": "港区",
    "area_name": "芝公園三丁目",
    "latitude": 35.659943,
    "longitude": 139.747207
  },
  // ...
  {
    "pref_name": "東京都",
    "city_name": "港区",
    "area_name": "芝公園二丁目",
    "latitude": 35.655131,
    "longitude": 139.751235
  }
]
```

With current location.

```shell
curl -sS \
  -X POST localhost:8080/api/geocoding \
  -d 'area_name=芝公園' \
  -d 'latitude=35.658584' \
  -d 'longitude=139.7454316' \
| jq .
```

Output:

```js
[
  {
    "pref_name": "東京都",
    "city_name": "港区",
    "area_name": "芝公園三丁目",
    "latitude": 35.659943,
    "longitude": 139.747207,
    "distance": 220.37123693585445
  },
  {
    "pref_name": "東京都",
    "city_name": "港区",
    "area_name": "芝公園四丁目",
    "latitude": 35.656459,
    "longitude": 139.74764,
    "distance": 309.2609283463855
  },
  // ...
  {
    "pref_name": "広島県",
    "city_name": "広島市西区",
    "area_name": "大芝公園",
    "latitude": 34.417138,
    "longitude": 132.460336,
    "distance": 677296.9280660086
  }
]
```

## Credits

[japanese-addresses](https://geolonia.github.io/japanese-addresses/) by [geolonia](https://github.com/geolonia) is licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/).

## License

Copyright (c) 2020 twihike. All rights reserved.

This project is licensed under the terms of the MIT license.
