package middleware

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/providerloader"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/providerloader/util"
	"github.com/oauth2-proxy/oauth2-proxy/v7/providers/utils"
)

// middleware that loads the provider and stores it in the context
func NewProviderLoader(loader providerloader.Loader) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

			providerID := utils.FromContext(req.Context())

			provider, err := loader.Load(req.Context(), providerID)
			if err != nil {
				logger.Error(fmt.Sprintf("unable to load provider, id='%s': %s", providerID, err.Error()))
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := util.AppendToContext(req.Context(), provider)
			next.ServeHTTP(rw, req.WithContext(ctx))
		})
	}
}
