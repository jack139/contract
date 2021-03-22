package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers contract-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	//registerTxHandlers(clientCtx, r)  // 不注册交易处理的http，另外实现

}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("/contract/contracts/{id}", getContractHandler(clientCtx)).Methods("GET")
	r.HandleFunc("/contract/contracts", listContractHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/contract/contracts", createContractHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/contract/contracts/{id}", updateContractHandler(clientCtx)).Methods("POST")
	r.HandleFunc("/contract/contracts/{id}", deleteContractHandler(clientCtx)).Methods("POST")

}
