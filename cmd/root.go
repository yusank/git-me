package cmd

import (
	"errors"
	"fmt"
	"os"

	"git-me/common"
	"git-me/extractors"
	"git-me/utils"

	"git-me/model"

	"bufio"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	OutputDir   string
	inputReader *bufio.Reader
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "git-me",
	Short: "git-me, Give It To Me.",
	Long: `git-me is a command-line tool which provide download service.
	This tool has nothing to do with git or any other version control tool.
	Git-me only focus on get media from web site to your computer.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if common.Name == "" && len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
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
			return
		}

		if len(args) < 1 {
			return
		}

		uri := args[0]
		extractors.MatchUrl(uri, OutputDir)
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
	RootCmd.Flags().StringVarP(&OutputDir, "outputDir", "o", ".", "The path you want save the file.")
	RootCmd.Flags().StringVarP(&utils.HttpProxy, "proxyPort", "x", "", "use agency when you need.")
	RootCmd.Flags().StringVarP(&utils.Socks5Proxy, "socketProxy", "s", "", "use agency when you need.")
	RootCmd.Flags().StringVarP(&utils.Cookie, "cookie", "c", "", "use agency when you need.")
	RootCmd.Flags().StringVarP(&common.Name, "name", "u", "", "account info of tool")
	RootCmd.Flags().StringVarP(&common.Pass, "password", "p", "", "account pass.")
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func handleUserTask() []string {
	if common.Name == "" || common.Pass == "" {
		return nil
	}

	var u model.InnerTaskResp
	u.Name = common.Name
	u.Pass = common.Pass

	urls, err := model.GetUserTask(u)
	if err != nil {
		fmt.Println("user info err", err)
		return nil
	}

	return urls
}
