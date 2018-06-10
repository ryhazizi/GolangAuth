package model

type User struct {
	ID       string `bson:"id"`
	FullName string `bson:"fullname"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

type Users []User
