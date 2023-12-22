package View

type AddressView struct {
	ID            uint    `json:"id"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	AddressDetail string  `json:"address_detail"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	PosterID      uint    `json:"poster_id"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}
