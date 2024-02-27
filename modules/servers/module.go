package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresHandlers"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresRepositories"
	"github.com/jetsadawwts/go-restapi/modules/middlewares/middlewaresUsecases"
	"github.com/jetsadawwts/go-restapi/modules/monitors/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
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
