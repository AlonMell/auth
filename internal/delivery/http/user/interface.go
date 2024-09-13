package user

type Provider interface {
	Deleter
	Creater
	Getter
	Updater
}
