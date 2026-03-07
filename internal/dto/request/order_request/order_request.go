package orderrequest

type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" form:"items"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" form:"product_id"`
	Quantity  int `json:"quantity" form:"quantity"`
}