package main

import (
	"fmt"
	"os"

	"github.com/Kukuruzoed/config-analyzer/internal/analyzer"
	"github.com/Kukuruzoed/config-analyzer/internal/parser"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
	"github.com/Kukuruzoed/config-analyzer/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	var silent, fromStdin, serve bool
	var port int

	cmd := &cobra.Command{
		Use:   "config-analyzer [file]",
		Short: "Анализирует конфигурационный файл",
		RunE: func(cmd *cobra.Command, args []string) error {
			var config map[string]any
			var err error

			var analyzer = analyzer.New(analyzer.DefaultRules())

			if serve {
				srv := server.New(analyzer, port)
				return srv.Run() // блокирует до остановки
			}

			if fromStdin {
				config, err = parser.ParseReader(os.Stdin)
			} else {
				if len(args) == 0 {
					return fmt.Errorf("Необходимо указать путь к файлу")
				}
				var filePath = args[0]

				_, err := os.Stat(filePath)

				if err != nil {
					return err
				}

				config, err = parser.ParseFile(filePath)

				if err != nil {
					return err
				}
			}

			var issues = analyzer.Analyze(config)

			printIssues(issues)

			return nil
		},
	}

	cmd.Flags().BoolVar(&serve, "serve", false, "запустить HTTP-сервер")
	cmd.Flags().IntVar(&port, "port", 8080, "порт HTTP-сервера")

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
