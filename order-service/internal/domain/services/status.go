package services

const (
	OrderPending    = "ORDER_PENDING"
	OrderPaid       = "ORDER_PAID"
	OrderPreparing  = "ORDER_PREPARING"
	OrderDelivering = "ORDER_DELIVERING"
	OrderCompleted  = "ORDER_COMPLETED"
	OrderCancelled  = "ORDER_CANCELLED"
)

const (
	PaymentPending   = "PAYMENT_PENDING"
	PaymentPaid      = "PAYMENT_PAID"
	PaymentFailed    = "PAYMENT_FAILED"
	PaymentCancelled = "PAYMENT_CANCELLED"
)
