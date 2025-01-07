package servers

import (
	"github.com/gofiber/fiber/v2"

	appinfohandlers "github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoHandlers"
	appinforepositories "github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoRepositories"
	appinfousecases "github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoUsecases"

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
	AppinfoModule()
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
	handler := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	respository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, respository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.m.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signup-admin", m.m.JwtAuth(), m.m.Authorize(2), handler.SignUpAdmin)
	router.Post("/signin", m.m.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.m.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.m.ApiKeyAuth(), handler.SignOut)

	router.Get("/admin/secret", m.m.JwtAuth(), m.m.Authorize(2), handler.GenerateAdminToken)
	router.Get("/:user_id", m.m.JwtAuth(), m.m.ParamsCheck(), handler.GetUserProfile)

}

func (m *moduleFactory) AppinfoModule() {
	respository := appinforepositories.AppinfoRepository(m.s.db)
	usecase := appinfousecases.AppinfoUsecase(respository)
	handler := appinfohandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")
	router.Get("/apikey", m.m.JwtAuth(), m.m.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.m.ApiKeyAuth(), handler.FindCategory)

	router.Post("/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.RemoveCategory)

}
