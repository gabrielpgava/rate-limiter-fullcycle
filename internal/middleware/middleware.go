package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/limiter/limiterIP"
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/limiter/limiterToken"
)

func RateLimiterMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))

		if token == "" {
			checkedIP, err := limiterIP.CheckIPLimit(ip)
			if err != nil {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
			fmt.Println("IP:", checkedIP)
		}

		if token != "" {
			checkToken, err := limiterToken.CheckTokenLimit(token)
			if err != nil {
				if err.Error() == "IsInvalid" {
					http.Error(w, "Invalid API Key Token", http.StatusUnauthorized)
					return
				}

				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			fmt.Println("Token:", checkToken)
		}

		next.ServeHTTP(w, r)
	})

}
