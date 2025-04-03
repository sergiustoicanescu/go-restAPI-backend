package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/sergiustoicanescu/go-restAPI-backend/go-ecommerce-backend/utils"
)

func ValidateBody[T any](handler func(w http.ResponseWriter, r *http.Request, body *T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := utils.ValidateStruct(body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		handler(w, r, &body)
	}
}
