package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func MakeHandler(ctx context.Context, gqs Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	var graphqlEnpoint endpoint.Endpoint = makeGraphqlEndpoint(gqs)

	graphqlHandler := kithttp.NewServer(
		graphqlEnpoint,
		decodeGraphqlRequest,
		encodeResponse,
		opts...,
	)

	return graphqlHandler
}

type graphqlRequest struct {
	query string
}

func makeGraphqlEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(graphqlRequest)
		res := s.Do(ctx, req.query)
		return res, nil
	}
}

var errBadRequest = errors.New("bad request")

func decodeGraphqlRequest(_ context.Context, r *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	var requestBody struct {
		Query string `json:"query"`
	}
	err := decoder.Decode(&requestBody)
	if err != nil {
		fmt.Println("cannot get query")
	}
	queryParam := requestBody.Query
	if queryParam == "" {
		return nil, errBadRequest
	}
	return graphqlRequest{
		query: queryParam,
	}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
