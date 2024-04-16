package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresHandlers"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresRepositories"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresUsecases"
	"github.com/jetsadawwts/go-restapi/modules/monitors/monitorHandlers"
	"github.com/jetsadawwts/go-restapi/modules/users/usersHandlers"
	"github.com/jetsadawwts/go-restapi/modules/users/usersRepositories"
	"github.com/jetsadawwts/go-restapi/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
	m middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, m middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
		m: m,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	respository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(respository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)

}

func (m *moduleFactory) MonitorModule() {
	handle := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handle.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	respository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, respository)
	handle := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", handle.SignUpCustomer)
	router.Post("/signin", handle.SignIn)
	router.Post("/refresh", handle.RefreshPassport)
}
