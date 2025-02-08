package warehouse

type Warehouse struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	ShopId  int    `db:"shop_id"`
}

type RegisterRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	ShopId  int    `json:"shop_id" validate:"required"`
}

type UpdateStatusRequest struct {
	Id     int    `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}
