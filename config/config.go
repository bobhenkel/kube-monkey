package config

import (
	"fmt"
	"github.com/asobti/kube-monkey/config/param"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"k8s.io/client-go/1.5/pkg/util/sets"
	"time"
)

const (
	configpath = "/etc/kube-monkey"
	configtype = "toml"
	configname = "config"

	// Currently, there does not appear to be
	// any value in making these configurable
	// so defining them as consts
	IdentLabelKey     = "kube-monkey/identifier"
	EnabledLabelKey   = "kube-monkey/enabled"
	EnabledLabelValue = "enabled"
	MtbfLabelKey      = "kube-monkey/mtbf"

	KubeSystemNamespace = "kube-system"
)

func SetDefaults() {
	viper.SetDefault(param.DryRun, true)
	viper.SetDefault(param.Timezone, "America/Los_Angeles")
	viper.SetDefault(param.RunHour, 8)
	viper.SetDefault(param.StartHour, 10)
	viper.SetDefault(param.EndHour, 16)
	viper.SetDefault(param.GracePeriodSec, 5)
	viper.SetDefault(param.BlacklistedNamespaces, []string{KubeSystemNamespace})

	viper.SetDefault(param.DebugEnabled, false)
	viper.SetDefault(param.DebugScheduleDelay, 30)
	viper.SetDefault(param.DebugForceShouldKill, false)
	viper.SetDefault(param.DebugScheduleImmediateKill, false)
}

func setupWatch() {
	// TODO: This does not appear to be working
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config change detected")
		ValidateConfigs()
	})
}

func Init() error {
	SetDefaults()
	viper.AddConfigPath(configpath)
	viper.SetConfigType(configtype)
	viper.SetConfigName(configname)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	ValidateConfigs()
	setupWatch()
	return nil
}

func DryRun() bool {
	return viper.GetBool(param.DryRun)
}

func Timezone() *time.Location {
	tz := viper.GetString(param.Timezone)
	location, err := time.LoadLocation(tz)
	if err != nil {
		panic(err.Error())
	}
	return location
}

func RunHour() int {
	return viper.GetInt(param.RunHour)
}

func StartHour() int {
	return viper.GetInt(param.StartHour)
}

func EndHour() int {
	return viper.GetInt(param.EndHour)
}

func GracePeriodSeconds() *int64 {
	gpInt64 := viper.GetInt64(param.GracePeriodSec)
	return &gpInt64
}

func BlacklistedNamespaces() sets.String {
	// Return as set for O(1) membership checks
	namespaces := viper.GetStringSlice(param.BlacklistedNamespaces)
	return sets.NewString(namespaces...)
}

func DebugEnabled() bool {
	return viper.GetBool(param.DebugEnabled)
}

func DebugScheduleDelay() time.Duration {
	delaySec := viper.GetInt(param.DebugScheduleDelay)
	return time.Duration(delaySec) * time.Second
}

func DebugForceShouldKill() bool {
	return viper.GetBool(param.DebugForceShouldKill)
}

func DebugScheduleImmediateKill() bool {
	return viper.GetBool(param.DebugScheduleImmediateKill)
}
