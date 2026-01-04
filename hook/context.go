package hook

type OrderContext struct {
	OrderID string
	UserID  string
	Amount  int64

	Metadata map[string]any
}
