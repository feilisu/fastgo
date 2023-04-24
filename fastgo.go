package fastgo

type App struct {
	appId       string
	router      *Router
	middlewares []Middleware
	server      *Server
}

func NewApp(appId string) *App {
	return &App{
		appId:  appId,
		router: NewRouter(),
		server: NewServer(),
	}
}

func (a *App) SetMiddleware(ws []Middleware) {
	a.middlewares = ws
}
func (a *App) SetServer(s *Server) {
	a.server = s
}

func (a *App) SetRouter(r *Router) {
	a.router = r
}

func (a *App) Run() error {
	return a.server.Run(a.router.serveMux)
}
