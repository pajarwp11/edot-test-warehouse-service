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
	Quantity        int `json:"quantity" validate:"required"`
}

type StockOperationRequest struct {
	ProductId   int `json:"product_id" validate:"required"`
	WarehouseId int `json:"warehouse_id" validate:"required"`
	Quantity    int `json:"quantity" validate:"required"`
}

type StockOperationOrderRequest struct {
	OrderId         int                     `json:"order_id"`
	StockOperations []StockOperationRequest `json:"stock_operations" validate:"required"`
}

type ProductShop struct {
	ProductId int `json:"product_id" validate:"required"`
	ShopId    int `json:"shop_id" validate:"required"`
}

type UpdateStatusRequest struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type StockOperationProductRequest struct {
	ProductId int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required"`
}

type OrderWarehouse struct {
	Id            int `db:"id"`
	OrderId       int `db:"order_id"`
	ProductId     int `db:"product_id"`
	WarehouseId   int `db:"warehouse_id"`
	ReservedStock int `db:"reserved_stock"`
}
