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
