pub type Point {
  Point(x: Float, y: Float)
}

pub fn move(p: Point, dx: Float, dy: Float) -> Point {
  Point(x: p.x +. dx, y: p.y +. dy)
}

pub type Person {
  Person(name: String, age: Int)
}

pub fn birthday(p: Person) -> Person {
  Person(..p, age: p.age + 1)
}

pub type Config {
  Config(host: String, port: Int, debug: Bool)
}

pub fn with_debug(c: Config) -> Config {
  Config(..c, debug: True)
}

pub fn with_port(c: Config, port: Int) -> Config {
  Config(..c, port: port)
}
