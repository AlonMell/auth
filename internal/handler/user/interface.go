package user

import (
	del "github.com/AlonMell/ProviderHub/internal/handler/user/delete"
	"github.com/AlonMell/ProviderHub/internal/handler/user/get"
	"github.com/AlonMell/ProviderHub/internal/handler/user/post"
	"github.com/AlonMell/ProviderHub/internal/handler/user/update"
)

type Provider interface {
	del.Deleter
	post.Creater
	get.Getter
	update.Updater
}
