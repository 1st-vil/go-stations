package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func GetAccessLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		res := &model.AccessLog{
			Timestamp: start,
			Latency:   time.Since(start).Milliseconds(),
			Path:      r.URL.Path,
		}
		OS, err := OSFromContext(r.Context())
		if err != nil {
			log.Println(err)
			return
		}
		res.OS = OS

		bytes, err := json.Marshal(res)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(bytes))
	}
	return http.HandlerFunc(fn)
}
