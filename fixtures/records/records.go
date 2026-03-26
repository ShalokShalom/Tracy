package records

type Point struct {
	X float64
	Y float64
}

type Person struct {
	Name string
	Age  int
}

func Move(p Point, dx float64, dy float64) Point {
	p.X = p.X + dx
	p.Y = p.Y + dy
	return p
}

func Birthday(p Person) Person {
	p.Age = p.Age + 1
	return p
}