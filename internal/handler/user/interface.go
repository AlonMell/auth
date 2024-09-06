package user

import (
	del "providerHub/internal/handler/user/delete"
	"providerHub/internal/handler/user/get"
	"providerHub/internal/handler/user/post"
	"providerHub/internal/handler/user/update"
)

type Interface interface {
	del.Deleter
	post.Creater
	get.Getter
	update.Updater
}
