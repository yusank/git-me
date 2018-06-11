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
 *     Initial: 2018/04/01        Yusan Kurban
 */

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/yusank/git-me/common"
	"github.com/yusank/git-me/extractors"
	"github.com/yusank/git-me/model"
	"github.com/yusank/git-me/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	OutputDir   string
	inputReader *bufio.Reader
	exitChan    = make(chan int)
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "git-me",
	Short: "git-me, Give It To Me.",
	Long: `	
	git-me is a command-line tool which provide download service.
	This tool has nothing to do with GIT or any other version control tool.
	Git-me only focus on get media from web site to your computer.

	Here is simple use example:
		git-me https://bilibili.com/ae232434

	Contact with author if you get any trouble  while using this tool.
	Yusank - yusankurban@gmail.com`,
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Name == "" && len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// listen channel ProcessChan and exitChan
		go func() {
			for {
				select {
				case u := <-common.ProcessChan:
					// only upload when have user account info
					if common.Name != "" && common.Pass != "" {
						resp := model.InnerTaskResp{
							Name:     common.Name,
							Pass:     utils.StringMd5(common.Pass),
							URL:      u.URL,
							Event:    u.Status,
							Schedule: u.Schedule,
						}

						if err := model.UploadCurrentTaskStatus(resp); err != nil {
							fmt.Println("upload err :", err)
						}
					}
					continue
				case <-exitChan:
					fmt.Println("exit")
					break
				}
				break
			}
		}()

		// init http-client
		utils.InitHttpClient()

		//utils.SetProxy(ProxyPort)
		// init map
		extractors.BeforeRun()

		tasks := handleUserTask()
		if len(tasks) > 0 {
			for _, v := range tasks {
				extractors.MatchUrl(v, OutputDir)
			}
			exitChan <- 1
			return
		}

		if len(args) < 1 && len(tasks) < 1 {
			fmt.Println("there is nothing in your task list.")
			return
		}

		uri := args[0]
		extractors.MatchUrl(uri, OutputDir)
		exitChan <- 1
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// PersistentFlags 是全局参数，即在所有的子命令也有效
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	RootCmd.Flags().StringP("author", "a", "YusanK", "Author name for copyright attribution")
	RootCmd.Flags().StringVarP(&OutputDir, "outputDir", "o", "current dir.", "The path you want save the file.")
	RootCmd.Flags().StringVarP(&utils.HttpProxy, "proxyPort", "x", "", "use http proxy agency when you need.")
	RootCmd.Flags().StringVarP(&utils.Socks5Proxy, "socketProxy", "s", "", "use socket proxy agency when you need.")
	RootCmd.Flags().StringVarP(&utils.Cookie, "cookie", "c", "", "use agency when you need.")
	RootCmd.Flags().StringVarP(&common.Name, "name", "u", "", "account info of tool")
	RootCmd.Flags().StringVarP(&common.Pass, "password", "p", "", "account pass.")
	RootCmd.Flags().StringVarP(&common.Format, "format", "f", "best quality", "format of download media")
	RootCmd.Flags().BoolVarP(&common.InfoOnly, "info", "i", false, "display all available format for choice")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	u := viper.Get("user")
	if u != nil {
		common.Name = u.(string)
	}

	p := viper.Get("pass")
	if p != nil {
		common.Pass = p.(string)
	}

	fmt.Println(common.Name, common.Pass)
}

func handleUserTask() []string {
	if common.Name == "" || common.Pass == "" {
		return nil
	}

	var u model.InnerTaskResp
	u.Name = common.Name
	u.Pass = utils.StringMd5(common.Pass)

	urls, err := model.GetUserTask(u)
	if err != nil {
		fmt.Println("user info err", err)
		return nil
	}

	return urls
}
