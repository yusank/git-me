/*
 * MIT License
 *
 * Copyright (c) 2018 Yusan Kurban <yusankurban@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/06/01        Yusan Kurban
 */

package cmd

import (
	"errors"

	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yusank/git-me/common"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configuration of tool",
	Long:  "you can add / rm any config",
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Name == "" && len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			s := strings.Split(v, "=")
			if len(s) != 2 {
				continue
			}

			addToConfigFile(s[0], s[1])
		}
	},
}

//func init() {
//	RootCmd.AddCommand(configCmd)
//}

// initConfig reads in config file and ENV variables if set.
func addToConfigFile(key, value string) {
	viper.Set(key, value)
}
