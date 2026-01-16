package scaffold

import (
	"embed"
	"fmt"
	"os"
)

//go:embed templates/**
var templateFS embed.FS

type Options struct {
	ProjectName string
	ModulePath  string
	DB          string
	Preset      string // "base" | "full"
}

func info(format string, args ...any) {
	fmt.Printf("• "+format+"\n", args...)
}

func Generate(outDir string, opts Options) error {
	info("Creating project: %s", outDir)

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	// 1) Render preset
	presetRoot := "templates/base"
	if opts.Preset == "full" {
		presetRoot = "templates/full"
	}

	info("Rendering preset (%s)...", opts.Preset)
	if err := renderDir(presetRoot, outDir, opts); err != nil {
		return err
	}

	// 2) Render DB overlay (override files)
	dbRoot := "templates/db/postgres"
	if opts.DB == "mysql" {
		dbRoot = "templates/db/mysql"
	}

	info("Rendering database driver (%s)...", opts.DB)
	if err := renderDir(dbRoot, outDir, opts); err != nil {
		return err
	}

	info("Running: go mod tidy")
	if err := run(outDir, "go", "mod", "tidy"); err != nil {
		return err
	}

	info("Initializing git repository")
	_ = run(outDir, "git", "init")

	info("Done ✅")
	return nil
}
