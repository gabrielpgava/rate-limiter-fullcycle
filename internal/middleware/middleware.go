package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	limiterIP "github.com/gabrielpgava/rate-limiter-fullcycle/internal/limiter"
)

func RateLimiterMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		ip, _, _ := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))

		checkedIP, err := limiterIP.CheckIPLimit(ip)
		if err != nil {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		fmt.Println("IP:", checkedIP, "Token:", token)




		if(token !=""){
			fmt.Fprintln(w, "Hello, World!")
		}
		

	})
}
