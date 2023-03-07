package config

// Config defines the configuration structure.
type Config struct {
	General struct {
		LogLevel int `mapstructure:"log_level"`
	} `mapstructure:"general"`

	Redis struct {
		Servers    []string `mapstructure:"servers"`
		Cluster    bool     `mapstructure:"cluster"`
		MasterName string   `mapstructure:"master_name"`
		PoolSize   int      `mapstructure:"pool_size"`
		Password   string   `mapstructure:"password"`
		Database   int      `mapstructure:"database"`
		TLSEnabled bool     `mapstructure:"tls_enabled"`
		KeyPrefix  string   `mapstructure:"key_prefix"`
	} `mapstructure:"redis"`

	CCAMPUS struct {
		BasePath   string `mapstructure:"basepath"`
		Username   string `mapstructure:"username"`
		Password   string `mapstructure:"password"`
		Workerssid string `mapstructure:"workerssid"`
	} `mapstructure:"ccampus"`

	IDM struct {
		BasePath    string `mapstructure:"basepath"`
		Username    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
		Service     string `mapstructure:"service"`
		ServicePath string `mapstructure:"servicepath"`
	} `mapstructure:"idm"`

	IoTAgent struct {
		HostName    string `mapstructure:"hostname"`
		IoTAPort    int16  `mapstructure:"iota_port"`
		JSONPort    int16  `mapstructure:"json_port"`
		APIKey      string `mapstructure:"apikey"`
		ForceUpdate bool   `mapstructure:"force_update"`
		IgnoreSites string `mapstructure:"ignore_sites"`
	} `mapstructure:"iotagent"`
}

// C holds the global configuration.
var C Config

// Get returns the configuration.
func Get() *Config {
	return &C
}

// Set sets the configuration.
func Set(c Config) {
	C = c
}
