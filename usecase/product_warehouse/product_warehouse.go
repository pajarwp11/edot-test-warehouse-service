package product_warehouse

import (
	"warehouse-service/entity"
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

type Publisher interface {
	PublishEvent(eventType string, data interface{}) error
}

type ProductWarehouseUsecase struct {
	productWarehouseRepo ProductWarehouseRepository
	publisher            Publisher
}

func NewProductWarehouseUsecase(productWarehouseRepo ProductWarehouseRepository, publisher Publisher) *ProductWarehouseUsecase {
	return &ProductWarehouseUsecase{
		productWarehouseRepo: productWarehouseRepo,
		publisher:            publisher,
	}
}

func (p *ProductWarehouseUsecase) Register(productWarehouseRegister *product_warehouse.RegisterRequest) error {
	return p.productWarehouseRepo.Insert(productWarehouseRegister)
}

func (p *ProductWarehouseUsecase) TransferStock(transferStock *product_warehouse.TransferStockRequest) error {
	return p.publisher.PublishEvent(entity.StockTransferEvent, transferStock)
}
