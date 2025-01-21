package drawing

var ElemTypeToRenderFunc = map[string]func(*Schematic, Point, Point){
	"line": func(s *Schematic, p1, p2 Point) {
		s.addPath(createPath(
			'M', p1.X, p1.Y,
			'L', p2.X, p2.Y,
		))
	},
	"resistor": func(s *Schematic, p1, p2 Point) {
		const (
			height = DefaultLength / 4
			width  = DefaultLength / 2
		)

		distance := p1.distanceTo(p2)

		s.addAndPivotPath(p1, p2, createPath(
			'M', p1.X, p1.Y,
			'L', p1.X+distance/2-width/2, p1.Y,
			'L', p1.X+distance/2-width/2, p1.Y+height/2,
			'L', p1.X+distance/2+width/2, p1.Y+height/2,
			'L', p1.X+distance/2+width/2, p1.Y-height/2,
			'L', p1.X+distance/2-width/2, p1.Y-height/2,
			'L', p1.X+distance/2-width/2, p1.Y,
			'M', p1.X+distance/2+width/2, p1.Y,
			'L', p1.X+distance, p1.Y,
		))
	},
	"battery": func(s *Schematic, p1, p2 Point) {
		const (
			termGap       = DefaultLength / 6
			negTermHeight = DefaultLength / 6
			posTermHeight = DefaultLength / 3
		)

		distance := p1.distanceTo(p2)
		negTermX := p1.X + distance/2 - termGap/2
		posTermX := negTermX + termGap

		s.addAndPivotPath(p1, p2, createPath(
			'M', p1.X, p1.Y,
			'L', negTermX, p1.Y,
			'M', negTermX, p1.Y-negTermHeight/2,
			'L', negTermX, p1.Y+negTermHeight/2,
			'M', posTermX, p1.Y-posTermHeight/2,
			'L', posTermX, p1.Y+posTermHeight/2,
			'M', posTermX, p1.Y,
			'L', p1.X+distance, p1.Y,
		))
	},
	"dot": func(s *Schematic, p1, p2 Point) {
		const radius = DefaultLength / 25
		s.addCircle(p1.X, p1.Y, radius)
	},
}
