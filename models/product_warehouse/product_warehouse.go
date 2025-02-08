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

type TransferStockRequest struct {
	ProductId       int `json:"product_id" validate:"required"`
	FromWarehouseId int `json:"from_warehouse_id" validate:"required"`
	ToWarehouseId   int `json:"to_warehouse_id" validate:"required"`
	Quantity        int `json:"1uantity" validate:"required"`
}
