package warehouse

import (
	"warehouse-service/models/warehouse"
)

type WarehouseRepository interface {
	Insert(warehouse *warehouse.RegisterRequest) error
}

type WarehouseUsecase struct {
	warehouseRepo WarehouseRepository
}

func NewWarehouseUsecase(warehouseRepo WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{
		warehouseRepo: warehouseRepo,
	}
}

func (w *WarehouseUsecase) Register(warehouseRegister *warehouse.RegisterRequest) error {
	return w.warehouseRepo.Insert(warehouseRegister)
}
