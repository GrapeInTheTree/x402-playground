package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/GrapeInTheTree/x402-demo/internal/config"
	"github.com/GrapeInTheTree/x402-demo/internal/tui"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/dashboard"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/explore"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/home"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/learn"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/practice"
)

func main() {
	mode := flag.String("mode", "", "Start directly in a mode: learn, explore, practice, dashboard")
	flow := flag.String("flow", "", "For practice mode: eip3009, permit2, sidebyside")
	flag.Parse()

	_ = godotenv.Load()

	// Load config (best-effort — Learn/Explore modes work without private keys)
	cfg, _ := config.LoadExplorer()

	flowFlag := *flow
	factories := map[tui.Page]tui.SubModelFactory{
		tui.PageHome: func(w, h int) tui.SubModel {
			return home.New(w, h)
		},
		tui.PageLearn: func(w, h int) tui.SubModel {
			return learn.New(w, h)
		},
		tui.PageExplore: func(w, h int) tui.SubModel {
			return explore.New(w, h)
		},
		tui.PagePractice: func(w, h int) tui.SubModel {
			if flowFlag != "" {
				return practice.NewWithFlow(w, h, cfg, flowFlag)
			}
			return practice.New(w, h, cfg)
		},
		tui.PageDashboard: func(w, h int) tui.SubModel {
			return dashboard.New(w, h, cfg)
		},
	}

	// Resolve --mode flag to starting page
	startPage := tui.PageHome
	switch *mode {
	case "learn":
		startPage = tui.PageLearn
	case "explore":
		startPage = tui.PageExplore
	case "practice":
		startPage = tui.PagePractice
	case "dashboard":
		startPage = tui.PageDashboard
	case "":
		// default: home
	default:
		fmt.Fprintf(os.Stderr, "Unknown mode: %q (valid: learn, explore, practice, dashboard)\n", *mode)
		os.Exit(1)
	}

	// Validate --flow flag
	if *flow != "" && *mode != "practice" {
		fmt.Fprintf(os.Stderr, "--flow requires --mode=practice\n")
		os.Exit(1)
	}
	if *flow != "" {
		switch *flow {
		case "eip3009", "permit2", "sidebyside":
			// valid
		default:
			fmt.Fprintf(os.Stderr, "Unknown flow: %q (valid: eip3009, permit2, sidebyside)\n", *flow)
			os.Exit(1)
		}
	}

	model := tui.NewRootModelWithStart(factories, startPage)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
