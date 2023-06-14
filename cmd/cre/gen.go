package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/goccy/go-yaml"

	"github.com/ychengcloud/cre/api"
	"github.com/ychengcloud/cre/gen"
)

var configPath string

var generateCmd = &cobra.Command{
	Use:     "generate [flags]",
	Short:   "generate go code for the database schema",
	Example: `cre generate -c ./config`,
	Args: func(_ *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig(configPath, strings.ToUpper("cre_"))

		if cfg.Overwrite {
			prompt := &survey.Confirm{
				Message: `[Warning]
The overwrite flag (Overwrite is True) is specified in the configuration file. If Yes is selected, the generated file will overwrite the existing file. 
Are you sure to continue?
				
配置文件中指定了覆写标志(Overwrite is True)，如果选择 Yes， 生成的文件将覆盖已有文件。
确认继续吗?`,
			}

			overwrite := false
			// ask the question
			err := survey.AskOne(prompt, &overwrite)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if err != nil || !overwrite {
				return
			}

		}

		if err := api.Generate(cfg); err != nil {
			fmt.Println("gen error:", err.Error())
			return
		}
		fmt.Println("Done")
	},
}

func init() {
	generateCmd.Flags().StringVarP(&configPath, "config", "c", "./config.yml", "config file path")

	cobra.OnInitialize()
	rootCmd.AddCommand(generateCmd)

}

func loadConfig(path string, prefix string) *gen.Config {
	var (
		v = viper.New()
	)

	v.AddConfigPath(".")
	v.AddConfigPath("./")
	v.AddConfigPath("/etc/")   // path to look for the config file in
	v.AddConfigPath("$HOME/.") // call multiple times to add many search paths

	v.AutomaticEnv()
	v.SetEnvPrefix(prefix)

	conf := &gen.Config{}

	//读取默认配置
	v.SetConfigName(string(path + ".template"))
	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
		if err := v.Unmarshal(conf); err != nil {
			fmt.Printf("unmarshal conf failed, err:%s \n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Can't read default config file -> %s\n", v.ConfigFileUsed())
		os.Exit(1)
	}

	//读取应用配置
	v.SetConfigName(string(path))
	if err := v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		fmt.Printf("unmarshal conf failed, err:%s \n", err)
		os.Exit(1)
	}

	if err := v.Unmarshal(conf); err != nil {
		fmt.Printf("unmarshal conf failed, err:%s \n", err)
		os.Exit(1)
	}

	return conf
}
