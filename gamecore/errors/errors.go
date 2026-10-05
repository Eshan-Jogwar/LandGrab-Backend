package errors

import stderrors "errors"

var (
	ErrNotPurchasable   = stderrors.New("tile is not purchasable")
	ErrNotRentable      = stderrors.New("tile is not rentable")
	ErrNotBuildable     = stderrors.New("tile is not buildable")
	ErrNoMortgage       = stderrors.New("tile cannot be mortgaged")
	ErrInvalidExtraInfo = stderrors.New("invalid extra info")
	ErrTileNotFound    = stderrors.New("tile not found")
)
