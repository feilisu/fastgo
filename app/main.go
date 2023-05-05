package main

import (
	"fastgo"
	"log"
)

type GetUserName struct {
	Id   int64  `json:"id,omitempty" default:"456" validate:"required,max=10,mix=1"`
	Name string `json:"name,omitempty"`
}

type GetUserNameV2 struct {
	Test GetUserName `json:"test,omitempty" `
}

type User struct {
}

func (u User) GetName(ctx *fastgo.Context) error {
	param := new(GetUserNameV2)
	err := ctx.Request.Params(param)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(param)
	return ctx.Response.Json("qwe")
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

	r.Host("127.0.0.1").Path("/user/name").POST(fastgo.HandlerFunc(User{}.GetName))
	app.SetRouter(r)

	err := app.Run()
	if err != nil {
		log.Println(err)
		return
	}
}
