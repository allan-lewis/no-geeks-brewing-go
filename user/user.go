package user

type User struct {
	authenticated bool
	name          string
	email         string
}

func UnauthenticatedUser() User {
	return User{
		authenticated: false,
		name:          "",
		email:         "",
	}
}

func AuthenticatedUser(name string, email string) User {
	return User{
		authenticated: true,
		name:          name,
		email:         email,
	}
}

func (u *User) Authenticated() bool {
	return u.authenticated
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}
