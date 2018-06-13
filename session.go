package beegosession

import (
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/go-session/session"
)

type (
	// ErrorHandleFunc error handling function
	ErrorHandleFunc func(*context.Context, error)
	// Config defines the config for Session middleware
	Config struct {
		// error handling when starting the session
		ErrorHandleFunc ErrorHandleFunc
	}
	storeKey struct{}
)

var (
	once            sync.Once
	internalManager *session.Manager

	// DefaultConfig is the default Recover middleware config.
	DefaultConfig = Config{
		ErrorHandleFunc: func(ctx *context.Context, err error) {
			ctx.Abort(500, err.Error())
		},
	}
)

func manager(opt ...session.Option) *session.Manager {
	once.Do(func() {
		internalManager = session.NewManager(opt...)
	})
	return internalManager
}

// New create a session middleware
func New(opt ...session.Option) beego.FilterFunc {
	return NewWithConfig(DefaultConfig, opt...)
}

// NewWithConfig create a session middleware
func NewWithConfig(config Config, opt ...session.Option) beego.FilterFunc {
	if config.ErrorHandleFunc == nil {
		config.ErrorHandleFunc = DefaultConfig.ErrorHandleFunc
	}

	return func(ctx *context.Context) {
		store, err := manager(opt...).Start(nil, ctx.ResponseWriter, ctx.Request)
		if err != nil {
			config.ErrorHandleFunc(ctx, err)
			return
		}
		ctx.Input.SetData(storeKey{}, store)
	}
}

// FromContext Get session storage from context
func FromContext(ctx *context.Context) session.Store {
	store := ctx.Input.GetData(storeKey{})
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
