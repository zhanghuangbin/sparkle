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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zhanghuangbin/sparkle/meta"
	"runtime"
)

var (
	global bool
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "修改环境变量，并更新shell的环境变量",
	Long:  `修改环境变量，并更新shell的环境变量`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var list meta.AliasList
		if err := viper.UnmarshalKey("alias", &list); err != nil {
			return err
		}

		alias := list.Get(args[0])
		if alias == nil {
			return errors.New(fmt.Sprintf("别名%s不存在", args[0]))
		}

		pType, _ := meta.OfOSType(runtime.GOOS)
		env, err := meta.New(pType, global)
		if err != nil {
			return err
		}

		if err := env.Apply(*alias); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	envCmd.Flags().BoolVarP(&global, "global", "g", false, "是否全局生效")

	rootCmd.AddCommand(envCmd)
}
