package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/arkste/elsi/utils"
	"github.com/spf13/cobra"
)

var fsSourceDir string
var fsExcludeFiles []string
var fsSizeLimit int
var fsPipeline bool

// Document Struct
type Document struct {
	ID        string    `json:"id,omitempty"`
	FilePath  string    `json:"filepath,omitempty"`
	Name      string    `json:"name,omitempty"`
	Extension string    `json:"extension,omitempty"`
	Path      string    `json:"path,omitempty"`
	Size      int64     `json:"size,omitempty"`
	Mode      string    `json:"filemode,omitempty"`
	ModeNum   string    `json:"filemode_num,omitempty"`
	ModTime   time.Time `json:"mod_time,omitempty"`
	IsDir     bool      `json:"is_dir,omitempty"`
	Data      string    `json:"data,omitempty"`
}

var filesystemCmd = &cobra.Command{
	Use:   "fs",
	Short: "Filesystem Indexer",
	RunE: func(cmd *cobra.Command, args []string) error {
		EsClient.UsePipeline = fsPipeline
		EsClient.Init()

		err := filepath.Walk(fsSourceDir, func(path string, info os.FileInfo, err error) error {
			// Skip Excludes
			for _, pattern := range fsExcludeFiles {
				match, err := filepath.Match(pattern, info.Name())
				if err != nil {
					return fmt.Errorf("bad pattern provided %s", pattern)
				}
				if match {
					return nil
				}
			}

			// don't believe info
			fileStat, err := os.Stat(path)
			if err != nil {
				return nil
			}

			// skip dirs
			if fileStat.IsDir() {
				return nil // filepath.SkipDir ?
			}

			var fileEncodedContent, filePipeline string
			// check if pipeline processor is enabled
			if fsPipeline {
				// open file, only if its <= fileSizeLimit MB
				if fileStat.Size() <= int64(fsSizeLimit*1024*1024) {
					f, err := ioutil.ReadFile(path)
					if err != nil {
						log.Printf("File could not be opened %s: %v", path, err)
						return nil
					}

					// convert file content to base64
					fileEncodedContent = base64.StdEncoding.EncodeToString(f)
					filePipeline = EsClient.PipelineName
				}
			}

			// prepare elasticsearch document
			document := Document{
				ID:        utils.CreateHashFromString(path),
				FilePath:  path,
				Name:      fileStat.Name(),
				Extension: strings.TrimLeft(filepath.Ext(fileStat.Name()), "."),
				Path:      filepath.Dir(path),
				Size:      fileStat.Size(),
				Mode:      fileStat.Mode().String(),
				ModeNum:   fmt.Sprintf("%04o", fileStat.Mode().Perm()),
				ModTime:   fileStat.ModTime(),
				IsDir:     fileStat.IsDir(),
				Data:      fileEncodedContent,
			}

			EsClient.AddDocument(document.ID, document, filePipeline)

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
	rootCmd.AddCommand(filesystemCmd)
	filesystemCmd.Flags().StringVarP(&fsSourceDir, "source", "s", fsSourceDir, "Source directory to read from")
	filesystemCmd.MarkFlagRequired("source")
	filesystemCmd.Flags().StringSliceVarP(&fsExcludeFiles, "exclude", "e", fsExcludeFiles, "Exclude File Patterns, comma separated (eg: \"*.log,*.epub,*.sdf*\")")
	filesystemCmd.Flags().IntVarP(&fsSizeLimit, "limit", "l", 10, "Limit Filesize in MB")
	filesystemCmd.Flags().BoolVarP(&fsPipeline, "pipeline", "p", false, "Use Elasticsearch Pipeline Processor (Ingest Attachment Plugin required)")
}
