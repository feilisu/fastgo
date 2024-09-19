package main

import (
	"github.com/feilisu/fastgo"
	"log"
	"strconv"
	"time"
)

func Mtest1(ctx *fastgo.Context) error {
	log.Print("Mtest1")
	return nil
}

type Mtest2 struct {
}

func (m *Mtest2) Exec(ctx *fastgo.Context) error {
	log.Print("Mtest2")
	return nil
}

type Json struct {
	Id   int64         `json:"id,omitempty"`
	Name string        `json:"name,omitempty"`
	Time time.Duration `json:"time,omitempty"`
}

func (receiver *GetUserName) Validate() error {
	fastgo.Logger().Info("GetUserName-Validate")
	return nil
}

type GetUserName struct {
	Id   int64         `json:"id,omitempty" default:"456" validate:"required,gt=1,lt=100,"`
	Name string        `json:"name,omitempty" validate:"required,lengthGt=1,lengthLt=10"`
	Time time.Duration `json:"time,omitempty"`
	Json Json          `json:"json,omitempty" validate:"required,lengthGt=1,lengthLt=10"`
}

type GetUserNameV2 struct {
	Test GetUserName `json:"test,omitempty" `
}

type User struct {
}

func (u User) GetName(ctx *fastgo.Context) error {
	param := new(GetUserName)
	err := ctx.Request.Params(param)
	if err != nil {
		fastgo.Logger().Info(err.Error())
	}
	return ctx.Response.Json(param)
}

func main() {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Print(msg)
		}
	}()

	server := fastgo.DefaultServer()
	server.Port = strconv.Itoa(11220)

	r := fastgo.NewRouter()
	r.Host(fastgo.DefaultHost).Path("/info").GET(func(ctx *fastgo.Context) error {
		err := ctx.Response.Text("你好")
		log.Println(ctx.Value("RemoteAddr"))
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	type TestResult struct {
		Result string
	}

	r.Host(fastgo.DefaultHost).Path("/icbcnewmis").POST(func(ctx *fastgo.Context) error {
		s := new(TestResult)
		s.Result = "1"
		err := ctx.Response.Json(s)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	r.Host("localhost").Path("/user/name").POST(User{}.GetName)
	err := server.Run(r)
	if err != nil {
		log.Println(err)
		return
	}

}

//func main() {
//	fastgo.Test("123")
//	fastgo.Test(1)
//	fastgo.Test(GetUserNameV2{Test: GetUserName{Id: 1, Name: "test"}})
//}

//
//package main
//
//import (
//"fmt"
//"reflect"
//)
//
//type I interface {
//	M()
//}
//
//type T struct{}
//
//func (t T) M() {}
//
//func main() {
//	//var t T
//	t := new(T)
//	fmt.Println(reflect.TypeOf(t).Implements(reflect.TypeOf((*I)(nil)).Elem()))
//}
