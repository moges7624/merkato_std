package order

import (
	"errors"
	"sync"
)

type FileStore struct {
	mu     sync.RWMutex
	orders map[int64]Order
	idSeq  int64
}

func NewFileStore() *FileStore {
	return &FileStore{
		orders: map[int64]Order{
			1: {
				ID:                 1,
				UserID:             1,
				Status:             StatusPending,
				TotalAmountInCents: 23999,
				Items: []OrderItem{
					{
						ProductID:            1,
						Quantity:             2,
						PurchasePriceInCents: 11598,
					},
					{
						ProductID:            2,
						Quantity:             1,
						PurchasePriceInCents: 7598,
					},
				},
			},
			2: {
				ID:                 2,
				UserID:             1,
				Status:             StatusShipped,
				TotalAmountInCents: 97329,
			},
		},
		idSeq: 2,
	}
}

func (fs *FileStore) getAll() ([]*Order, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	orders := make([]*Order, 0, len(fs.orders))
	for _, order := range fs.orders {
		orders = append(orders, &order)
	}

	return orders, nil
}

func (fs *FileStore) getByID(id int64) (*Order, error) {
	return nil, errors.New("not implemented")
}

func (fs *FileStore) insert(order *Order) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.idSeq++
	order.ID = fs.idSeq
	fs.orders[order.ID] = *order

	return nil
}
