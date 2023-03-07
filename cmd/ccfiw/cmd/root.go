package cmd

import (
	"ccfiw/internal/config"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string
var version string

var rootCmd = &cobra.Command{
	Use:   "ccfiw",
	Short: "ccfiw",
	Long:  `ccfiw`,
	RunE:  run,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "path to configuration file (optional)")
	rootCmd.PersistentFlags().Int("log-level", 4, "debug=5, info=4, error=2, fatal=1, panic=0")

	// Bind flag to config vars.
	viper.BindPFlag("general.log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.SetDefault("ccampus.basepath", "https://dashboard.lanaas.eu:18002")
	viper.SetDefault("ccampus.workerssid", "Corp,test1") // splitted by , without spaces
	viper.SetDefault("idm.basepath", "https://auth.iotplatform.telefonica.com:15001")
	viper.SetDefault("iotagent.hostname", "iota.iotplatform.telefonica.com")
	viper.SetDefault("iotagent.iota_port", 8088)
	viper.SetDefault("iotagent.json_port", 8185)
	viper.SetDefault("iotagent.ignore_sites", "")
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.servers", "redis:6379")
	viper.SetDefault("iotagent.force_update", false)
	rootCmd.AddCommand(versionCmd)
}

// Execute executes the root command.
func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	for _, pair := range os.Environ() {
		d := strings.SplitN(pair, "=", 2)
		if strings.Contains(d[0], ".") {
			log.Warning("Using dots in env variable is illegal and deprecated. Please use double underscore `__` for: ", d[0])
			underscoreName := strings.ReplaceAll(d[0], ".", "__")
			// Set only when the underscore version doesn't already exist.
			if _, exists := os.LookupEnv(underscoreName); !exists {
				os.Setenv(underscoreName, d[1])
			}
		}
	}

	viperBindEnvs(config.C)

	if err := viper.Unmarshal(&config.C); err != nil {
		log.WithError(err).Fatal("unmarshal config error")
	}

}

func viperBindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = strings.ToLower(t.Name)
		}
		if tv == "-" {
			continue
		}

		switch v.Kind() {
		case reflect.Struct:
			viperBindEnvs(v.Interface(), append(parts, tv)...)
		default:
			// Bash doesn't allow env variable names with a dot so
			// bind the double underscore version.
			keyDot := strings.Join(append(parts, tv), ".")
			keyUnderscore := strings.Join(append(parts, tv), "__")
			viper.BindEnv(keyDot, strings.ToUpper(keyUnderscore))
		}
	}
}
