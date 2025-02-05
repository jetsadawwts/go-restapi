package servers

import (
	"github.com/gofiber/fiber/v2"

	appinfohandlers "github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoHandlers"
	"github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoRepositories"
	"github.com/jetsadawwts/go-restapi/modules/appinfo/appinfoUsecases"
	"github.com/jetsadawwts/go-restapi/modules/orders/ordersHandlers"
	"github.com/jetsadawwts/go-restapi/modules/orders/ordersRepositories"
	"github.com/jetsadawwts/go-restapi/modules/orders/ordersUsecases"

	"github.com/jetsadawwts/go-restapi/modules/files/filesHandlers"
	"github.com/jetsadawwts/go-restapi/modules/files/filesUsecases"

	"github.com/jetsadawwts/go-restapi/modules/products/productsHandlers"
	"github.com/jetsadawwts/go-restapi/modules/products/productsRepositories"
	"github.com/jetsadawwts/go-restapi/modules/products/productsUsecases"

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
	FilesModule()
	ProductsModule()
	OrdersModule()
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
	respository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(respository)
	handler := appinfohandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")
	router.Get("/apikey", m.m.JwtAuth(), m.m.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.m.ApiKeyAuth(), handler.FindCategory)

	router.Post("/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.RemoveCategory)

}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)
	router := m.r.Group("/files")

	router.Post("/upload", m.m.JwtAuth(), m.m.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.m.JwtAuth(), m.m.Authorize(2), handler.DeleteFiles)

}

func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)

	productsRespository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)
	productsUsecase := productsUsecases.ProductsUsecase(productsRespository)
	productsHandler := productsHandlers.ProductsHandler(m.s.cfg, productsUsecase, filesUsecase)

	router := m.r.Group("/products")
	
	router.Get("/", m.m.ApiKeyAuth(), productsHandler.FindProduct)
	router.Get("/:product_id", m.m.ApiKeyAuth(), productsHandler.FindOneProduct)
	router.Post("/", m.m.JwtAuth(), m.m.Authorize(2), productsHandler.AddProduct)
	router.Patch("/:product_id", m.m.JwtAuth(), m.m.Authorize(2), productsHandler.UpdateProduct)
	router.Delete("/:product_id", m.m.JwtAuth(), m.m.Authorize(2), productsHandler.DeleteProduct)

}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)
	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrdersHandler(m.s.cfg, ordersUsecase)

	router := m.r.Group("/orders")

	router.Get("/:user_id/:order_id", m.m.JwtAuth(), m.m.ParamsCheck(), ordersHandler.FindOneOrder)
	router.Get("/", m.m.JwtAuth(), m.m.Authorize(2), ordersHandler.FindOrder)
	router.Post("/", m.m.JwtAuth(), ordersHandler.InsertOrder)
	router.Patch("/:user_id/:order_id", m.m.JwtAuth(), m.m.ParamsCheck(), ordersHandler.UpdateOrder)

}
