package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

//-----------------------------------------------------------------------------

// in millimeters

type row []keycap

func newRow(inX, inZ, gap float64, keys ...float64) row {
	row := make([]keycap, 0, 0)
	x, z := inX, inZ
	for _, ele := range keys {
		row = append(row, keycap{x, x + ele, z, z + CAPWIDTH})
		x = x + gap + ele
	}
	/*
		for _, ele := range keys {
			dimensions = append(dimensions, sdf.V2{x, z})
			dimensions = append(dimensions, sdf.V2{x + ele, z})
			dimensions = append(dimensions, sdf.V2{x + ele, z + CAPWIDTH})
			dimensions = append(dimensions, sdf.V2{x, z + CAPWIDTH})
			x = x + gap + ele
		}*/
	return row
}

type keycap struct {
	x1, x2, z1, z2 float64
}

func (kc keycap) Dimensions() []sdf.V2 {
	return []sdf.V2{
		{kc.x1, kc.z1},
		{kc.x2, kc.z1},
		{kc.x2, kc.z2},
		{kc.x1, kc.z2},
	}
}

const (
	// plate
	PLATELENGTH = 285.0 //x
	PLATEHEIGHT = 1.2   //y
	PLATEWIDTH  = 94.6  //z
	// keyswitch
	SWITCHWIDTH  = 15 // square
	SWITCHLENGTH = 15 // square
	// keycap
	CAP1LENGTH  = 17.5 //x
	CAP15LENGTH = 17.5 * 1.5
	CAP2LENGTH  = 17.5 * 2
	CAPHEIGHT   = 1    //y for reference only
	CAPWIDTH    = 16.5 //z
	// row
	ROWCOUNT = 5
)

func keyCaps() []row {
	// this is about to get a lil silly
	sideWidthGaps := gapLength(PLATEWIDTH, (SWITCHWIDTH * ROWCOUNT), ROWCOUNT+1)                // 5 rows, so 6 gaps including exterior, so we add 1
	row1Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // we want consistent gaps on the outside, so remove those, we don't include exterior at all so we subtract 1.
	row2Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // split for unneccessary futureproofing
	row3Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row4Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row5Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (3*CAP2LENGTH)+(1*CAP15LENGTH)+(7*CAP1LENGTH)+(0.5*CAP1LENGTH), 11) // this one has an additional gapLength to split the arrow keys
	rows := make([]row, ROWCOUNT)

	// this can be turned into a loop
	rows[0] = newRow(sideWidthGaps, sideWidthGaps, row1Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	rows[1] = newRow(sideWidthGaps, (sideWidthGaps*2)+(CAPWIDTH*1), row2Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	rows[2] = newRow(sideWidthGaps, (sideWidthGaps*3)+(CAPWIDTH*2), row3Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	rows[3] = newRow(sideWidthGaps, (sideWidthGaps*4)+(CAPWIDTH*3), row4Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	rows[4] = newRow(sideWidthGaps, (sideWidthGaps*5)+(CAPWIDTH*4), row5Gaps, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH, CAP1LENGTH, CAP2LENGTH, CAP2LENGTH, CAP2LENGTH, CAP1LENGTH/2, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH)
	return rows
}

func drawKeyCaps() (sdf.SDF3, error) {
	keyCaps := keyCaps()
	var s sdf.SDF3
	caps := make([]sdf.SDF3, 0)
	for _, row := range keyCaps {
		for _, keycap := range row {
			cap := sdf.NewPolygon()
			cap.AddV2Set(keycap.Dimensions())
			profile, err := sdf.Polygon2D(cap.Vertices())
			if err != nil {
				return nil, err
			}
			caps = append(caps, sdf.Extrude3D(profile, CAPHEIGHT))
		}
	}
	s = caps[0]
	for i := range caps {
		s = sdf.Union3D(s, caps[i])
	}

	return s, nil
}

func plate() (sdf.SDF3, error) {
	plate, err := rect(0, 0, PLATELENGTH, PLATEHEIGHT, PLATEWIDTH)
	if err != nil {
		return nil, err
	}
	keycaps := keyCaps()
	for _, row := range keycaps {
		for _, cap := range row {
			x, z := (cap.x1 + (cap.x2-cap.x1)/2), (cap.z1 + (cap.z2-cap.z1)/2)
			keySwitch := sdf.NewBox2(sdf.V2{x, z}, sdf.V2{SWITCHLENGTH / 2, SWITCHWIDTH / 2})
			profile, err := sdf.Polygon2D(keySwitch.Vertices())
			if err != nil {
				return nil, err
			}
			keySwitchHole := sdf.Extrude3D(profile, PLATELENGTH)
			plate = sdf.Difference3D(plate, keySwitchHole)
		}
	}
	return plate, nil
}

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

	render.ToSTL(keycaps, 300, "keycaps.stl", &render.MarchingCubesOctree{})
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

func gapLength(plateLength float64, keyLength float64, gaps float64) float64 {
	return (plateLength - keyLength) / gaps
}

func createRow(inX, inZ, y, gap float64, keys ...float64) []sdf.V2 {
	dimensions := make([]sdf.V2, 0, 0)
	x, z := inX, inZ
	for _, ele := range keys {
		dimensions = append(dimensions, sdf.V2{x, z})
		dimensions = append(dimensions, sdf.V2{x + ele, z})
		dimensions = append(dimensions, sdf.V2{x + ele, z + CAPWIDTH})
		dimensions = append(dimensions, sdf.V2{x, z + CAPWIDTH})
		x = x + gap + ele
	}
	return dimensions
}

//-----------------------------------------------------------------------------
