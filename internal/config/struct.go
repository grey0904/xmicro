package config

type AppConfig struct {
	Nacos Nacos `yaml:"nacos"`
	Mysql Mysql `yaml:"mysql"`
	Redis Redis `yaml:"redis"`
}

type Nacos struct {
	Endpoints           []string `yaml:"endpoints"`
	Username            string   `yaml:"username"`
	Password            string   `yaml:"password"`
	TimeoutMs           uint64   `yaml:"timeoutMs"`
	NamespaceId         string   `yaml:"namespaceId"`
	NotLoadCacheAtStart bool     `yaml:"otLoadCacheAtStart"`
	CacheDir            string   `yaml:"cacheDir"`
	LogDir              string   `yaml:"logDir"`
	LogLevel            string   `yaml:"logLevel"`
}

type Mysql struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Database string `json:"database" yaml:"database"`
}

type Redis struct {
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
	Select   int    `json:"select" yaml:"select"`
}
