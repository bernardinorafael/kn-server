package contract

type AuthService interface {
	Login(email, password string) error
	Register(name, email, password string) error
}
