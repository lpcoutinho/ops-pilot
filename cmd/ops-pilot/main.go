package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	debugMode     bool
	dangerousMode bool
)

var rootCmd = &cobra.Command{
	Use:   "ops-pilot",
	Short: "Ops-Pilot is an AI-powered CLI for Linux administration",
	Long: `Ops-Pilot is an open-source tool that acts as a natural language co-pilot
for Linux system administration and auditing. It translates user intent into 
safe system commands using LLMs.`,
}

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask the pilot a question or give a command",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]
		
		// In a real implementation, we would choose the provider based on config
		// For now, we'll use a placeholder or the first available implementation
		fmt.Printf("Analyzing: %s\n", question)
		fmt.Println("(Note: Integration with real LLM APIs will require setting LLM_API_KEY)")
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ops-pilot.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&dangerousMode, "dangerous-mode", false, "allow execution of potentially dangerous commands")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("dangerous_mode", rootCmd.PersistentFlags().Lookup("dangerous-mode"))

	rootCmd.AddCommand(askCmd)
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
