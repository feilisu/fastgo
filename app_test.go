package fastgo

import (
	"log"
	"testing"
)

func Test_app(t *testing.T) {

	defer func() {
		msg := recover()
		if msg != nil {
			log.Print(msg)
		}
	}()

	app := NewApp("测试应用1")
	app.router.Host("http://127.0.0.1").Path("/info").GET(func(ctx *Context) {
		err := ctx.Response.Text("你好")
		if err != nil {
			log.Println(err)
			return
		}
	})

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
