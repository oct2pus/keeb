package main

import "github.com/deadsy/sdfx/sdf"

const (
	// plate
	PLATELENGTH = 285.0 //x
	PLATEHEIGHT = 1.2   //y
	PLATEWIDTH  = 94.6  //z
	// keyswitch
	SWITCHWIDTH  = 15 // square
	SWITCHLENGTH = 15 // square
)

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
			p := sdf.NewPolygon()
			p.AddV2(keySwitch.Min)
			p.AddV2(sdf.V2{keySwitch.Max.X, keySwitch.Min.Y})
			p.AddV2(keySwitch.Max)
			p.AddV2(sdf.V2{keySwitch.Min.X, keySwitch.Max.Y})
			profile, err := sdf.Polygon2D(p.Vertices())
			if err != nil {
				return nil, err
			}
			keySwitchHole := sdf.Extrude3D(profile, PLATELENGTH)
			plate = sdf.Difference3D(plate, keySwitchHole)
		}
	}
	return plate, nil
}
