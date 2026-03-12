package order

type OrderStatus string

var (
	StatusPending   OrderStatus = "pending"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
)

type Order struct {
	ID                 int64       `json:"id"`
	UserID             int64       `json:"user_id"`
	Status             OrderStatus `json:"status"`
	TotalAmountInCents int32       `json:"total_amount_in_cents"`
	Items              []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ProductID            int64 `json:"product_id"`
	Quantity             int   `json:"quantity"`
	PurchasePriceInCents int32 `json:"purchase_price_in_cents"`
}
