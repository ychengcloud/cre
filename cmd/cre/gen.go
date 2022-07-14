package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"
	// "github.com/goccy/go-yaml"

	"github.com/ychengcloud/cre/api"
	"github.com/ychengcloud/cre/gen"
)

var configPath string

var generateCmd = &cobra.Command{
	Use:     "generate [flags]",
	Short:   "generate go code for the database schema",
	Example: `heidou generate -c ./config.yml`,
	Args: func(_ *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig(configPath)

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

func loadConfig(filename string) *gen.Config {
	d, err := os.ReadFile(filename)
	if err != nil {
		logrus.Fatalf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}

	// 支持环境变量
	d = []byte(os.ExpandEnv(string(d)))

	cfg := &gen.Config{}

	err = yaml.Unmarshal(d, cfg)
	if err != nil {
		logrus.Fatalf("Config unmarshal fail: %s \n", err)
		os.Exit(1)
	}

	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filepath.Base(filename), ext)
	overwriteTemplate := filepath.Join(dir, name+".template"+ext)

	d, err = os.ReadFile(overwriteTemplate)
	if err != nil {
		logrus.Fatalf("Fatal error overwrite template config file: %s \n", err)
		os.Exit(1)
	}

	// 支持环境变量
	d = []byte(os.ExpandEnv(string(d)))

	err = yaml.Unmarshal(d, cfg)
	if err != nil {
		logrus.Fatalf("Config unmarshal overwrite template fail: %s \n", err)
		os.Exit(1)
	}

	return cfg
}
