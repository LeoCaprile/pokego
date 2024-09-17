package strimage

import (
	"image"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

func ToString(width int, img image.Image) string {
	img = imaging.Resize(img, width, 0, imaging.NearestNeighbor)
	b := img.Bounds()
	imageWidth := b.Max.X - 5
	h := b.Max.Y - 8
	str := strings.Builder{}

	for heightCounter := 8; heightCounter < h; heightCounter += 2 {

		for x := 5; x < imageWidth; x++ {
			c1, _ := colorful.MakeColor(img.At(x, heightCounter))
			color1 := lipgloss.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, heightCounter+1))
			color2 := lipgloss.Color(c2.Hex())
			str.WriteString(lipgloss.NewStyle().Foreground(color1).
				Background(color2).Render("▀"))
		}

		str.WriteString("\n")
	}

	return str.String()
}
