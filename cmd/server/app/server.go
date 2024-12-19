package app

import (
	"encoding/json"
	"errors"
	"github.com/Yessentemir256/crud/cmd/server/app/middleware"
	"github.com/Yessentemir256/crud/pkg/customers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
)

// Server представляет собой логический сервер нашего приложения.
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
}

func NewServer(mux *mux.Router, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init инициализирует сервер (регистрирует все Handler'ы)
func (s *Server) Init() {
	// Применяем middleware для проверки заголовка Content-Type
	chMd := middleware.CheckHeader("Content-Type", "application/json")
	s.mux.Handle("/customers", chMd(http.HandlerFunc(s.handleSaveCustomer))).Methods(POST) // убрали старый хендлер и добавили с middleware

	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveCustomerByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActive).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)

	s.mux.Use(middleware.Logger) // использование middleware

	//s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	//s.mux.HandleFunc("/customers.save", s.handleSaveCustomer) // Новый обработчик
	//s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
}

func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	item, err := s.customerSvc.ByID(request.Context(), id)
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
	if err != nil {
		log.Print(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleGetAllCustomers(writer http.ResponseWriter, request *http.Request) {
	// Вызов бизнес-логики
	customers, err := s.customerSvc.GetAll(request.Context())
	if err != nil {
		http.Error(writer, "Failed to get customers", http.StatusInternalServerError)
		return
	}

	// Преобразование данных в JSON
	data, err := json.Marshal(customers)
	if err != nil {
		http.Error(writer, "Failed to marshal customers", http.StatusInternalServerError)
		return
	}

	// Отправка данных в ответ
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err = s.customerSvc.Save(request.Context(), item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleRemoveCustomerByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	err = s.customerSvc.RemoveByID(request.Context(), int(id))
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Успешный ответ, без содержимого, только статус
	writer.WriteHeader(http.StatusNoContent)

}

func (s *Server) handleGetAllActive(writer http.ResponseWriter, request *http.Request) {
	// Вызов бизнес-логики
	customers, err := s.customerSvc.GetAllActive(request.Context())
	if err != nil {
		http.Error(writer, "Failed to get customers", http.StatusInternalServerError)
		return
	}

	// Преобразование данных в JSON
	data, err := json.Marshal(customers)
	if err != nil {
		http.Error(writer, "Failed to marshal customers", http.StatusInternalServerError)
		return
	}

	// Отправка данных в ответ
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
}

func (s *Server) handleBlockByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	err = s.customerSvc.UnBlockByID(request.Context(), int64(id))
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Успешный ответ, без содержимого, только статус
	writer.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUnBlockByID(writer http.ResponseWriter, request *http.Request) {
	//idParam := request.URL.Query().Get("id")
	idParam, ok := mux.Vars(request)["id"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
	}

	err = s.customerSvc.UnBlockByID(request.Context(), int64(id))
	if errors.Is(err, customers.ErrNotFound) {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Успешный ответ, без содержимого, только статус
	writer.WriteHeader(http.StatusNoContent)

}
