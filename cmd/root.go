package cmd

import (
	"fmt"
	"os"

	"github.com/Kukuruzoed/config-analyzer/internal/analyzer"
	"github.com/Kukuruzoed/config-analyzer/internal/config"
	"github.com/Kukuruzoed/config-analyzer/internal/grpcserver"
	"github.com/Kukuruzoed/config-analyzer/internal/parser"
	"github.com/Kukuruzoed/config-analyzer/internal/rules"
	"github.com/Kukuruzoed/config-analyzer/internal/server"
	"github.com/spf13/cobra"
)

var (
	silent    bool
	fromStdin bool
	serve     bool
	port      int
	grpcServe bool
	grpcPort  int
)

var rootCmd = &cobra.Command{
	Use:   "config-analyzer [file]",
	Short: "Анализирует конфиг на небезопасные настройки",
	RunE:  run,
}

// Execute — единственная публичная функция пакета,
// её вызывает main.go
func Execute() error {
	return rootCmd.Execute()
}

// init вызывается Go автоматически при импорте пакета —
// это стандартное место для регистрации флагов в cobra
func init() {
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "не выходить с ошибкой при наличии проблем")
	rootCmd.Flags().BoolVar(&fromStdin, "stdin", false, "читать конфиг из stdin")
	rootCmd.Flags().BoolVar(&serve, "serve", false, "запустить HTTP-сервер")
	rootCmd.Flags().IntVar(&port, "port", 8080, "порт HTTP-сервера")
	rootCmd.Flags().BoolVar(&grpcServe, "grpc", false, "запустить gRPC-сервер")
	rootCmd.Flags().IntVar(&grpcPort, "grpc-port", 50051, "порт gRPC-сервера")
}

func run(cmd *cobra.Command, args []string) error {
	a := analyzer.New(analyzer.DefaultRules())

	if serve {
		srv := server.New(a, port)
		return srv.Run()
	}

	if grpcServe {
		srv := grpcserver.New(a, grpcPort)
		return srv.Run()
	}

	return runCLI(a, args)
}

func runCLI(a *analyzer.Analyzer, args []string) error {
	var (
		cfg config.Config
		err error
	)

	if fromStdin {
		cfg, err = parser.ParseReader(os.Stdin)
	} else {
		if len(args) == 0 {
			return fmt.Errorf("укажи путь к файлу или используй --stdin")
		}
		cfg, err = parser.ParseFile(args[0])
	}
	if err != nil {
		return fmt.Errorf("ошибка чтения конфига: %w", err)
	}

	issues := a.Analyze(cfg)
	printIssues(issues)

	if len(issues) > 0 && !silent {
		os.Exit(1)
	}
	return nil
}

func printIssues(issues []rules.Issue) {
	if len(issues) == 0 {
		fmt.Println("проблем не найдено")
		return
	}
	for _, issue := range issues {
		fmt.Printf("[%s] %s\n  → %s\n\n",
			issue.Severity,
			issue.Description,
			issue.Recommendation,
		)
	}
}
