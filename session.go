package beegosession

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gopkg.in/session.v2"
)

var (
	once            sync.Once
	internalManager *session.Manager
	internalError   ErrorHandleFunc
)

type sessionKey struct{}

func init() {
	internalError = func(ctx *context.Context, err error) {
		ctx.Abort(500, err.Error())
	}
}

// ErrorHandleFunc error handling function
type ErrorHandleFunc func(ctx *context.Context, err error)

// SetErrorHandler Set error handling
func SetErrorHandler(handler ErrorHandleFunc) {
	internalError = handler
}

func manager(opt ...session.Option) *session.Manager {
	once.Do(func() {
		internalManager = session.NewManager(opt...)
	})
	return internalManager
}

// New Create a session middleware
func New(opt ...session.Option) beego.FilterFunc {
	return func(ctx *context.Context) {
		store, err := manager(opt...).Start(nil, ctx.ResponseWriter, ctx.Request)
		if err != nil {
			internalError(ctx, err)
			return
		}
		ctx.Input.SetData(sessionKey{}, store)
	}
}

// FromContext Get session storage from context
func FromContext(ctx *context.Context) session.Store {
	store := ctx.Input.GetData(sessionKey{})
	if store != nil {
		return store.(session.Store)
	}
	return nil
}

// Destroy a session
func Destroy(ctx *context.Context) error {
	return manager().Destroy(nil, ctx.ResponseWriter, ctx.Request)
}

// Refresh a session and return to session storage
func Refresh(ctx *context.Context) (session.Store, error) {
	return manager().Refresh(nil, ctx.ResponseWriter, ctx.Request)
}
