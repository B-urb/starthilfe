package main

import (
	"github.com/B-urb/starthilfe/pkg/gitops"
	"github.com/B-urb/starthilfe/pkg/projectconfig"
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		slog.Info("Usage: starthilfe <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		cfg := projectconfig.DefaultConfig()
		if err := projectconfig.SaveConfig(cfg, "starthilfe.yml"); err != nil {
			slog.Error("Failed to save default config", "error", err)
		} else {
			slog.Info("Default configuration saved", "file", "starthilfe.yml")
		}
	case "add":
		if err := gitops.AddSubtrees("starthilfe.yml"); err != nil {
			slog.Error("Error adding subtrees", "error", err)
		}
	case "update":
		if err := gitops.UpdateSubtrees(".starthilfe_state"); err != nil {
			slog.Error("Error updating subtrees", "error", err)
		}
	default:
		slog.Info("Unknown command", "command", os.Args[1])
	}
}
