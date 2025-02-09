package product_warehouse

import (
	"errors"
	"warehouse-service/entity"
	"warehouse-service/models/product_warehouse"

	"github.com/jmoiron/sqlx"
)

type ProductWarehouseRepository interface {
	Insert(productWarehouse *product_warehouse.RegisterRequest) error
	AddAvailableStock(tx *sqlx.Tx, productId int, warehouseId int, addedAvailableStock int) error
	SubstractAvailableStock(tx *sqlx.Tx, productId int, warehouseId int, substractedAvailableStock int) error
	AddAvailableStockSubsReservedStock(productId int, warehouseId int, addedAvailableStock int, substractedReservedStock int) error
	SubsAvailableStockAddReservedStock(productId int, warehouseId int, substractedAvailableStock int, addedReservedStock int) error
	SubstractReservedStock(productId int, warehouseId int, substractedReservedStock int) error
	GetByProductAndWarehouseId(productId int, wareHouseId int) (*product_warehouse.ProductWarehouse, error)
	GetAvailableStockBulk(availableStockRequest *product_warehouse.GetAvailableStockRequest) (map[int]int, error)
}

type Publisher interface {
	PublishEvent(eventType string, data interface{}) error
}

type ProductWarehouseUsecase struct {
	productWarehouseRepo ProductWarehouseRepository
	publisher            Publisher
	mysql                *sqlx.DB
}

func NewProductWarehouseUsecase(productWarehouseRepo ProductWarehouseRepository, publisher Publisher, mysql *sqlx.DB) *ProductWarehouseUsecase {
	return &ProductWarehouseUsecase{
		productWarehouseRepo: productWarehouseRepo,
		publisher:            publisher,
		mysql:                mysql,
	}
}

func (p *ProductWarehouseUsecase) Register(productWarehouseRegister *product_warehouse.RegisterRequest) error {
	return p.productWarehouseRepo.Insert(productWarehouseRegister)
}

func (p *ProductWarehouseUsecase) TransferStockRequest(transferStock *product_warehouse.TransferStockRequest) error {
	return p.publisher.PublishEvent(entity.StockTransferEvent, transferStock)
}

func (p *ProductWarehouseUsecase) TransferStock(transferStock *product_warehouse.TransferStockRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	sourceProduct, err := p.productWarehouseRepo.GetByProductAndWarehouseId(transferStock.ProductId, transferStock.FromWarehouseId)
	if err != nil {
		return err
	}

	if sourceProduct.AvailableStock < transferStock.Quantity {
		return errors.New("stock is less than quantity")
	}

	err = p.productWarehouseRepo.AddAvailableStock(tx, transferStock.ProductId, transferStock.ToWarehouseId, transferStock.Quantity)
	if err != nil {
		return err
	}

	err = p.productWarehouseRepo.SubstractAvailableStock(tx, transferStock.ProductId, transferStock.FromWarehouseId, transferStock.Quantity)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) AddStockRequest(addStock *product_warehouse.StockOperationRequest) error {
	return p.publisher.PublishEvent(entity.StockAddEvent, addStock)
}

func (p *ProductWarehouseUsecase) AddStock(addStock *product_warehouse.StockOperationRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = p.productWarehouseRepo.AddAvailableStock(tx, addStock.ProductId, addStock.WarehouseId, addStock.Quantity)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) DeductStockRequest(deductStock *product_warehouse.StockOperationRequest) error {
	return p.publisher.PublishEvent(entity.StockDeductEvent, deductStock)
}

func (p *ProductWarehouseUsecase) DeductStock(deductStock *product_warehouse.StockOperationRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = p.productWarehouseRepo.SubstractAvailableStock(tx, deductStock.ProductId, deductStock.WarehouseId, deductStock.Quantity)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) ReleaseReservedStock(operationStock *product_warehouse.StockOperationRequest) error {
	err := p.productWarehouseRepo.SubstractReservedStock(operationStock.ProductId, operationStock.WarehouseId, operationStock.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseUsecase) ReturnReservedStock(operationStock *product_warehouse.StockOperationRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = p.productWarehouseRepo.AddAvailableStockSubsReservedStock(operationStock.ProductId, operationStock.WarehouseId, operationStock.Quantity, operationStock.Quantity)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) GetAvailableStockBulk(getAvailableStock *product_warehouse.GetAvailableStockRequest) (map[int]int, error) {
	return p.productWarehouseRepo.GetAvailableStockBulk(getAvailableStock)
}
