package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nadeeshame/repograph_platform/internal/config"
	"github.com/nadeeshame/repograph_platform/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cfgFile string
	verbose bool
	rootCmd = &cobra.Command{
		Use:   "repograph-cli",
		Short: "RepoGraph AI - Intelligent Document Processing Platform",
		Long: `RepoGraph AI is an enterprise-grade document processing and RAG system.
It processes multiple file formats, generates summaries, and enables semantic search.`,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(healthCmd)
}

func initConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	logLevel := cfg.App.LogLevel
	if verbose {
		logLevel = "debug"
	}

	if err := logger.Initialize(logLevel); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index documents",
	Long:  `Scan and index documents from a directory into the vector database.`,
	Run: func(cmd *cobra.Command, args []string) {
		directory, err := cmd.Flags().GetString("directory")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting directory flag: %v\n", err)
			return
		}
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting force flag: %v\n", err)
			return
		}

		logger.Info("Starting indexing",
			zap.String("directory", directory),
			zap.Bool("force", force))

		fmt.Printf("📂 Indexing documents from: %s\n", directory)
		fmt.Printf("⚙️  Force reprocess: %v\n\n", force)

		// TODO: Call orchestrator service to process directory
		fmt.Println("✨ Indexing complete!")
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query the knowledge base",
	Long:  `Query the knowledge base using natural language or search for documents.`,
}

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask a question",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]
		topK, err := cmd.Flags().GetInt("top-k")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting top-k flag: %v\n", err)
			return
		}

		logger.Info("Asking question",
			zap.String("question", question),
			zap.Int("top_k", topK))

		fmt.Printf("🤔 Question: %s\n\n", question)

		// TODO: Call query service
		fmt.Println("💡 Answer: [Implementation pending]")
	},
}

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search documents",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		topK, err := cmd.Flags().GetInt("top-k")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting top-k flag: %v\n", err)
			return
		}
		fileType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting type flag: %v\n", err)
			return
		}

		logger.Info("Searching documents",
			zap.String("query", query),
			zap.Int("top_k", topK),
			zap.String("file_type", fileType))

		fmt.Printf("🔍 Searching for: %s\n\n", query)

		// TODO: Call query service
		fmt.Println("📄 Results: [Implementation pending]")
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive query mode",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🎯 Interactive Query Mode")
		fmt.Println("Type your questions or 'exit' to quit")

		// TODO: Implement interactive mode
		fmt.Println("Implementation pending...")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check indexing status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 RepoGraph AI Status")

		// TODO: Call orchestrator service for status
		fmt.Println("Total Documents: 0")
		fmt.Println("Indexed: 0")
		fmt.Println("Pending: 0")
		fmt.Println("Failed: 0")
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check service health",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			return
		}

		fmt.Println("🏥 Health Check")
		fmt.Println()

		ctx := context.Background()

		// Check services
		services := []string{
			"Azure OpenAI",
			"Pinecone",
			"Google Vision",
			"Database",
			"Redis",
		}

		for _, service := range services {
			fmt.Printf("%-20s ", service+":")
			// TODO: Implement actual health checks
			fmt.Println("✅ Healthy")
		}

		_ = ctx
		_ = cfg
	},
}

func init() {
	// Index command flags
	indexCmd.Flags().StringP("directory", "d", "./data/diagrams", "Directory to index")
	indexCmd.Flags().BoolP("force", "f", false, "Force reprocess all documents")

	// Query command flags
	askCmd.Flags().IntP("top-k", "k", 5, "Number of sources to retrieve")
	searchCmd.Flags().IntP("top-k", "k", 10, "Number of results to return")
	searchCmd.Flags().StringP("type", "t", "", "Filter by file type")

	queryCmd.AddCommand(askCmd)
	queryCmd.AddCommand(searchCmd)
	queryCmd.AddCommand(interactiveCmd)
}
