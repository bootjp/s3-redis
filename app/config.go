package app

type Config struct {
	Version       string `mapstructure:"version"`
	UseCache      bool   `mapstructure:"use_cache"`
	UseDirtyWrite bool   `mapstructure:"use_dirty_write"`
	Redis         struct {
		Port     string      `mapstructure:"port"`
		Bind     string      `mapstructure:"bind"`
		Password interface{} `mapstructure:"password"`
	} `mapstructure:"redis"`
	S3 struct {
		Version   string `mapstructure:"version"`
		Bucket    string `mapstructure:"bucket"`
		Region    string `mapstructure:"region"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Prefix    string `mapstructure:"prefix"`
		Endpoint  string `mapstructure:"endpoint"`
	} `mapstructure:"s3"`
}
