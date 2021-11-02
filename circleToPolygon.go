package circletopolygon

import (
	"math"
)

const (
	defaultEarthRadius float64 = 6378137 // equatorial Earth radius
	defaultCount       float64 = 32
)

type CircleToPolygon interface {
	SetEarthRadius(float64) CircleToPolygon
	SetRadius(float64) CircleToPolygon
	SetCenter([2]float64) CircleToPolygon
	SetBearing(float64) CircleToPolygon
	SetDirection(float64) CircleToPolygon
	Draw() [][2]float64
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

func (ctp *circleToPolygon) SetRadius(radius float64) CircleToPolygon {
	ctp.radius = radius
	return ctp
}

func (ctp *circleToPolygon) getRadius() float64 {
	return ctp.radius
}

func (ctp *circleToPolygon) SetCenter(center [2]float64) CircleToPolygon {
	ctp.center = center
	return ctp
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

func toRadians(angleInDegrees float64) float64 {
	return (angleInDegrees * math.Pi) / 180
}

func toDegrees(angleInRadians float64) float64 {
	return (angleInRadians * 180) / math.Pi
}

func offset(c1 [2]float64, distance float64, earthRadius float64, bearing float64) [2]float64 {
	lat1 := toRadians(c1[1])
	lon1 := toRadians(c1[0])
	dByR := distance / earthRadius
	lat := math.Asin(math.Sin(lat1)*math.Cos(dByR) + math.Cos(lat1)*math.Sin(dByR)*math.Cos(bearing))
	lon := lon1 + math.Atan2(
		math.Sin(bearing)*math.Sin(dByR)*math.Cos(lat1),
		math.Cos(dByR)-math.Sin(lat1)*math.Sin(lat),
	)
	return [2]float64{toDegrees(lon), toDegrees(lat)}
}
