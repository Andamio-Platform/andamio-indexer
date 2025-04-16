package address_handlers

import (
	"encoding/json"
	"net/http"

	"log/slog"

	database "github.com/Andamio-Platform/andamio-indexer/database"
	"github.com/gorilla/mux"
)

type AddressRequest struct {
	Address string `json:"address"`
}

type Response struct {
	Message string `json:"message"`
}

func AddAddressHandler(db *database.Database, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var addressRequest AddressRequest
		err := json.NewDecoder(r.Body).Decode(&addressRequest)
		if err != nil {
			logger.Error("failed to decode request body", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Invalid request body"})
			return
		}

		if addressRequest.Address == "" {
			logger.Error("address is required")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{Message: "Address is required"})
			return
		}

		txn := db.Metadata().DB().Begin()
		defer func() {
			if r := recover(); r != nil {
				txn.Rollback()
				panic(r)
			}
		}()

		err = db.Metadata().AddAddress(addressRequest.Address, txn.DB())
		if err != nil {
			txn.Rollback()
			logger.Error("failed to add address to database", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Failed to add address"})
			return
		} else {
			logger.Error("failed to commit transaction", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Failed to add address"})
			return
		}

		if txn != nil {
			err = txn.Commit().Error
			if err != nil {
				txn.Rollback()
				logger.Error("failed to commit transaction", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(Response{Message: "Failed to add address"})
				return
			}
		}
		logger.Info("address added successfully", "address", addressRequest.Address)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{Message: "Address added successfully"})
	}
}

func AddAddressRoute(router *mux.Router, db *database.Database, logger *slog.Logger) {
	router.HandleFunc("/v1/addresses", AddAddressHandler(db, logger)).Methods("POST")
}
