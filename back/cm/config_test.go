package cm

import (
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConf(t *testing.T) {
	assert.NotZero(t, viper.GetInt("app.back_port"))
	assert.NotZero(t, viper.GetString("app.data_dir"))

	assert.NotZero(t, viper.GetString("pocket.token"))
	assert.NotZero(t, viper.GetInt("pocket.interval"))

	assert.NotZero(t, viper.GetString("aai.token"))

}

func TestSetConfig(t *testing.T) {
	viper.Set("pocket.app_info", "old_app_info")

	deviceId := strings.ToUpper(uuid.NewString())
	config := &Config{
		Pocket: PocketConfig{
			AppInfo: PocketAppInfo{
				DeviceId: deviceId,
			},
		},
	}

	changedCount := SetConfig(config)
	assert.Equal(t, 1, changedCount)
	assert.Equal(t, deviceId, viper.GetString("pocket.app_info.deviceId"))
}
