package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/arkste/elsi/elasticsearch"
	"github.com/ghodss/yaml"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config File
var cfgFile string

// EsClient is Elasticsearch Client
var EsClient elasticsearch.Client

var rootCmd = &cobra.Command{
	Use:   "elsi",
	Short: "Elasticsearch Indexer (elsi)",
}

// Execute ist the Root Cmd Exection
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.elsi.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalln(err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".elsi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.Unmarshal(&EsClient)

	// @TODO needs to be improved
	f, _ := ioutil.ReadFile(viper.ConfigFileUsed())
	yaml.Unmarshal(f, &EsClient)
	jsonConfig, err := yaml.YAMLToJSON(f)
	if err != nil {
		log.Fatalln(err)
	}

	var objmap map[string]*json.RawMessage
	json.Unmarshal(jsonConfig, &objmap)

	EsClient.Mapping = string(*objmap["mapping"])
	EsClient.Pipeline = string(*objmap["pipeline"])
}
