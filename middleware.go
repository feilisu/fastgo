package fastgo

type MiddlewareParam struct {
}

type MiddlewareRes struct {
}

type Middleware interface {
	Run(p MiddlewareParam) (res MiddlewareRes)
}

var (
	Middlewares []Middleware
)

func MiddlewareRun(p MiddlewareParam) {
	for _, m := range Middlewares {
		m.Run(p)
	}
}
