package cloud

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/safedep/vet/internal/auth"
	"github.com/safedep/vet/internal/ui"
	"github.com/safedep/vet/pkg/cloud/query"
	"github.com/safedep/vet/pkg/common/logger"
	"github.com/spf13/cobra"
)

var (
	querySql string
)

func newQueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query risks by executing SQL queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newQueryExecuteCommand())

	return cmd
}

func newQueryExecuteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Execute a query",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := executeQuery()
			if err != nil {
				logger.Errorf("Failed to execute query: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&querySql, "sql", "s", "", "SQL query to execute")

	return cmd
}

func executeQuery() error {
	if querySql == "" {
		return errors.New("SQL string is required")
	}

	client, err := auth.ControlPlaneClientConnection("vet-cloud-query")
	if err != nil {
		return err
	}

	queryService, err := query.NewQueryService(client)
	if err != nil {
		return err
	}

	response, err := queryService.ExecuteSql(querySql)
	if err != nil {
		return err
	}

	return renderQueryResponseAsTable(response)
}

func renderQueryResponseAsTable(response *query.QueryResponse) error {
	tbl := table.NewWriter()
	tbl.SetOutputMirror(os.Stdout)
	tbl.SetStyle(table.StyleLight)

	if response.Count() == 0 {
		logger.Infof("No results found")
		return nil
	}

	ui.PrintSuccess(fmt.Sprintf("Query returned %d results", response.Count()))

	// Header
	headers := []string{}
	response.GetRow(0).ForEachField(func(key string, _ interface{}) {
		headers = append(headers, key)
	})

	sort.Strings(headers)

	headerRow := []interface{}{}
	for _, header := range headers {
		headerRow = append(headerRow, header)
	}

	tbl.AppendHeader(headerRow)

	// Ensure we have a consistent order of columns
	response.ForEachRow(func(row *query.QueryRow) {
		rowValues := []interface{}{}
		for _, header := range headers {
			rowValues = append(rowValues, row.GetField(header))
		}

		tbl.AppendRow(rowValues)
		tbl.AppendSeparator()
	})

	tbl.Render()
	return nil
}
