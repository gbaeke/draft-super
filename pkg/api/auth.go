package api

import (
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-10-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func (s *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	// parse subscription id from request
	subscriptionId := r.URL.Query().Get("subscriptionId")
	if subscriptionId == "" {
		s.logger.Infow("Failed to get subscriptionId from request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := resources.NewGroupsClient(subscriptionId)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		s.logger.Error("Error: ", zap.Error(err))
		return
	}
	client.Authorizer = authorizer

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	s.logger.Infow("Getting groups")
	defer cancel()
	groups, err := client.ListComplete(ctx, "", nil)
	if err != nil {
		s.logger.Error("Error: ", zap.Error(err))
		return
	}

	// output all groups
	for groups.NotDone() {
		w.Write([]byte(*groups.Value().Name))
		err = groups.NextWithContext(ctx)
		if err != nil {
			s.logger.Error("Error: ", zap.Error(err))
		}
	}
}
