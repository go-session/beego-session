package beegosession

import (
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
	storeKey  struct{}
	manageKey struct{}
)

var (
	// DefaultConfig is the default Recover middleware config.
	DefaultConfig = Config{
		ErrorHandleFunc: func(ctx *context.Context, err error) {
			ctx.Abort(500, err.Error())
		},
	}
)

// New create a session middleware
func New(opt ...session.Option) beego.FilterFunc {
	return NewWithConfig(DefaultConfig, opt...)
}

// NewWithConfig create a session middleware
func NewWithConfig(config Config, opt ...session.Option) beego.FilterFunc {
	if config.ErrorHandleFunc == nil {
		config.ErrorHandleFunc = DefaultConfig.ErrorHandleFunc
	}

	manage := session.NewManager(opt...)
	return func(ctx *context.Context) {
		ctx.Input.SetData(manageKey{}, manage)
		store, err := manage.Start(nil, ctx.ResponseWriter, ctx.Request)
		if err != nil {
			config.ErrorHandleFunc(ctx, err)
			return
		}
		ctx.Input.SetData(storeKey{}, store)
	}
}

// FromContext Get session storage from context
func FromContext(ctx *context.Context) session.Store {
	return ctx.Input.GetData(storeKey{}).(session.Store)
}

// Destroy a session
func Destroy(ctx *context.Context) error {
	return ctx.Input.GetData(storeKey{}).(*session.Manager).
		Destroy(nil, ctx.ResponseWriter, ctx.Request)
}

// Refresh a session and return to session storage
func Refresh(ctx *context.Context) (session.Store, error) {
	return ctx.Input.GetData(storeKey{}).(*session.Manager).
		Refresh(nil, ctx.ResponseWriter, ctx.Request)
}
