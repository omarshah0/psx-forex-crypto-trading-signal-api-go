package models

import (
	"time"
)

type SignalType string
type SignalResult string
type AssetClass string
type DurationType string

const (
	SignalTypeLong  SignalType = "LONG"
	SignalTypeShort SignalType = "SHORT"
)

const (
	SignalResultWin       SignalResult = "WIN"
	SignalResultLoss      SignalResult = "LOSS"
	SignalResultBreakeven SignalResult = "BREAKEVEN"
)

const (
	AssetClassForex  AssetClass = "FOREX"
	AssetClassCrypto AssetClass = "CRYPTO"
	AssetClassPSX    AssetClass = "PSX"
)

const (
	DurationTypeShortTerm DurationType = "SHORT_TERM"
	DurationTypeLongTerm  DurationType = "LONG_TERM"
)

type TradingSignal struct {
	ID              int64         `json:"id" db:"id"`
	Symbol          string        `json:"symbol" db:"symbol" validate:"required"`
	AssetClass      AssetClass    `json:"asset_class" db:"asset_class" validate:"required,oneof=FOREX CRYPTO PSX"`
	DurationType    DurationType  `json:"duration_type" db:"duration_type" validate:"required,oneof=SHORT_TERM LONG_TERM"`
	StopLossPrice   float64       `json:"stop_loss_price" db:"stop_loss_price" validate:"required,gt=0"`
	EntryPrice      float64       `json:"entry_price" db:"entry_price" validate:"required,gt=0"`
	TakeProfitPrice float64       `json:"take_profit_price" db:"take_profit_price" validate:"required,gt=0"`
	Type            SignalType    `json:"type" db:"type" validate:"required,oneof=LONG SHORT"`
	Result          *SignalResult `json:"result" db:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return          *float64      `json:"return" db:"return"`
	FreeForAll      bool          `json:"free_for_all" db:"free_for_all"`
	Comments        *string       `json:"comments,omitempty" db:"comments"`
	CreatedBy       int64         `json:"created_by" db:"created_by"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

// TradingSignalCreate represents the data needed to create a new trading signal
type TradingSignalCreate struct {
	Symbol          string        `json:"symbol" validate:"required"`
	AssetClass      AssetClass    `json:"asset_class" validate:"required,oneof=FOREX CRYPTO PSX"`
	DurationType    DurationType  `json:"duration_type" validate:"required,oneof=SHORT_TERM LONG_TERM"`
	StopLossPrice   float64       `json:"stop_loss_price" validate:"required,gt=0"`
	EntryPrice      float64       `json:"entry_price" validate:"required,gt=0"`
	TakeProfitPrice float64       `json:"take_profit_price" validate:"required,gt=0"`
	Type            SignalType    `json:"type" validate:"required,oneof=LONG SHORT"`
	Result          *SignalResult `json:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return          *float64      `json:"return"`
	FreeForAll      bool          `json:"free_for_all"`
	Comments        *string       `json:"comments,omitempty"`
}

// TradingSignalUpdate represents the data needed to update a trading signal
type TradingSignalUpdate struct {
	Symbol          *string       `json:"symbol"`
	AssetClass      *AssetClass   `json:"asset_class" validate:"omitempty,oneof=FOREX CRYPTO PSX"`
	DurationType    *DurationType `json:"duration_type" validate:"omitempty,oneof=SHORT_TERM LONG_TERM"`
	StopLossPrice   *float64      `json:"stop_loss_price" validate:"omitempty,gt=0"`
	EntryPrice      *float64      `json:"entry_price" validate:"omitempty,gt=0"`
	TakeProfitPrice *float64      `json:"take_profit_price" validate:"omitempty,gt=0"`
	Type            *SignalType   `json:"type" validate:"omitempty,oneof=LONG SHORT"`
	Result          *SignalResult `json:"result" validate:"omitempty,oneof=WIN LOSS BREAKEVEN"`
	Return          *float64      `json:"return"`
	FreeForAll      *bool         `json:"free_for_all"`
	Comments        *string       `json:"comments,omitempty"`
}

