package product_warehouse

import (
	"encoding/json"
	"net/http"
	"warehouse-service/models/product_warehouse"

	"github.com/go-playground/validator/v10"
)

type ProductWarehouseUsecase interface {
	Register(productWarehouseRegister *product_warehouse.RegisterRequest) error
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
