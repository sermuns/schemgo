package drawing

func (s *Schematic) line(p1, p2 Point) {
	s.addPath(createPath(
		'M', p1.x, p1.y,
		'L', p2.x, p2.y,
	))
}

func (s *Schematic) resistor(p1, p2 Point) {
	const (
		height = defaultLength / 4
		width  = defaultLength / 2
	)

	distance := p1.distanceTo(p2)

	s.addAndPivotPath(p1, p2, createPath(
		'M', p1.x, p1.y,
		'L', p1.x+distance/2-width/2, p1.y,
		'L', p1.x+distance/2-width/2, p1.y+height/2,
		'L', p1.x+distance/2+width/2, p1.y+height/2,
		'L', p1.x+distance/2+width/2, p1.y-height/2,
		'L', p1.x+distance/2-width/2, p1.y-height/2,
		'L', p1.x+distance/2-width/2, p1.y,
		'M', p1.x+distance/2+width/2, p1.y,
		'L', p1.x+distance, p1.y,
	))
}

func (s *Schematic) dot(p1 Point) {
	const radius = defaultLength / 25
	s.addCircle(p1.x, p1.y, radius)
}

func (s *Schematic) battery(p1, p2 Point) {
	const (
		termGap       = defaultLength / 6
		negTermHeight = defaultLength / 6
		posTermHeight = defaultLength / 3
	)

	distance := p1.distanceTo(p2)
	negTermX := p1.x + distance/2 - termGap/2
	posTermX := negTermX + termGap

	s.addAndPivotPath(p1, p2, createPath(
		'M', p1.x, p1.y,
		'L', negTermX, p1.y,
		'M', negTermX, p1.y-negTermHeight/2,
		'L', negTermX, p1.y+negTermHeight/2,
		'M', posTermX, p1.y-posTermHeight/2,
		'L', posTermX, p1.y+posTermHeight/2,
		'M', posTermX, p1.y,
		'L', p1.x+distance, p1.y,
	))
}

func (s *Schematic) capacitor(p1, p2 Point) {
	const (
		gap    = defaultLength / 6
		height = defaultLength / 3
	)

	distance := p1.distanceTo(p2)
	negTermX := p1.x + distance/2 - gap/2
	posTermX := negTermX + gap

	s.addAndPivotPath(p1, p2, createPath(
		'M', p1.x, p1.y,
		'L', negTermX, p1.y,
		'M', negTermX, p1.y-height/2,
		'L', negTermX, p1.y+height/2,
		'M', posTermX, p1.y-height/2,
		'L', posTermX, p1.y+height/2,
		'M', posTermX, p1.y,
		'L', p1.x+distance, p1.y,
	))
}
