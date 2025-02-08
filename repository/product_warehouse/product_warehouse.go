package product_warehouse

import (
	"warehouse-service/models/product_warehouse"

	"github.com/jmoiron/sqlx"
)

type ProductWarehouseRepository struct {
	mysql *sqlx.DB
}

func NewProductWarehouseRepository(mysql *sqlx.DB) *ProductWarehouseRepository {
	return &ProductWarehouseRepository{
		mysql: mysql,
	}
}

func (p *ProductWarehouseRepository) Insert(productWarehouse *product_warehouse.RegisterRequest) error {
	_, err := p.mysql.Exec("INSERT INTO product_warehouses (product_id,warehouse_id,available_stock) VALUES (?,?,?)", productWarehouse.ProductId, productWarehouse.WarehouseId, productWarehouse.AvailableStock)
	return err
}
