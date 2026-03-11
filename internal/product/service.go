package product

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetProducts() (*[]Product, error) {
	products, err := s.store.getProducts()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Service) GetProduct(id int64) (*Product, error) {
	prod, err := s.store.getProduct(id)
	if err != nil {
		return nil, err
	}

	return prod, nil
}

func (s *Service) CreateProducts(cpr *CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:         *cpr.Name,
		PriceInCents: int32(*cpr.PriceInDollar * 100),
		Quantity:     *cpr.Qauntity,
	}
	err := s.store.createProduct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}
