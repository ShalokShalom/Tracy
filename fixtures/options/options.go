package options

type User struct {
    Name string
}

func GetUser(id int) *User {
    if id <= 0 {
        return nil
    }
    u := User{Name: "User123"}
    return &u
}

func GreetUser(user *User) string {
    if user == nil {
        return "Hello, stranger"
    }
    return "Hello, " + user.Name
}

func FirstUser(a *User, b *User) *User {
    if a != nil {
        return a
    }
    return b
}

func IncrementMaybeAge(age *int) *int {
    if age == nil {
        return nil
    }
    result := *age + 1
    return &result
}

func FindInt(items []int, target int) *int {
    for _, x := range items {
        if x == target {
            result := x
            return &result
        }
    }
    return nil
}