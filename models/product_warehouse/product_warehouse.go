package product_warehouse

type ProductWarehouse struct {
	Id             int `db:"id"`
	ProductId      int `db:"product_id"`
	WarehouseId    int `db:"warehouse_id"`
	AvailableStock int `db:"available_stock"`
	ReservedStock  int `db:"reserved_stock"`
}

type RegisterRequest struct {
	ProductId      int `json:"product_id" validate:"required"`
	WarehouseId    int `json:"warehouse_id" validate:"required"`
	AvailableStock int `json:"available_stock"`
}
