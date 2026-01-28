package limiterIP

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/models"
	"github.com/gabrielpgava/rate-limiter-fullcycle/internal/storage"
)





func CheckIPLimit(ip string) (bool, error) {

	getIPState, exists := storage.GetIPState(ip)
	if !exists {
		newState := &models.IPstate{
			Count:       1,
			WindowStart: time.Now(),
			BannedUntil: time.Time{},
		}
		storage.SetIPState(ip, newState)
		return true, nil
	}

	fmt.Printf("ip %s,count: %d, window_start: %v, banned_until: %v\n",
		ip,
		getIPState.Count,
		getIPState.WindowStart,
		getIPState.BannedUntil)


	if(!getIPState.BannedUntil.IsZero() && time.Now().After(getIPState.BannedUntil)){
		storage.SetIPState(ip,
			&models.IPstate{
				Count:       1,
				WindowStart: time.Now(),
				BannedUntil: time.Time{},
			})
			return true, nil
	}

	if(getIPState.Count >= 5 && getIPState.BannedUntil.IsZero()){

		       duration := os.Getenv("timeout_ip_block_inSeconds")
		       seconds, err := time.ParseDuration(duration + "s")
		       if err != nil {
			       return false, errors.New("invalid timeout_ip_block_inSeconds value")
		       }
		       storage.SetIPState(ip,
			       &models.IPstate{
				       Count:       getIPState.Count,
				       WindowStart: getIPState.WindowStart,
				       BannedUntil: time.Now().Add(seconds),
			       })
		       return false,  errors.New("IP blocked due to too many requests")
	
	}


	
	storage.SetIPState(ip,
		&models.IPstate{
			Count:       getIPState.Count + 1,
			WindowStart: getIPState.WindowStart,
			BannedUntil: getIPState.BannedUntil,
		})

	


	return true, nil
}

