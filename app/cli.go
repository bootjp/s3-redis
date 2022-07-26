package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	cfgFile string
	conf    Config
	rootCmd = &cobra.Command{
		Use:   "s3-redis",
		Short: "Redis scales infinitely with S3",
		//Long: ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			file, err := ioutil.ReadFile(cfgFile)
			if err != nil {
				fmt.Println(cfgFile)
				cobra.CheckErr(err)
			}
			err = yaml.Unmarshal(file, &conf)
			if err != nil {
				cobra.CheckErr(err)
			}

			s3 := NewS3(conf)
			redis := NewRedis(s3, conf)
			s := NewServer(conf, s3, redis)
			return s.Run()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./s3_redis.yaml", "config file (default is $HOME/.cobra.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find current dir.
		dir, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(dir)
		viper.SetConfigName("s3_redis")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("config file read error")
		fmt.Println(err)
		os.Exit(1)
	}
}
