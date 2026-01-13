package main

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var OurUser = User{
	Username: "John Doe",
	Password: "password", // yes I know, very strong password
}
