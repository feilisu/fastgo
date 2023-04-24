package main

import (
	"fastgo"
	"log"
)

func main() {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Print(msg)
		}
	}()

	app := fastgo.NewApp("测试应用1")
	r := fastgo.NewRouter()
	r.Host("http://127.0.0.1").Path("/info").GET(func(ctx *fastgo.Context) {
		err := ctx.Response.Text("你好")
		if err != nil {
			log.Println(err)
			return
		}
	})
	app.SetRouter(r)

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
