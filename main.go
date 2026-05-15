package main

import (
	"fmt"
	"os"

	"github.com/Kukuruzoed/config-analyzer/internal/analyzer"
	"github.com/Kukuruzoed/config-analyzer/internal/parser"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
	"github.com/spf13/cobra"
)

func main() {
	var silent bool
	var fromStdin bool

	cmd := &cobra.Command{
		Use:   "config-analyzer [file]",
		Short: "Анализирует конфигурационный файл",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config map[string]any
			var err error

			if fromStdin {
				config, err = parser.ParseReader(os.Stdin)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("Необходимо указать путь к файлу")
				}
				config, err = parser.ParseFile(args[0])
			}

			if err != nil {
				return err
			}

			var analyzer = analyzer.New(analyzer.DefaultRules())

			var issues = analyzer.Analyze(config)

			printIssues(issues)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "Не выходить с ошибкой")
	cmd.Flags().BoolVar(&fromStdin, "stdin", false, "Читать из stdin")
	cmd.Execute()
}

func printIssues(issues []rules.Issue) {
	for _, issue := range issues {
		fmt.Printf("[%s] %s\n -> %s\n\n",
			issue.Severity,
			issue.Description,
			issue.Recommendation,
		)
	}
}
