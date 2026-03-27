// Generated from Go records

pub type Point {
  X Int
  Y Int
}

pub type Person {
  Name String
  Age String
}

pub fn Birthday(p: codeberg.org/shalokshalom/Tracy/fixtures/records.Person) -> codeberg.org/shalokshalom/Tracy/fixtures/records.Person {
  codeberg.org/shalokshalom/Tracy/fixtures/records.Person(..p, Age: (t2 + 1))
}

