package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/lpcoutinho/ops-pilot/internal/agent"
	"github.com/lpcoutinho/ops-pilot/internal/agent/providers"
	"github.com/lpcoutinho/ops-pilot/internal/tools"
	"github.com/lpcoutinho/ops-pilot/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	debugMode     bool
	dangerousMode bool
)

var rootCmd = &cobra.Command{
	Use:   "ops-pilot [question]",
	Short: "Ops-Pilot is an AI-powered CLI for Linux administration",
	Long: `Ops-Pilot acts as a natural language co-pilot for Linux system administration.
If no sub-command is provided, it defaults to 'ask'.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			askFunc(args[0])
		} else {
			cmd.Help()
		}
	},
}

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the pilot a question or give a command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		askFunc(args[0])
	},
}

func askFunc(question string) {
	provider, err := providers.NewProviderFromConfig()
	if err != nil {
		color.Red("❌ Error: Failed to initialize LLM provider: %v", err)
		os.Exit(1)
	}

	v := &validator.CommandValidator{DangerousMode: viper.GetBool("dangerous_mode")}
	a := agent.NewAgent(provider, v)
	a.RegisterTool(&tools.GetSystemHealthTool{})
	a.RegisterTool(&tools.GetTopProcessesTool{})
	a.RegisterTool(&tools.AuditNetworkTool{})
	a.RegisterTool(&tools.AnalyzeLogsTool{})
	a.RegisterTool(&tools.GetHardwareInfoTool{})

	c := color.New(color.FgCyan).Add(color.Bold)
	c.Printf("🚀 Ops-Pilot is analyzing: %s\n", question)

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Thinking..."
	s.Start()

	response, err := a.Process(context.Background(), question)
	s.Stop()

	if err != nil {
		color.Red("\n❌ Agent failed: %v", err)
		os.Exit(1)
	}

	color.New(color.FgHiGreen).Println("\n🤖 Response:")
	fmt.Println(response)
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models for the configured provider",
	Run: func(cmd *cobra.Command, args []string) {
		provider, err := providers.NewProviderFromConfig()
		if err != nil {
			slog.Error("Failed to initialize LLM provider", "error", err)
			os.Exit(1)
		}

		fmt.Printf("🔍 Fetching available models for provider: %s...\n", viper.GetString("llm_provider"))
		
		models, err := provider.ListModels(context.Background())
		if err != nil {
			slog.Error("Failed to list models", "error", err)
			os.Exit(1)
		}

		fmt.Println("\nAvailable Models:")
		for _, m := range models {
			fmt.Printf(" - %s\n", m)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ops-pilot.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&dangerousMode, "dangerous-mode", false, "allow execution of potentially dangerous commands")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("dangerous_mode", rootCmd.PersistentFlags().Lookup("dangerous-mode"))

	// LLM Configs
	viper.SetDefault("llm_provider", "mock")
	
	rootCmd.AddCommand(askCmd)
	rootCmd.AddCommand(modelsCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ops-pilot")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	setupLogging()
}

func setupLogging() {
	level := slog.LevelInfo
	if viper.GetBool("debug") {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("CLI execution failed", "error", err)
		os.Exit(1)
	}
}
