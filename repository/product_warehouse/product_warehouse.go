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

func (p *ProductWarehouseRepository) AddAvailableStock(tx *sqlx.Tx, productId int, warehouseId int, addedAvailableStock int) error {
	_, err := tx.Exec("UPDATE warehouses SET available_stock = available_stock + ? WHERE product_id=? and warehouse_id=?", addedAvailableStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubstractAvailableStock(tx *sqlx.Tx, productId int, warehouseId int, substractedAvailableStock int) error {
	_, err := tx.Exec("UPDATE warehouses SET available_stock = available_stock - ? WHERE product_id=? and warehouse_id=?", substractedAvailableStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) AddAvailableStockSubsReservedStock(productId int, warehouseId int, addedAvailableStock int, substractedReservedStock int) error {
	_, err := p.mysql.Exec("UPDATE warehouses SET available_stock = available_stock + ?, reserved_stock = reserved_stock - ? WHERE product_id=? and warehouse_id=?", addedAvailableStock, substractedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubsAvailableStockAddReservedStock(productId int, warehouseId int, substractedAvailableStock int, addedReservedStock int) error {
	_, err := p.mysql.Exec("UPDATE warehouses SET available_stock = available_stock - ?, reserved_stock = reserved_stock + ? WHERE product_id=? and warehouse_id=?", substractedAvailableStock, addedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubstractReservedStock(productId int, warehouseId int, substractedReservedStock int) error {
	_, err := p.mysql.Exec("UPDATE warehouses SET reserved_stock = reserved_stock - ? WHERE product_id=? and warehouse_id=?", substractedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) GetByProductAndWarehouseId(productId int, wareHouseId int) (*product_warehouse.ProductWarehouse, error) {
	data := product_warehouse.ProductWarehouse{}
	err := p.mysql.Get(&data, "SELECT id,product_id,warehouse_id,available_stock,reserved_stock FROM shops WHERE product_id=? and warehouse_id=?", productId, wareHouseId)
	return &data, err
}
