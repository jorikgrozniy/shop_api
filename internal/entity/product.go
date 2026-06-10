package entity

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID             uuid.UUID  `db:"id"`
	Name           string     `db:"name"`
	Category       string     `db:"category"`
	Price          float64    `db:"price"`
	AvailableStock int        `db:"available_stock"`
	LastUpdate     time.Time  `db:"last_update_date"`
	SupplierID     *uuid.UUID `db:"supplier_id"`
	ImageID        *uuid.UUID `db:"image_id"`
}
