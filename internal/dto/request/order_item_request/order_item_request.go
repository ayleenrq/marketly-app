package orderitemrequest

type CreateOrderItemRequest struct {
	OrderID   int `json:"order_id" form:"order_id"`
	ProductID int `json:"product_id" form:"product_id"`
	Quantity  int `json:"quantity" form:"quantity"`
}