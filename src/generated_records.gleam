// Generated from Go records

pub type Point {
  Point(x: Float, y: Float)
}

pub type Person {
  Person(name: String, age: Int)
}

pub fn move(p: Point, dx: Float, dy: Float) -> Point {
  Point(x: p.x +. dx, y: p.y +. dy)
}

pub fn birthday(p: Person) -> Person {
  Person(..p, age: p.age + 1)
}

