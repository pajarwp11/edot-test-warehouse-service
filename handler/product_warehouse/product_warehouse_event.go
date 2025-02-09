package product_warehouse

import (
	"encoding/json"
	"errors"
	"warehouse-service/models/product_warehouse"
)

func (p *ProductWarehouseHandler) TransferStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := product_warehouse.TransferStockRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.TransferStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) AddStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := product_warehouse.StockOperationRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.AddStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) DeductStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := product_warehouse.StockOperationRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.DeductStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) ReleaseReservedStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := product_warehouse.StockOperationRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.ReleaseReservedStock(&request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) ReturnReservedStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := []product_warehouse.StockOperationRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.ReturnReservedStock(request)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductWarehouseHandler) ReserveStock(data interface{}) error {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return errors.New("invalid body request")
	}
	request := product_warehouse.StockOperationOrderRequest{}
	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		return err
	}
	if err := validate.Struct(request); err != nil {
		return err
	}

	err = p.productWarehouseUsecase.ReserveStock(&request)
	if err != nil {
		return err
	}
	return nil
}
