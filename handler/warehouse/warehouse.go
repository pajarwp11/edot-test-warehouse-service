package warehouse

import (
	"encoding/json"
	"net/http"
	"strconv"
	"warehouse-service/models/warehouse"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type WarehouseUsecase interface {
	Register(warehouseRegister *warehouse.RegisterRequest) error
	UpdateStatus(updateStatus *warehouse.UpdateStatusRequest) error
}

type WarehouseHandler struct {
	warehouseUsecase WarehouseUsecase
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var validate = validator.New()

func NewWarehouseHandler(warehouseUsecase WarehouseUsecase) *WarehouseHandler {
	return &WarehouseHandler{
		warehouseUsecase: warehouseUsecase,
	}
}

func (wa *WarehouseHandler) Register(w http.ResponseWriter, req *http.Request) {
	request := warehouse.RegisterRequest{}
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

	err := wa.warehouseUsecase.Register(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "warehouse registered"
	json.NewEncoder(w).Encode(response)
}

func (wa *WarehouseHandler) UpdateStatus(w http.ResponseWriter, req *http.Request) {
	request := warehouse.UpdateStatusRequest{}
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

	vars := mux.Vars(req)
	id := vars["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "id is required"
		json.NewEncoder(w).Encode(response)
		return
	}
	var err error
	request.Id, err = strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "id must be numeric"
		json.NewEncoder(w).Encode(response)
		return
	}

	err = wa.warehouseUsecase.UpdateStatus(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.Message = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response.Message = "warehouse status updated"
	json.NewEncoder(w).Encode(response)
}
