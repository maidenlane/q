package orderbook

import (
	"fmt"
)

type Match struct {
	TradeID      int
	Sequence     int64
	MakerOrderID string
	TakerOrderID string
	Time         fcxClient.Time
	ProductID    string
	Size         decimal.Decimal
	Price        decimal.Decimal
	MatchSide    common.Side
}

type Change struct {
	Time      fcxClient.Time
	Sequence  int64
	OrderID   string
	ProductID string
	NewSize   decimal.Decimal
	OldSize   decimal.Decimal
	Price     decimal.Decimal
	Side      common.Side
}

func NewMatchMessage(tradeID int, sequence int64, makerOrderID string, takerOrderID string,
	time fcxClient.Time, productID string, size string, price string, side string) (Match, error) {
	orderSide, err := ToSide(side)
	if err != nil {
		return Match{}, err
	}
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		return Match{}, errors.Wrap(err, "Alert: Unable to convert Price to decimal")
	}

	sizeDec, err := decimal.NewFromString(size)
	if err != nil {
		return Match{}, errors.Wrap(err, "Alert: Unable to convert size to decimal")
	}

	return Match{
		TradeID:      tradeID,
		Sequence:     sequence,
		MakerOrderID: makerOrderID,
		TakerOrderID: takerOrderID,
		Time:         time,
		ProductID:    productID,
		Size:         sizeDec,
		Price:        priceDec,
		MatchSide:    orderSide,
	}, nil
}

func (msg Match) ToString() string {
	return fmt.Sprintf(`TradeID: %d, Sequence %v, MakerOrderID: %s, TakerOrderID: %s, Time: %v, ProductID: %v,
                       Size: %s, Price: %s, MatchSide: %s`, msg.TradeID, msg.Sequence, msg.MakerOrderID, msg.TakerOrderID,
		msg.Time, msg.ProductID, msg.Size.String(), msg.Price.String(), common.ToString(msg.MatchSide))
}

func NewChangeMessage(time fcxClient.Time, sequence int64, orderID string, productID string, newSize string, oldSize string, price string, side string) (Change, error) {
	orderSide, err := ToSide(side)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return Change{}, err
	}
	var priceDec decimal.Decimal
	if price == "" {
		priceDec = decimal.NewFromFloat(0.0)
	} else {
		priceDec, err = decimal.NewFromString(price)
		if err != nil {
			return Change{}, errors.Wrap(err, "Alert: Unable to convert Price to decimal")
		}
	}

	newSizeDec, err := decimal.NewFromString(newSize)
	if err != nil {
		return Change{}, errors.Wrap(err, "Alert: Unable to convert newSize to decimal")
	}

	oldSizeDec, err := decimal.NewFromString(oldSize)
	if err != nil {
		return Change{}, errors.Wrap(err, "Alert: Unable to convert oldSize to decimal")
	}

	return Change{
		Time:      time,
		Sequence:  sequence,
		OrderID:   orderID,
		ProductID: productID,
		NewSize:   newSizeDec,
		OldSize:   oldSizeDec,
		Price:     priceDec,
		Side:      orderSide,
	}, nil
}