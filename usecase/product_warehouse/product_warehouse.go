package product_warehouse

import (
	"warehouse-service/models/product_warehouse"
)

type ProductWarehouseRepository interface {
	Insert(productWarehouse *product_warehouse.RegisterRequest) error
	AddAvailableStock(productId int, warehouseId int, addedAvailableStock int) error
	SubstractAvailableStock(productId int, warehouseId int, substractedAvailableStock int) error
	AddAvailableStockSubsReservedStock(productId int, warehouseId int, addedAvailableStock int, substractedReservedStock int) error
	SubsAvailableStockAddReservedStock(productId int, warehouseId int, substractedAvailableStock int, addedReservedStock int) error
	SubstractReservedStock(productId int, warehouseId int, substractedReservedStock int) error
}

type ProductWarehouseUsecase struct {
	productWarehouseRepo ProductWarehouseRepository
}

func NewProductWarehouseUsecase(productWarehouseRepo ProductWarehouseRepository) *ProductWarehouseUsecase {
	return &ProductWarehouseUsecase{
		productWarehouseRepo: productWarehouseRepo,
	}
}

func (p *ProductWarehouseUsecase) Register(warehouseRegister *product_warehouse.RegisterRequest) error {
	return p.productWarehouseRepo.Insert(warehouseRegister)
}
