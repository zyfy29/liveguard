package pocket

import (
	"github.com/spf13/viper"
	"sync"
)

var clientCreateLock = &sync.Mutex{}

var instance *Client // singleton

// GetClient returns a singleton instance of the Client, ensuring thread-safe initialization.
func GetClient() *Client {
	if instance == nil {
		clientCreateLock.Lock()
		defer clientCreateLock.Unlock()
		if instance == nil {
			instance = newClient(
				viper.GetString("pocket.token"),
				viper.GetStringMapString("pocket.app_info"),
				viper.GetInt("pocket.interval"),
			)
		}
	}
	return instance
}
