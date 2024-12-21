/*
Copyright © 2024 zhanghuangbin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var (
	storeFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sparkle",
	Short: "一个快速修改切换系统环境变量的命令行工具",
	Long: `sparkle 是一个快速修改切换系统环境变量的命令行工具。它的实现思路如下：
1.注册一系列环境变量的信息到一个持久化文件中，并为这些变量信息取一个别名
2.修改环境变量时，通过这些别名，从而减少修改变量时需要输入的信息，从而快速的切换环境变量
3.同时，sparkle还支持对环境变量进行分组，从而可以快速切换到某个分组下的环境变量
4.此外，sparkle在修改系统环境变量的同时，会更新当前bash的环境变量，从而无需重启bash
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&storeFile, "store", "s", "", "store file used by aliasCommand(default is $HOME/.sparkle.yaml)")

	rootCmd.AddCommand(aliasCmd)
}

func initConfig() {
	if storeFile != "" {
		viper.SetConfigFile(storeFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configFile := path.Join(home, ".sparkle.yaml")
		// 检查文件是否存在，不存在则创建
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(configFile), 0755); err != nil {
				cobra.CheckErr(err)
			}
			if err := os.WriteFile(configFile, []byte{}, 0644); err != nil {
				cobra.CheckErr(err)
			}
			fmt.Fprintln(os.Stdout, "store file not exists, create empty store file:", configFile)
		}

		viper.SetConfigFile(configFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stdout, "Using store file:", viper.ConfigFileUsed())
	}
}
