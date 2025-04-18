package server

import (
	"fmt"
	"net/http"

	"PRACTICAS-GO-WEB/internal/handlers"
	"PRACTICAS-GO-WEB/internal/repository"
	"PRACTICAS-GO-WEB/internal/service"

	"github.com/go-chi/chi/v5"
)

type ConfigServer struct {
	// ServerAddress es la dirección en la que se ejecutará el servidor
	ServerAddress string
	// StaticFilesPath es la ruta de los archivos estáticos
	StaticFilesPath string
}

type Server struct {
	// ServerAddress es la dirección en la que se ejecutará el servidor
	serverAddress string
	// StaticFilesPath es la ruta de los archivos estáticos
	staticFilesPath string
}

func NewServer(cfg *ConfigServer) *Server {

	defaultConfig := &ConfigServer{
		ServerAddress:   ":8080",
		StaticFilesPath: "./docs/db",
	}

	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
		if cfg.StaticFilesPath != "" {
			defaultConfig.StaticFilesPath = cfg.StaticFilesPath
		}
	}

	return &Server{
		serverAddress:   defaultConfig.ServerAddress,
		staticFilesPath: defaultConfig.StaticFilesPath,
	}

}

func (s *Server) Run() error {
	pr, err := repository.NewProductRepository(s.staticFilesPath)
	if err != nil {
		return fmt.Errorf("Error al crear el repositorio de productos: %s", err.Error())
	}

	ps, err := service.NewProductService(pr)
	if err != nil {
		return fmt.Errorf("Error al crear el servicio de productos: %s", err.Error())
	}

	ph := handlers.NewProductHandler(ps)

	router := chi.NewRouter()

	//router.Use(middleware.Logger)
	//router.Use(middleware.Recoverer)

	router.Group(func(router chi.Router) {
		router.Get("/ping", ph.HandlerPing)
	})

	router.Route("/products", func(router chi.Router) {

		router.Group(func(router chi.Router) {
			router.Get("/", ph.HandlerGetAllProduct)
			router.Get("/{id}", ph.HandlerGetProductByID)
			router.Get("/search", ph.HandlerSearchProductByPrice)
		})

		router.Group(func(router chi.Router) {
			router.Post("/", ph.HandlerCreateProduct)
			router.Patch("/{id}", ph.HandlerUpdatePartialProduct)
			router.Put("/{id}", ph.HandlerUpdateProduct)
			router.Delete("/{id}", ph.HandlerDeleteProduct)
		})

	})

	return http.ListenAndServe(s.serverAddress, router)

}
