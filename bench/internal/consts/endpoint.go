package consts

import "fmt"

const (
	// GET
	SearchTrainsPath     = "/train/search"
	ListTrainSeatsPath   = "/train/seats"
	ListReservationsPath = "/reservation"

	// POST
	RegisterPath = "/register"
	LoginPath    = "/login"
	ReservePath  = "/reserve"
)

var (
	// POST
	BuildCommitReservationPath = func(id int) string {
		return fmt.Sprintf("/reservation/%d/commit", id)
	}

	// DELETE
	BuildCancelReservationPath = func(id int) string {
		return fmt.Sprintf("/reservation/%d/cancel", id)
	}
)

var (
	// Mock
	MockCommitReservationPath = `=~^/reservation/(\d+)/commit\z`
	MockCancelReservationPath = `=~^/reservation/(\d+)/cancel\z`
)
