package circleToPolygon

// version 1.0.3

import (
	"encoding/json"
	"math"
)

const (
	defaultEarthRadius float64 = 6378137 // equatorial Earth radius
	defaultCount       float64 = 32
)

// create new circle
func NewCircle(latitude float64, longtitude float64, radius float64) CircleToPolygon {
	return &circleToPolygon{
		center: [2]float64{latitude, longtitude},
		radius: radius,
	}
}

type CircleToPolygon interface {
	SetEarthRadius(float64) CircleToPolygon // set earth radius
	SetBearing(float64) CircleToPolygon     // set bearing
	SetDirection(float64) CircleToPolygon   // set circle direction
	Draw() [][2]float64                     // draw circle by options
	DrawGeoJson() []byte                    // draw in geoJson
}

type circleToPolygon struct {
	earthRadius float64
	radius      float64
	center      [2]float64
	bearing     float64
	direction   float64
	count       float64
}

func (ctp *circleToPolygon) SetEarthRadius(earthRadius float64) CircleToPolygon {
	ctp.earthRadius = earthRadius
	return ctp
}

func (ctp *circleToPolygon) getEarthRadius() float64 {
	if ctp.earthRadius == 0 {
		return defaultEarthRadius
	}
	return ctp.earthRadius
}

func (ctp *circleToPolygon) getRadius() float64 {
	return ctp.radius
}

func (ctp *circleToPolygon) getCenter() [2]float64 {
	return ctp.center
}

func (ctp *circleToPolygon) SetBearing(bearing float64) CircleToPolygon {
	ctp.bearing = bearing
	return ctp
}

func (ctp *circleToPolygon) getBearing() float64 {
	return ctp.bearing
}

func (ctp *circleToPolygon) SetDirection(direction float64) CircleToPolygon {
	ctp.direction = direction
	return ctp
}

func (ctp *circleToPolygon) getDirection() float64 {
	if ctp.direction != -1 && ctp.direction != 1 {
		return 1
	}
	return ctp.direction
}

func (ctp *circleToPolygon) SetCount(count float64) CircleToPolygon {
	ctp.count = count
	return ctp
}

func (ctp *circleToPolygon) getCount() float64 {
	if ctp.count <= 3 {
		return defaultCount
	}
	return ctp.count
}

func (ctp *circleToPolygon) Draw() [][2]float64 {
	count := ctp.getCount()
	bearing := ctp.getBearing()
	direction := ctp.getDirection()
	earthRadius := ctp.getEarthRadius()
	radius := ctp.getRadius()
	center := ctp.getCenter()
	start := toRadians(bearing)
	coordinates := make([][2]float64, 0)
	for i := 0.0; i < count; i++ {
		coordinates = append(coordinates, offset(center, radius, earthRadius, start+(direction*2*math.Pi*-i)/count))
	}
	return coordinates
}

func (ctp *circleToPolygon) DrawGeoJson() []byte {
	data, _ := json.Marshal(ctp.Draw()) // always correct
	result := []byte{}
	result = append(result, []byte(`{"coordinates":[`)...)
	result = append(result, data...)
	result = append(result, []byte(`],"type":"Polygon"}`)...)
	return result
}

func toRadians(angleInDegrees float64) float64 {
	return (angleInDegrees * math.Pi) / 180
}

func toDegrees(angleInRadians float64) float64 {
	return (angleInRadians * 180) / math.Pi
}

func offset(c1 [2]float64, distance float64, earthRadius float64, bearing float64) [2]float64 {
	lat1 := toRadians(c1[0])
	lon1 := toRadians(c1[1])
	dByR := distance / earthRadius
	lat := math.Asin(math.Sin(lat1)*math.Cos(dByR) + math.Cos(lat1)*math.Sin(dByR)*math.Cos(bearing))
	lon := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(dByR)*math.Cos(lat1),
		math.Cos(dByR)-math.Sin(lat1)*math.Sin(lat),
	)
	return [2]float64{toDegrees(lon), toDegrees(lat)}
}
