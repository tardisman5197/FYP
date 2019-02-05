package view

import (
	"strconv"

	"github.com/fogleman/gg"
)

// GenImg draws a simple top down image of the simulation.
func GenImg(waypoints [][]float64, agents [][]float64, tick int) {
	filename := "img/" + strconv.Itoa(tick)

	dc := gg.NewContext(300, 300)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(1)
	dc.MoveTo(waypoints[0][0]*10, waypoints[0][1]*10)
	for _, waypoint := range waypoints {
		dc.LineTo(waypoint[0]*10, waypoint[1]*10)
	}
	dc.SetHexColor("FF2D00")
	dc.SetLineWidth(1)
	dc.Stroke()
	dc.SetRGB(0, 0, 1)
	for _, agent := range agents {
		dc.DrawPoint(agent[0]*10, agent[1]*10, 2)
		dc.Fill()
	}
	dc.SavePNG(filename + ".png")
}
