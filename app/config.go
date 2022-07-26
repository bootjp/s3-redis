package app

type Config struct {
	Version       string `yaml:"version"`
	UseCache      bool
	UseDirtyWrite bool
	Redis         struct {
		Port     string      `yaml:"port"`
		Bind     string      `yaml:"bind"`
		Password interface{} `yaml:"password"`
	} `yaml:"redis"`
	S3 struct {
		Version   string `yaml:"version"`
		Bucket    string `yaml:"bucket"`
		Region    string `yaml:"region"`
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Prefix    string `yaml:"prefix"`
		Endpoint  string `yaml:"endpoint"`
	} `yaml:"s3"`
}
