package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

// all measurements should be in millimeters

func main() {
	plate, err := plate()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	render.ToSTL(plate, 1000, "keeb.stl", &render.MarchingCubesUniform{})
	keycaps, err := drawKeyCaps()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	render.ToSTL(keycaps, 1000, "keycaps.stl", &render.MarchingCubesOctree{})
}

func rect(x1, z1, x2, y, z2 float64) (sdf.SDF3, error) {
	dimensions := []sdf.V2{
		{x1, z1},
		{x2, z1},
		{x2, z2},
		{x1, z2},
	}
	p := sdf.NewPolygon()
	p.AddV2Set(dimensions)
	profile, err := sdf.Polygon2D(p.Vertices())
	if err != nil {
		return nil, err
	}
	return sdf.Extrude3D(profile, y), err
}
