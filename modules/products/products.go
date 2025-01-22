package products

import (
	"github.com/jetsadawwts/go-restapi/modules/appinfo"
	"github.com/jetsadawwts/go-restapi/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreateAt    string            `json:"created_at"`
	UpdateAt    string            `json:"update_at"`
	Price       float64           `json:"price"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id string `query:"id"`
	Search string `query:"search"` //title & description
	*entities.PaginationReq
	*entities.SortReq
}
