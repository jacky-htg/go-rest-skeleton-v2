package httpresponse

import (
	"context"
	"net/http"
	"rest-skeleton/pkg/redis"

	"github.com/bytedance/sonic"
)

type Response struct {
	Cache *redis.Cache
}

func (r Response) SetMarshal(ctx context.Context, w http.ResponseWriter, statusCode int, response interface{}, key string) {
	data, err := sonic.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
	if len(key) > 0 {
		r.Cache.Add(ctx, key, data)
	}
}

func (r Response) Set(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(response.(string)))
}
