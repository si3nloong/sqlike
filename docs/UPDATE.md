# Update

### Create Single Record

```go
var user User{
    Name string
    Age  int
}

user.Name = "John Cena"
user.Age = 24

db.Table("users").InsertOne(context.Background(), &user)
```
