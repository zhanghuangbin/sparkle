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
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"strings"
)

type EnvType int

const (
	OVERWRITE EnvType = iota
	APPEND    EnvType = iota
)

var (
	format     string
	desc       string
	longDesc   string
	aliasType  int
	aliasKey   string
	aliasValue string
)

type Alias struct {
	Alias    string  `json:"alias"`
	Desc     string  `json:"desc"`
	LongDesc string  `json:"longDesc"`
	Type     EnvType `json:"type"`
	Key      string  `json:"key"`
	Value    string  `json:"value"`
}

type AliasList []Alias

func (list *AliasList) Add(alias *Alias) bool {
	idx := list.IndexOf(alias.Alias)
	if idx == -1 {
		// 添加alias到列表
		*list = append(*list, *alias)
		return true
	}
	// 修改alias
	(*list)[idx] = *alias
	return false
}

func (list *AliasList) IndexOf(alias string) int {
	for idx, v := range *list {
		if v.Alias == alias {
			return idx
		}
	}
	return -1
}

func (list *AliasList) Query(alias string, key string) []Alias {
	var result = make([]Alias, 0, len(*list))
	for _, v := range *list {
		if (alias == "" || v.Alias == alias) && (key == "" || strings.Contains(v.Key, key)) {
			result = append(result, v)
		}
	}
	return result
}

func (list *AliasList) Remove(alias string) bool {
	idx := list.IndexOf(alias)
	if idx == -1 {
		return false
	}
	*list = append((*list)[:idx], (*list)[idx+1:]...)
	return true
}

var (
	aliasCmd = &cobra.Command{
		Use:   "alias",
		Short: "别名信息命令",
		Long:  `对别名信息进行操作，包括查询,新增，修改，删除`,
	}

	listAliasCmd = &cobra.Command{
		Use:   "list [aliasName]",
		Short: "列举别名信息",
		Long:  "查询别名信息",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var list AliasList
			err = viper.UnmarshalKey("alias", &list)
			if err != nil {
				return
			}

			alias := ""
			if len(args) > 0 {
				alias = args[0]
			}
			list = list.Query(alias, aliasKey)

			fullTpl := fmt.Sprintf(`{{range $index, $value := .}}%s{{end}}`, format)
			var tpl *template.Template
			if tpl, err = template.New("tpl").Parse(fullTpl); err != nil {
				return
			}

			var buf bytes.Buffer
			if err = tpl.Execute(&buf, list); err != nil {
				return
			}

			fmt.Printf("总共%d个别名：\n", len(list))
			fmt.Println(string(buf.Bytes()))

			return
		},
	}
	// 增加sparkle alias add aliasName -k "" -v "" --desc "" --long-desc "" -t 0|1
	addAliasCmd = &cobra.Command{
		Use:   "add aliasName",
		Short: "增加别名",
		Long:  `增加/修改别名`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var list AliasList
			err = viper.UnmarshalKey("alias", &list)
			if err != nil {
				return
			}
			list.Add(&Alias{
				Alias:    args[0],
				Key:      aliasKey,
				Value:    aliasValue,
				Type:     EnvType(aliasType),
				Desc:     desc,
				LongDesc: longDesc,
			})
			viper.Set("alias", &list)

			err = viper.WriteConfig()
			return
		},
	}

	removeAliasCmd = &cobra.Command{
		Use:   "rm aliasName",
		Short: "删除别名",
		Long:  `删除别名`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var list AliasList
			err = viper.UnmarshalKey("alias", &list)
			if err != nil {
				return
			}
			list.Remove(args[0])
			viper.Set("alias", &list)
			err = viper.WriteConfig()
			return
		},
	}
)

func init() {
	const defaultFormat string = `别名{{$index}}:
	alias:{{.Alias}}
	key:{{.Key}}
	value:{{.Value}}
	type:{{.Type}}
`

	listAliasCmd.Flags().StringVarP(&format, "format", "f", defaultFormat, "list alias format(use 'text/template' syntax). available values:{{$index}}{{.Alias}}{{.Key}}{{.Value}}{{.Type}}")
	listAliasCmd.Flags().StringVarP(&aliasKey, "key", "k", "", "filter the return result with the key value")
	aliasCmd.AddCommand(listAliasCmd)

	addAliasCmd.Flags().StringVarP(&aliasKey, "key", "k", "", "alias key")
	addAliasCmd.MarkFlagRequired("key")
	addAliasCmd.Flags().StringVarP(&aliasValue, "value", "v", "", "alias value")
	addAliasCmd.MarkFlagRequired("value")
	addAliasCmd.Flags().IntVarP(&aliasType, "type", "t", -1, "alias type")
	addAliasCmd.MarkFlagRequired("type")
	addAliasCmd.Flags().StringVarP(&desc, "desc", "", "", "alias short desc")
	addAliasCmd.Flags().StringVarP(&longDesc, "long-desc", "", "", "alias long desc")
	aliasCmd.AddCommand(addAliasCmd)

	aliasCmd.AddCommand(removeAliasCmd)

}
