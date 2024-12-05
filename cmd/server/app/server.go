package app

import (
	"encoding/json"
	"errors"
	"github.com/Yessentemir256/crud/pkg/customers"
	"log"
	"net/http"
	"strconv"
)

// Server представляет собой логический сервер нашего приложения.
type Server struct {
	mux         *http.ServeMux
	customerSvc *customers.Service
}

func NewServer(mux *http.ServeMux, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

// Init инициализирует сервер (регистрирует все Handler'ы)
func (s *Server) Init() {
	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/customers.save", s.handleSaveCustomer) // Новый обработчик
	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)

}

func (s *Server) handleGetCustomerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

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

func (s *Server) handleSaveCustomer(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := request.FormValue("name")
	phone := request.FormValue("phone")
	idParam := request.FormValue("id")

	var id int
	if idParam != "" {
		var err error
		id, err = strconv.Atoi(idParam)
		if err != nil {
			http.Error(writer, "Invalid ID", http.StatusBadRequest)
			return
		}
	}

	err := s.customerSvc.Save(request.Context(), id, name, phone)
	if err != nil {
		log.Print(err)
		http.Error(writer, "Failed to save customer", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Customer saved successfully"))
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
