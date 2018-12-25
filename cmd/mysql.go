package cmd

import (
	"database/sql"
	"log"

	"github.com/arkste/elsi/utils"
	_ "github.com/go-sql-driver/mysql" // Import MySQL Driver
	"github.com/spf13/cobra"
)

var mysqlDSN string
var mysqlQuery string

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "MySQL Indexer",
	Run: func(cmd *cobra.Command, args []string) {
		EsClient.Init()

		db, err := sql.Open("mysql", mysqlDSN)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		rows, err := db.Query(mysqlQuery)
		if err != nil {
			log.Fatalln(err)
		}

		columns, err := rows.Columns()
		if err != nil {
			log.Fatalln(err)
		}

		count := len(columns)
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for rows.Next() {
			for i := 0; i < count; i++ {
				valuePtrs[i] = &values[i]
			}

			rows.Scan(valuePtrs...)

			id, jsonDataString := utils.ConvertToJSON(columns, values)

			EsClient.AddDocument(id, jsonDataString, "")
		}

		EsClient.Flush()
	},
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
	mysqlCmd.Flags().StringVarP(&mysqlDSN, "dsn", "d", "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&collation=utf8mb4_unicode_ci", "MySQL DSN")
	mysqlCmd.MarkFlagRequired("dsn")
	mysqlCmd.Flags().StringVarP(&mysqlQuery, "query", "q", "SELECT * FROM table", "MySQL Query")
	mysqlCmd.MarkFlagRequired("query")
}
