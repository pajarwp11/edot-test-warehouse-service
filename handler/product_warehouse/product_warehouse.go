package product_warehouse

import (
	"encoding/json"
	"errors"
	"net/http"
	"warehouse-service/models/product_warehouse"

	"github.com/go-playground/validator/v10"
)

type ProductWarehouseUsecase interface {
	Register(productWarehouseRegister *product_warehouse.RegisterRequest) error
	TransferStockRequest(transferStock *product_warehouse.TransferStockRequest) error
	TransferStock(transferStock *product_warehouse.TransferStockRequest) error
	AddStockRequest(addStock *product_warehouse.StockOperationRequest) error
	AddStock(addStock *product_warehouse.StockOperationRequest) error
	DeductStockRequest(deductStock *product_warehouse.StockOperationRequest) error
	DeductStock(deductStock *product_warehouse.StockOperationRequest) error
	ReleaseReservedStock(operationStock *product_warehouse.StockOperationRequest) error
	ReturnReservedStock(operationStock *product_warehouse.StockOperationRequest) error
	GetAvailableStock(getAvailableStock *product_warehouse.GetAvailableStockRequest) (int, error)
}

type ProductWarehouseHandler struct {
	productWarehouseUsecase ProductWarehouseUsecase
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var validate = validator.New()

func NewProductWarehouseHandler(productWarehouseUsecase ProductWarehouseUsecase) *ProductWarehouseHandler {
	return &ProductWarehouseHandler{
		productWarehouseUsecase: productWarehouseUsecase,
	}
}

func (p *ProductWarehouseHandler) Register(w http.ResponseWriter, req *http.Request) {
	request := product_warehouse.RegisterRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := p.productWarehouseUsecase.Register(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "product warehouse registered"
	json.NewEncoder(w).Encode(response)
}

func (p *ProductWarehouseHandler) TranserStockRequest(w http.ResponseWriter, req *http.Request) {
	request := product_warehouse.TransferStockRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := p.productWarehouseUsecase.TransferStockRequest(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "stock is transferred"
	json.NewEncoder(w).Encode(response)
}

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

func (p *ProductWarehouseHandler) AddStockRequest(w http.ResponseWriter, req *http.Request) {
	request := product_warehouse.StockOperationRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := p.productWarehouseUsecase.AddStockRequest(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "stock is added"
	json.NewEncoder(w).Encode(response)
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

func (p *ProductWarehouseHandler) DeductStockRequest(w http.ResponseWriter, req *http.Request) {
	request := product_warehouse.StockOperationRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	err := p.productWarehouseUsecase.DeductStockRequest(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "stock is deducted"
	json.NewEncoder(w).Encode(response)
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

func (p *ProductWarehouseHandler) GetAvailableStock(w http.ResponseWriter, req *http.Request) {
	request := product_warehouse.GetAvailableStockRequest{}
	response := Response{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := validate.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	availableStock, err := p.productWarehouseUsecase.GetAvailableStock(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "get available stock success"
	response.Data = availableStock
	json.NewEncoder(w).Encode(response)
}
