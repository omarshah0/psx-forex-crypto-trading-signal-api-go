package models

import (
	"time"
)

type SignalType string
type SignalResult string

const (
	SignalTypeLong  SignalType = "LONG"
	SignalTypeShort SignalType = "SHORT"
)

const (
	SignalResultWin      SignalResult = "WIN"
	SignalResultLoss     SignalResult = "LOSS"
	SignalResultBreakeven SignalResult = "BREAKEVEN"
)

type TradingSignal struct {
	ID             int64        `json:"id" db:"id"`
	Symbol         string       `json:"symbol" db:"symbol" validate:"required"`
	StopLossPrice  float64      `json:"stop_loss_price" db:"stop_loss_price" validate:"required,gt=0"`
	EntryPrice     float64      `json:"entry_price" db:"entry_price" validate:"required,gt=0"`
	TakeProfitPrice float64     `json:"take_profit_price" db:"take_profit_price" validate:"required,gt=0"`
	Type           SignalType   `json:"type" db:"type" validate:"required,oneof=LONG SHORT"`
	Result         *SignalResult `json:"result" db:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return         *float64     `json:"return" db:"return"`
	CreatedBy      int64        `json:"created_by" db:"created_by"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

// TradingSignalCreate represents the data needed to create a new trading signal
type TradingSignalCreate struct {
	Symbol          string       `json:"symbol" validate:"required"`
	StopLossPrice   float64      `json:"stop_loss_price" validate:"required,gt=0"`
	EntryPrice      float64      `json:"entry_price" validate:"required,gt=0"`
	TakeProfitPrice float64      `json:"take_profit_price" validate:"required,gt=0"`
	Type            SignalType   `json:"type" validate:"required,oneof=LONG SHORT"`
	Result          *SignalResult `json:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return          *float64     `json:"return"`
}

// TradingSignalUpdate represents the data needed to update a trading signal
type TradingSignalUpdate struct {
	Symbol          *string       `json:"symbol"`
	StopLossPrice   *float64      `json:"stop_loss_price" validate:"omitempty,gt=0"`
	EntryPrice      *float64      `json:"entry_price" validate:"omitempty,gt=0"`
	TakeProfitPrice *float64      `json:"take_profit_price" validate:"omitempty,gt=0"`
	Type            *SignalType   `json:"type" validate:"omitempty,oneof=LONG SHORT"`
	Result          *SignalResult `json:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return          *float64      `json:"return"`
}

