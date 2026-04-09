package invoice

import "encoding/json"

type BroadcastEnvelope struct {
	Event   string          `json:"event"`
	Data    json.RawMessage `json:"data"`
	Channel string          `json:"channel"`
}

type Content struct {
	PartnerName  string  `json:"partner_name"`
	CashierName  string  `json:"cashier_name"`
	InvoiceID    string  `json:"invoice_id"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
	Discount     float64 `json:"discount"`
	DiscountPerc float64 `json:"discount_perc"`
	Tax          float64 `json:"tax"`
	TaxPerc      float64 `json:"tax_perc"`
	GrandTotal   float64 `json:"grand_total"`
	PayedTotal   float64 `json:"payed_total"`
	ReturnTotal  float64 `json:"return_total"`
	OrderDate    string  `json:"order_date"`
}

type Item struct {
	Description string  `json:"description"`
	Qty         int     `json:"qty"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}
