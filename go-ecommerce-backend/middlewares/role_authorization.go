package middlewares

import (
	"net/http"
	"slices"
)

func RoleAuthorizationMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextUserRole).(string)
			if !ok {
				http.Error(w, "User role not found", http.StatusUnauthorized)
				return
			}
			if slices.Contains(allowedRoles, role) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
