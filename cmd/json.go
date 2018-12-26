package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/arkste/elsi/utils"
	"github.com/spf13/cobra"
)

var jsonSourceDir string

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON Indexer",
	RunE: func(cmd *cobra.Command, args []string) error {
		EsClient.Init()

		err := filepath.Walk(jsonSourceDir, func(path string, info os.FileInfo, err error) error {
			// skip all files which dont have a "*.json" file extension
			if filepath.Ext(info.Name()) != ".json" {
				return nil
			}

			// don't believe info
			fileStat, err := os.Stat(path)
			if err != nil {
				log.Printf("Could not stat() file %s: %v", path, err)
				return nil
			}

			// skip dirs
			if fileStat.IsDir() {
				return nil // filepath.SkipDir ?
			}

			f, err := ioutil.ReadFile(path)
			if err != nil {
				log.Printf("File could not be opened %s: %v", path, err)
				return nil
			}

			EsClient.AddDocument(utils.CreateHashFromString(path), utils.IsJSONArray(string(f)), "")

			return nil
		})
		if err != nil {
			return err
		}

		EsClient.Flush()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.Flags().StringVarP(&jsonSourceDir, "source", "s", jsonSourceDir, "Source directory to read from")
	jsonCmd.MarkFlagRequired("source")
}
