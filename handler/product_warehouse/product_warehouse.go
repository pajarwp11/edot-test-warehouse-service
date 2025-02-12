package product_warehouse

import (
	"encoding/json"
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
	ReleaseReservedStock(order *product_warehouse.Order) error
	ReturnReservedStock(order *product_warehouse.Order) error
	GetAvailableStockBulk(getAvailableStock []product_warehouse.ProductShop) (map[int]int, error)
	ReserveStock(operationStock *product_warehouse.StockOperationOrderRequest) error
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

	if request.FromWarehouseId == request.ToWarehouseId {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = "transfer must be to different warehouse"
		json.NewEncoder(w).Encode(response)
		return
	}

	err := p.productWarehouseUsecase.TransferStockRequest(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response.Message = "transfer in progress"
	json.NewEncoder(w).Encode(response)
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
	w.WriteHeader(http.StatusOK)
	response.Message = "add stock in progress"
	json.NewEncoder(w).Encode(response)
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
	w.WriteHeader(http.StatusOK)
	response.Message = "deduct stock in progress"
	json.NewEncoder(w).Encode(response)
}

func (p *ProductWarehouseHandler) GetAvailableStock(w http.ResponseWriter, req *http.Request) {
	var request []product_warehouse.ProductShop
	response := Response{}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "invalid request body"
		json.NewEncoder(w).Encode(response)
		return
	}

	for _, item := range request {
		if err := validate.Struct(item); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
			return
		}
	}

	availableStock, err := p.productWarehouseUsecase.GetAvailableStockBulk(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Message = "get available stock success"
	response.Data = availableStock
	json.NewEncoder(w).Encode(response)
}
