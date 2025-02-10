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
	AddAvailableStockSubsReservedStock(tx *sqlx.Tx, productId int, warehouseId int, addedAvailableStock int, substractedReservedStock int) error
	SubsAvailableStockAddReservedStock(tx *sqlx.Tx, productId int, warehouseId int, substractedAvailableStock int, addedReservedStock int) error
	SubstractReservedStock(productId int, warehouseId int, substractedReservedStock int) error
	GetByProductAndWarehouseId(productId int, wareHouseId int) (*product_warehouse.ProductWarehouse, error)
	GetAvailableStockBulk(availableStockRequest []product_warehouse.ProductShop) (map[int]int, error)
	GetAllByProductId(productId int) ([]product_warehouse.ProductWarehouse, error)
	GetAvailableStock(productId int) (int, error)
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
		// TO DO: send notif to user
		return errors.New(entity.ErrorInsufficientStock)
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
	sourceProduct, err := p.productWarehouseRepo.GetByProductAndWarehouseId(deductStock.ProductId, deductStock.WarehouseId)
	if err != nil {
		return err
	}

	if sourceProduct.AvailableStock < deductStock.Quantity {
		// TO DO: send notif to user
		return errors.New(entity.ErrorInsufficientStock)
	}

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

func (p *ProductWarehouseUsecase) ReturnReservedStock(operationStock []product_warehouse.StockOperationProductRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, operation := range operationStock {
		productWarehouses, err := p.productWarehouseRepo.GetAllByProductId(operation.ProductId)
		if err != nil {
			tx.Rollback()
			return err
		}

		returnedStock := operation.Quantity

		for i := 0; i < len(productWarehouses) && returnedStock > 0; i++ {
			warehouse := productWarehouses[i]

			queryQuantity := min(warehouse.ReservedStock, returnedStock)

			err = p.productWarehouseRepo.AddAvailableStockSubsReservedStock(tx, operation.ProductId, warehouse.WarehouseId, queryQuantity, queryQuantity)
			if err != nil {
				tx.Rollback()
				return err
			}

			returnedStock -= queryQuantity
		}
	}

	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) ReserveStock(operationStock *product_warehouse.StockOperationOrderRequest) error {
	tx, err := p.mysql.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, operation := range operationStock.StockOperations {
		availableStock, err := p.productWarehouseRepo.GetAvailableStock(operation.ProductId)
		if err != nil {
			tx.Rollback()
			return err
		}
		if availableStock < operation.Quantity {
			tx.Rollback()
			updateOrderRequest := product_warehouse.UpdateStatusRequest{
				Id:     operationStock.OrderId,
				Status: "cancel",
			}
			go p.publisher.PublishEvent(entity.OrderUpdateStatusEvent, updateOrderRequest)
			return errors.New(entity.ErrorInsufficientStock)
		}
		productWarehouses, err := p.productWarehouseRepo.GetAllByProductId(operation.ProductId)
		if err != nil {
			tx.Rollback()
			return err
		}

		reservedStock := operation.Quantity

		for i := 0; i < len(productWarehouses) && reservedStock > 0; i++ {
			warehouse := productWarehouses[i]

			queryQuantity := min(warehouse.AvailableStock, reservedStock)

			err = p.productWarehouseRepo.SubsAvailableStockAddReservedStock(tx, operation.ProductId, warehouse.WarehouseId, queryQuantity, queryQuantity)
			if err != nil {
				tx.Rollback()
				return err
			}

			reservedStock -= queryQuantity
		}
	}

	tx.Commit()
	return nil
}

func (p *ProductWarehouseUsecase) GetAvailableStockBulk(getAvailableStock []product_warehouse.ProductShop) (map[int]int, error) {
	return p.productWarehouseRepo.GetAvailableStockBulk(getAvailableStock)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
