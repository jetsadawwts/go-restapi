package ordersUsecases

import (
	"fmt"
	"math"

	"github.com/jetsadawwts/go-restapi/modules/entities"
	"github.com/jetsadawwts/go-restapi/modules/orders"
	"github.com/jetsadawwts/go-restapi/modules/orders/ordersRepositories"
	"github.com/jetsadawwts/go-restapi/modules/products/productsRepositories"
)

type IOrdersUsecase interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
	UpdateOrder(req *orders.Order) (*orders.Order, error)
}

type ordersUsecase struct {
	ordersRepository   ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(ordersRepository ordersRepositories.IOrdersRepository, productsRepository productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository:   ordersRepository,
		productsRepository: productsRepository,
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *ordersUsecase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.ordersRepository.FindOrder(req)
	return &entities.PaginateRes{
		Data:      orders,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	//Check if products is exists
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is nil")
		}

		prod, err := u.productsRepository.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, err
		}

		//Set price
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (u *ordersUsecase) UpdateOrder(req *orders.Order) (*orders.Order, error) {
	if err := u.ordersRepository.UpdateOrder(req); err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOneOrder(req.Id)
	if err != nil {
		return nil, err
	}

	return order, nil

}
