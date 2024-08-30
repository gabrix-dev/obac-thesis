package models

import (
	"errors"
)

type BookingInfo struct {
	Seat int    `json:"seat"`
	Zone string `json:"zone"`
}

func (b BookingInfo) ValidateBooking(ticketType string, tokenId float64) error {
	switch ticketType {
	case "Basic":
		if b.Zone != "Zone 2" && b.Zone != "Zone 3" {
			return errors.New("Basic ticket only gives you access to Zone 2 and Zone 3")
		}
	case "Basic plus":
		if b.Zone != "Zone 1" {
			return errors.New("Basic  plus ticket only gives you access to Zone 1")
		}
	case "VIP premium":
		if b.Zone != "VIP Zone" {
			return errors.New("VIP Premium ticket only gives you access to the VIP Zone")

		}
	default:
		return errors.New("Invalid ticket type")
	}
	return nil
}
