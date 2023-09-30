package tracker_test

import (
	"testing"

	"github.com/tessellated-io/mail-in-rebates/paymaster/tracker"
)

func TestFileWrite(t *testing.T) {
	fileName := "addresses.csv"
	a := tracker.NewAddressTracker(fileName)
	a.AddAddress("cosmos123")
	a.AddAddress("osmosis123")

	a.AddAddress("cosmos123")

	a2 := tracker.NewAddressTracker(fileName)
	a2.AddAddress("axelar456")
}
