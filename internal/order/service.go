package order

import (
	"github.com/moges7624/merkato_std/internal/product"
	"github.com/moges7624/merkato_std/internal/user"
)

type Service struct {
	store          Store
	productService product.Service
	userService    user.Service
}

func NewService(
	store Store,
	productService product.Service,
	userService user.Service,
) *Service {
	return &Service{
		store:          store,
		productService: productService,
		userService:    userService,
	}
}

func (s *Service) GetOrders() ([]*Order, error) {
	return s.store.getAll()
}

func (s *Service) GetOrderByID(id int64) (*Order, error) {
	return s.store.getByID(id)
}

func (s *Service) CreateOrder(req *CreateOrderRequest) (*Order, error) {
	_, err := s.userService.GetUser(int(*req.UserID))
	if err != nil {
		return nil, err
	}
	totalPriceInCents := 0
	items := make([]OrderItem, 0, len(*req.Items))

	for _, item := range *req.Items {
		_, err := s.productService.GetProduct(*item.ProductID)
		if err != nil {
			return nil, err
		}

		err = s.productService.DeductStock(*item.ProductID, *item.Quantity)
		if err != nil {
			return nil, err
		}

		priceInCent := int32(*item.PriceInUSD) * 100

		items = append(items, OrderItem{
			ProductID:            *item.ProductID,
			Quantity:             *item.Quantity,
			PurchasePriceInCents: priceInCent,
		})

		totalPriceInCents += int(priceInCent) * *item.Quantity
	}

	order := &Order{
		UserID:             *req.UserID,
		Status:             StatusPending,
		Items:              items,
		TotalAmountInCents: int32(totalPriceInCents),
	}

	err = s.store.insert(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}
