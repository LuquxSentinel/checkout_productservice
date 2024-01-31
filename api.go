package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/luqu/productservice/types"
)

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error

type APIServer struct {
	listenAddress string
	router        *mux.Router
	service       Service
}

func NewAPIServer(listenAddress string, service Service) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		router:        mux.NewRouter(),
		service:       service,
	}
}

func (api *APIServer) Run() error {

	// get product request
	api.router.HandleFunc("/product/{product-id}", handler(api.getProduct))
	api.router.HandleFunc("/products", handler(api.getAllProducts))

	return http.ListenAndServe(api.listenAddress, api.router)
}

func handler(fn APIFunc) http.HandlerFunc {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request-id", uuid.New())
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(ctx, w, r)
		if err != nil {
			writeErrorResponse(w, err)
		}
	}
}

func (api *APIServer) getProduct(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	// get product id from request
	productID := mux.Vars(r)["product-id"]
	if productID == "" {
		return &types.Error{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid product id",
		}
	}

	// call the product service to get the product with id
	product, err := api.service.GetProduct(ctx, productID)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "product with id not found",
		}
	}

	err = writeResponse(w, product)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "response error",
		}
	}

	return nil

}

func (api *APIServer) getAllProducts(ctx context.Context, w http.ResponseWriter, r *http.Request) *types.Error {
	products, err := api.service.GetAllProducts(ctx)
	if err != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to fetch products",
		}
	}

	err1 := writeResponse(w, products)
	if err1 != nil {
		return &types.Error{
			StatusCode: http.StatusInternalServerError,
			Message:    "response failed",
		}
	}

	return nil
}

func writeErrorResponse(w http.ResponseWriter, err *types.Error) error {
	w.WriteHeader(err.StatusCode)
	return json.NewEncoder(w).Encode(err.Error())
}

func writeResponse(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}
