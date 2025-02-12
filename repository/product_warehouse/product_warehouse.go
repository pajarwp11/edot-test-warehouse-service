package product_warehouse

import (
	"warehouse-service/entity"
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
	_, err := tx.Exec("UPDATE product_warehouses SET available_stock = available_stock + ? WHERE product_id=? and warehouse_id=?", addedAvailableStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubstractAvailableStock(tx *sqlx.Tx, productId int, warehouseId int, substractedAvailableStock int) error {
	_, err := tx.Exec("UPDATE product_warehouses SET available_stock = available_stock - ? WHERE product_id=? and warehouse_id=?", substractedAvailableStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) AddAvailableStockSubsReservedStock(tx *sqlx.Tx, productId int, warehouseId int, addedAvailableStock int, substractedReservedStock int) error {
	_, err := tx.Exec("UPDATE product_warehouses SET available_stock = available_stock + ?, reserved_stock = reserved_stock - ? WHERE product_id=? and warehouse_id=?", addedAvailableStock, substractedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubsAvailableStockAddReservedStock(tx *sqlx.Tx, productId int, warehouseId int, substractedAvailableStock int, addedReservedStock int) error {
	_, err := tx.Exec("UPDATE product_warehouses SET available_stock = available_stock - ?, reserved_stock = reserved_stock + ? WHERE product_id=? and warehouse_id=?", substractedAvailableStock, addedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) SubstractReservedStock(tx *sqlx.Tx, productId int, warehouseId int, substractedReservedStock int) error {
	_, err := p.mysql.Exec("UPDATE product_warehouses SET reserved_stock = reserved_stock - ? WHERE product_id=? and warehouse_id=?", substractedReservedStock, productId, warehouseId)
	return err
}

func (p *ProductWarehouseRepository) GetByProductAndWarehouseId(productId int, wareHouseId int) (*product_warehouse.ProductWarehouse, error) {
	data := product_warehouse.ProductWarehouse{}
	err := p.mysql.Get(&data, "SELECT id,product_id,warehouse_id,available_stock,reserved_stock FROM product_warehouses WHERE product_id=? and warehouse_id=?", productId, wareHouseId)
	return &data, err
}

func (p *ProductWarehouseRepository) GetAvailableStockBulk(availableStockRequest []product_warehouse.ProductShop) (map[int]int, error) {
	productIds := []int{}
	shopIds := []int{}
	for _, productShopMap := range availableStockRequest {
		productIds = append(productIds, productShopMap.ProductId)
		shopIds = append(shopIds, productShopMap.ShopId)
	}

	query, args, err := sqlx.In(`
		SELECT pw.product_id, w.shop_id, COALESCE(SUM(pw.available_stock), 0) AS total_stock
		FROM product_warehouses pw
		JOIN warehouses w ON pw.warehouse_id = w.id
		WHERE pw.product_id IN (?) AND w.shop_id IN (?) AND w.status = ?
		GROUP BY pw.product_id, w.shop_id
	`, productIds, shopIds, entity.WarehouseActive)
	if err != nil {
		return nil, err
	}

	query = p.mysql.Rebind(query)
	rows, err := p.mysql.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stockMap := make(map[int]int)
	for rows.Next() {
		var productId, shopId, totalStock int
		if err := rows.Scan(&productId, &shopId, &totalStock); err != nil {
			return nil, err
		}
		stockMap[productId] = totalStock
	}

	return stockMap, nil
}

func (p *ProductWarehouseRepository) GetAllByProductId(productId int) ([]product_warehouse.ProductWarehouse, error) {
	query := `
		SELECT pw.id, pw.product_id, pw.warehouse_id, pw.available_stock, pw.reserved_stock
		FROM product_warehouses pw
		JOIN warehouses w ON pw.warehouse_id = w.id
		WHERE pw.product_id = ? AND w.status = ?
		ORDER BY pw.id asc
	`

	var productWarehouses []product_warehouse.ProductWarehouse
	err := p.mysql.Select(&productWarehouses, query, productId, entity.WarehouseActive)
	if err != nil {
		return nil, err
	}
	return productWarehouses, nil
}

func (p *ProductWarehouseRepository) GetAvailableStock(productId int) (int, error) {
	query := `
		SELECT COALESCE(SUM(pw.available_stock), 0)
		FROM product_warehouses pw
		JOIN warehouses w ON pw.warehouse_id = w.id
		WHERE pw.product_id = ? AND w.status = ?
	`

	var availableStock int
	err := p.mysql.Get(&availableStock, query, productId, entity.WarehouseActive)
	if err != nil {
		return 0, err
	}
	return availableStock, nil
}

func (p *ProductWarehouseRepository) InsertOrderWarehouse(tx *sqlx.Tx, orderWarehouse *product_warehouse.OrderWarehouse) error {
	_, err := tx.Exec("INSERT INTO order_warehouses (order_id,product_id,warehouse_id,reserved_stock) VALUES (?,?,?,?)", orderWarehouse.OrderId, orderWarehouse.ProductId, orderWarehouse.WarehouseId, orderWarehouse.ReservedStock)
	return err
}

func (p *ProductWarehouseRepository) GetOrderWarehouseByOrderId(orderId int) ([]product_warehouse.OrderWarehouse, error) {
	query := `
		SELECT id, order_id, product_id, warehouse_id, reserved_stock
		FROM order_warehouses
		WHERE order_id = ?
	`

	var orderWarehouses []product_warehouse.OrderWarehouse
	err := p.mysql.Select(&orderWarehouses, query, orderId)
	if err != nil {
		return nil, err
	}

	return orderWarehouses, nil
}
