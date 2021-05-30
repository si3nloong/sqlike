# Create

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

### Batch Insert

```go
var users := []User{
    User{Name: "User A", Age: 20},
    User{Name: "User B", Age: 40},
}

db.Table("users").Insert(context.Background(), &users)
```
