package product_warehouse

import (
	"errors"
	"warehouse-service/models/product_warehouse"
)

func (p *ProductWarehouseHandler) TransferStock(data interface{}) error {
	request, ok := data.(product_warehouse.TransferStockRequest)
	if !ok {
		return errors.New("invalid body request")
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err := p.productWarehouseUsecase.TransferStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) AddStock(data interface{}) error {
	request, ok := data.(product_warehouse.StockOperationRequest)
	if !ok {
		return errors.New("invalid body request")
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err := p.productWarehouseUsecase.AddStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) DeductStock(data interface{}) error {
	request, ok := data.(product_warehouse.StockOperationRequest)
	if !ok {
		return errors.New("invalid body request")
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err := p.productWarehouseUsecase.DeductStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) ReleaseReservedStock(data interface{}) error {
	request, ok := data.(product_warehouse.StockOperationRequest)
	if !ok {
		return errors.New("invalid body request")
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err := p.productWarehouseUsecase.ReleaseReservedStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) ReturnReservedStock(data interface{}) error {
	request, ok := data.(product_warehouse.StockOperationRequest)
	if !ok {
		return errors.New("invalid body request")
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err := p.productWarehouseUsecase.ReturnReservedStock(&request)
	if err != nil {
		return err
	}
	return nil
}
