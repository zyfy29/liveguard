package cm

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConf(t *testing.T) {
	assert.NotZero(t, viper.GetInt("app.back_port"))
	assert.NotZero(t, viper.GetString("app.data_dir"))

	assert.NotZero(t, viper.GetString("pocket.token"))
	assert.NotZero(t, viper.GetString("pocket.app_info"))
	assert.NotZero(t, viper.GetInt("pocket.interval"))

	assert.NotZero(t, viper.GetString("aai.token"))

	assert.NotZero(t, viper.GetString("medium.token"))
}

func TestSetConfig(t *testing.T) {
	viper.Set("pocket.token", "old_token")
	viper.Set("pocket.app_info", "old_app_info")

	config := &Config{
		Pocket: PocketConfig{
			Token:   "new_token",
			AppInfo: "new_app_info",
		},
	}

	changedCount := SetConfig(config)
	assert.Equal(t, 2, changedCount)
	assert.Equal(t, "new_token", viper.GetString("pocket.token"))
	assert.Equal(t, "new_app_info", viper.GetString("pocket.app_info"))
}
