package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OwnerVerifierFunc func(resourceID int) (int, error)

func OwnerOnlyMiddleware(paramName string, getOwnerID OwnerVerifierFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(ContextUserID).(int)
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(r)
			resourceIDStr, exists := vars[paramName]
			if !exists {
				http.Error(w, "Resource identifier not provided", http.StatusBadRequest)
				return
			}

			resourceID, err := strconv.Atoi(resourceIDStr)
			if err != nil {
				http.Error(w, "Invalid resource identifier", http.StatusBadRequest)
				return
			}

			ownerID, err := getOwnerID(resourceID)
			if err != nil {
				http.Error(w, "Resource not found", http.StatusNotFound)
				return
			}

			if userID != ownerID {
				http.Error(w, "Not authorized to access this resource", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
