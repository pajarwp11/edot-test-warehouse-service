package warehouse

import (
	"warehouse-service/models/warehouse"
)

type WarehouseRepository interface {
	Insert(warehouse *warehouse.RegisterRequest) error
	UpdateStatus(id int, status string) error
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

func (w *WarehouseUsecase) UpdateStatus(updateStatus *warehouse.UpdateStatusRequest) error {
	return w.warehouseRepo.UpdateStatus(updateStatus.Id, updateStatus.Status)
}
