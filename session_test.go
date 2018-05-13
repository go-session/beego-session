package beegosession

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gopkg.in/session.v2"
)

func TestSession(t *testing.T) {
	cookieName := "test_beego_session"

	app := beego.NewApp()

	app.Handlers.InsertFilter("*", beego.BeforeRouter, New(
		session.SetCookieName(cookieName),
		session.SetSign([]byte("sign")),
	))

	app.Handlers.Get("/", func(ctx *context.Context) {
		store := FromContext(ctx)
		if ctx.Input.Query("login") == "1" {
			foo, ok := store.Get("foo")
			fmt.Fprintf(ctx.ResponseWriter, "%s:%v", foo, ok)
			return
		}

		store.Set("foo", "bar")
		err := store.Save()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Fprint(ctx.ResponseWriter, "ok")
	})

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.Handlers.ServeHTTP(w, req)

	res := w.Result()
	cookie := res.Cookies()[0]
	if cookie.Name != cookieName {
		t.Error("Not expected value:", cookie.Name)
		return
	}

	buf, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(buf) != "ok" {
		t.Error("Not expected value:", string(buf))
		return
	}

	req, err = http.NewRequest("GET", "/?login=1", nil)
	if err != nil {
		t.Error(err)
		return
	}
	req.AddCookie(cookie)

	w = httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, req)

	res = w.Result()
	buf, _ = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if string(buf) != "bar:true" {
		t.Error("Not expected value:", string(buf))
		return
	}
}
