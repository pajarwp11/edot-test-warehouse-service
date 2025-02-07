package warehouse

import (
	"warehouse-service/models/warehouse"

	"github.com/jmoiron/sqlx"
)

type WarehouseRepository struct {
	mysql *sqlx.DB
}

func NewWarehouseRepository(mysql *sqlx.DB) *WarehouseRepository {
	return &WarehouseRepository{
		mysql: mysql,
	}
}

func (w *WarehouseRepository) Insert(warehouse *warehouse.RegisterRequest) error {
	_, err := w.mysql.Exec("INSERT INTO warehouses (name,address,user_id) VALUES (?,?,?)", warehouse.Name, warehouse.Address, warehouse.ShopId)
	return err
}
