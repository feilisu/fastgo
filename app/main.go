package main

import (
	"fastgo"
	"log"
)

type User struct {
}

func (u User) GetName(ctx *fastgo.Context) error {
	return ctx.Response.Json("费力苏")
}

func main() {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Print(msg)
		}
	}()

	app := fastgo.NewApp("测试应用1")
	app.SetMiddleware([]fastgo.Middleware{
		fastgo.MiddlewareFunc(fastgo.Mtest1),
		new(fastgo.Mtest2),
	})
	r := fastgo.NewRouter()
	r.Host("127.0.0.1").Path("/info").GET(fastgo.HandlerFunc(func(ctx *fastgo.Context) error {
		err := ctx.Response.Text("你好")
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}))

	r.Host("127.0.0.1").Path("/user/name").GET(fastgo.HandlerFunc(User{}.GetName))
	app.SetRouter(r)

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
