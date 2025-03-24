package models

import (
	"fmt"
	"strings"
	"time"
)

type CurrencyQuotation struct {
	Code       string   `json:"code"`
	Codein     string   `json:"codein"`
	Name       string   `json:"name"`
	High       string   `json:"high"`
	Low        string   `json:"low"`
	VarBid     string   `json:"varBid"`
	PctChange  string   `json:"pctChange"`
	Bid        BidValue `json:"bid"`
	Ask        string   `json:"ask"`
	Timestamp  string   `json:"timestamp"`
	CreateDate string   `json:"create_date"`
}

type BidResponse struct {
	Bid BidValue `json:"bid"`
}

type QuotationRequest struct {
	FromCurrency Currency
	ToCurrency   Currency
	RequestData  time.Time
}

func (qr *QuotationRequest) GetCurrencyPair() string {
	return fmt.Sprintf("%s-%s", qr.FromCurrency.ISO4217Code, qr.ToCurrency.ISO4217Code)
}

func (qr *QuotationRequest) GetQuotationKey() string {
	return strings.ToUpper(qr.FromCurrency.ISO4217Code) + strings.ToUpper(qr.ToCurrency.ISO4217Code)
}
