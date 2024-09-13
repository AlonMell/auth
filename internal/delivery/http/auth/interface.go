package auth

type Auth interface {
	UserProvider
	UserRefresher
	UserRegister
}
