package view

import (
	"strconv"

	"github.com/fogleman/gg"
)

// GenImg draws a simple top down image of the simulation.
func GenImg(waypoints [][]float64, agents [][]float64, tick int) {
	filename := "img/" + strconv.Itoa(tick)

	dc := gg.NewContext(400, 400)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.MoveTo(waypoints[0][0], waypoints[0][1])
	for _, waypoint := range waypoints {
		dc.LineTo(waypoint[0], waypoint[1])
	}
	dc.SetHexColor("FF2D00")
	dc.SetLineWidth(1)
	dc.Stroke()
	for i, agent := range agents {
		switch i % 3 {
		case 0:
			dc.SetRGB(0, 0, 1)
		case 1:
			dc.SetRGB(0, 1, 0)
		case 2:
			dc.SetRGB(1, 0, 0)
		}

		dc.DrawPoint(agent[0], agent[1], 4)
		dc.Fill()
	}
	dc.SavePNG(filename + ".png")
}
