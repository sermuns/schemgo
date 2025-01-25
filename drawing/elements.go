package drawing

import "fmt"

var (
	ElemTypeToRenderFunc = map[string]func(*Schematic, Point, Point){
		"line": func(s *Schematic, p1, p2 Point) {
			s.addPath(createPath(
				'M', Point{p1.X, p1.Y},
				'L', Point{p2.X, p2.Y},
			))
		},
		"resistor": func(s *Schematic, p1, p2 Point) {
			const (
				height = DefaultLength / 4
				width  = DefaultLength / 2
			)

			distance := p1.distanceTo(p2)

			s.addAndPivotPath(p1, p2, createPath(
				'M', Point{p1.X, p1.Y},
				'L', Point{p1.X + distance/2 - width/2, p1.Y},
				'L', Point{p1.X + distance/2 - width/2, p1.Y + height/2},
				'L', Point{p1.X + distance/2 + width/2, p1.Y + height/2},
				'L', Point{p1.X + distance/2 + width/2, p1.Y - height/2},
				'L', Point{p1.X + distance/2 - width/2, p1.Y - height/2},
				'L', Point{p1.X + distance/2 - width/2, p1.Y},
				'M', Point{p1.X + distance/2 + width/2, p1.Y},
				'L', Point{p1.X + distance, p1.Y},
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
				'M', Point{p1.X, p1.Y},
				'L', Point{negTermX, p1.Y},
				'M', Point{negTermX, p1.Y - negTermHeight/2},
				'L', Point{negTermX, p1.Y + negTermHeight/2},
				'M', Point{posTermX, p1.Y - posTermHeight/2},
				'L', Point{posTermX, p1.Y + posTermHeight/2},
				'M', Point{posTermX, p1.Y},
				'L', Point{p1.X + distance, p1.Y},
			))
		},
		"dot": func(s *Schematic, p1, p2 Point) {
			const radius = DefaultLength / 25
			s.addCircle(p1.X, p1.Y, radius)
		},
		"capacitor": func(s *Schematic, p1, p2 Point) {
			const (
				gap    = DefaultLength / 6
				height = DefaultLength / 3
			)

			distance := p1.distanceTo(p2)
			negTermX := p1.X + distance/2 - gap/2
			posTermX := negTermX + gap

			s.addAndPivotPath(p1, p2, createPath(
				'M', Point{p1.X, p1.Y},
				'L', Point{negTermX, p1.Y},
				'M', Point{negTermX, p1.Y - height/2},
				'L', Point{negTermX, p1.Y + height/2},
				'M', Point{posTermX, p1.Y - height/2},
				'L', Point{posTermX, p1.Y + height/2},
				'M', Point{posTermX, p1.Y},
				'L', Point{p1.X + distance, p1.Y},
			))
		},
		"inductor": func(s *Schematic, p1, p2 Point) {
			const (
				width        = DefaultLength
				coilLoopSize = width / 8
			)

			distance := p1.distanceTo(p2)

			s.addAndPivotPath(p1, p2, createPath(
				'M', Point{p1.X, p1.Y},
				'l', Point{distance/2 - width/2, 0},

				'q',
				Point{0, -coilLoopSize},
				Point{coilLoopSize, -coilLoopSize},
				'q',
				Point{coilLoopSize, 0},
				Point{coilLoopSize, coilLoopSize},
				'q',
				Point{0, -coilLoopSize},
				Point{coilLoopSize, -coilLoopSize},
				'q',
				Point{coilLoopSize, 0},
				Point{coilLoopSize, coilLoopSize},
				'q',
				Point{0, -coilLoopSize},
				Point{coilLoopSize, -coilLoopSize},
				'q',
				Point{coilLoopSize, 0},
				Point{coilLoopSize, coilLoopSize},
				'q',
				Point{0, -coilLoopSize},
				Point{coilLoopSize, -coilLoopSize},
				'q',
				Point{coilLoopSize, 0},
				Point{coilLoopSize, coilLoopSize},

				'L', Point{p1.X + distance, p1.Y},
			))
		},
		"sourceV": func(s *Schematic, p1, p2 Point) {
			const (
				radius = DefaultLength / 8
			)

			// distance := p1.distanceTo(p2)

			circleStyle := fmt.Sprintf(`fill="none" stroke="black" stroke-width="%d"`, DefaultStrokeWidth)
			s.addCircle(p1.X, p1.Y, radius, circleStyle)

			// s.addAndPivotPath(p1, p2, createPath(
			// 	'M', Point{p1.X, p1.Y},
			// 	'l', Point{distance/2 - width/2, 0},
			//
			// 	'q',
			// 	Point{0, -coilLoopSize},
			// 	Point{coilLoopSize, -coilLoopSize},
			// 	'q',
			// 	Point{coilLoopSize, 0},
			// 	Point{coilLoopSize, coilLoopSize},
			// 	'q',
			// 	Point{0, -coilLoopSize},
			// 	Point{coilLoopSize, -coilLoopSize},
			// 	'q',
			// 	Point{coilLoopSize, 0},
			// 	Point{coilLoopSize, coilLoopSize},
			// 	'q',
			// 	Point{0, -coilLoopSize},
			// 	Point{coilLoopSize, -coilLoopSize},
			// 	'q',
			// 	Point{coilLoopSize, 0},
			// 	Point{coilLoopSize, coilLoopSize},
			// 	'q',
			// 	Point{0, -coilLoopSize},
			// 	Point{coilLoopSize, -coilLoopSize},
			// 	'q',
			// 	Point{coilLoopSize, 0},
			// 	Point{coilLoopSize, coilLoopSize},
			//
			// 	'L', Point{p1.X + distance, p1.Y},
			// ))
		},
	}

	CommandTypeToFunc = map[string]func(*Schematic){
		"push": func(s *Schematic) {
			s.Push(s.Pos)
		},
		"pop": func(s *Schematic) {
			s.Pos = s.Pop()
		},
	}

	MovementFuncs = map[string]func(*Schematic, float64){
		"right": func(s *Schematic, value float64) {
			s.Translate(value, 0)
		},
		"up": func(s *Schematic, value float64) {
			s.Translate(0, -value)
		},
		"down": func(s *Schematic, value float64) {
			s.Translate(0, value)
		},
		"left": func(s *Schematic, value float64) {
			s.Translate(-value, 0)
		},
	}
)
