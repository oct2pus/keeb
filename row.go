package main

import "github.com/deadsy/sdfx/sdf"

const (
	CAP1LENGTH  = 17.5 //x
	CAP15LENGTH = 17.5 * 1.5
	CAP2LENGTH  = 17.5 * 2
	CAPHEIGHT   = 1    //y for reference only
	CAPWIDTH    = 16.5 //z
	ROWCOUNT    = 5
)

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

/*
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
} */

func gapLength(plateLength float64, keyLength float64, gaps float64) float64 {
	return (plateLength - keyLength) / gaps
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

func keyCaps() []row {
	// this is about to get a lil silly
	sideWidthGaps := gapLength(PLATEWIDTH-4.6, (SWITCHWIDTH * ROWCOUNT), ROWCOUNT+2)            // 5 rows, so 6 gaps including exterior, so we add 1
	row1Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // we want consistent gaps on the outside, so remove those, we don't include exterior at all so we subtract 1.
	row2Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP15LENGTH)+(12*CAP1LENGTH), 14-1) // split for unneccessary futureproofing
	row3Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row4Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (2*CAP2LENGTH)+(11*CAP1LENGTH), 13-1)
	row5Gaps := gapLength(PLATELENGTH-(sideWidthGaps*2), (3*CAP2LENGTH)+(1*CAP15LENGTH)+(7*CAP1LENGTH)+(0.5*CAP1LENGTH), 11) // this one has an additional gapLength to split the arrow keys
	rows := make([]row, ROWCOUNT)

	// this can be turned into a loop
	rows[0] = newRow(sideWidthGaps, 0, row1Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	rows[1] = newRow(sideWidthGaps, (sideWidthGaps*1)+(CAPWIDTH*1), row2Gaps, CAP15LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH)
	rows[2] = newRow(sideWidthGaps, (sideWidthGaps*2)+(CAPWIDTH*2), row3Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	rows[3] = newRow(sideWidthGaps, (sideWidthGaps*3)+(CAPWIDTH*3), row4Gaps, CAP2LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP2LENGTH)
	rows[4] = newRow(sideWidthGaps, (sideWidthGaps*4)+(CAPWIDTH*4), row5Gaps, CAP1LENGTH, CAP1LENGTH, CAP15LENGTH, CAP1LENGTH, CAP2LENGTH, CAP2LENGTH, CAP2LENGTH, CAP1LENGTH/2, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH, CAP1LENGTH)
	return rows
}
