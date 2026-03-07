package paymentrequest

type CreatePaymentRequest struct {
	OrderID int    `json:"order_id" form:"order_id"`
	Method  string `json:"method" form:"method"`
}
