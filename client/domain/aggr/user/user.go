package user

type User struct {
	Name     string
	Password string
	Token    string
}

func NewUser(username, password string) *User {
	return &User{
		Name:     username,
		Password: password,
		Token:    "",
	}
}
