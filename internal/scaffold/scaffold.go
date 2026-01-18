package scaffold

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed all:templates/**
var templateFS embed.FS

type Scaffold interface {
	Generate() error
}

type ScaffoldOptions struct {
	ProjectName string
	ModulePath  string
	DB          string
	Preset      string
	OutDir      string
}

type scaffoldImpl struct {
	templateFS embed.FS
	opts       ScaffoldOptions
}

func NewScaffold(opts ScaffoldOptions) Scaffold {
	return &scaffoldImpl{
		templateFS: templateFS,
		opts:       opts,
	}
}


func (s *scaffoldImpl) Generate() error {
	presetRoot := "templates/base"
	if s.opts.Preset == "full" {
		presetRoot = "templates/full"
	}

	dbRoot := "templates/db/postgres"
	if s.opts.DB == "mysql" {
		dbRoot = "templates/db/mysql"
	}

	dstBase := filepath.Join("internal", "infra", "database")

	fmt.Println()
	header("Goscaff • Project Generator")
	info("Folder : " + s.opts.OutDir)
	info("Preset : " + s.opts.Preset)
	info("DB     : " + s.opts.DB)
	fmt.Println()

	action("Creating project directory")
	if err := os.MkdirAll(s.opts.OutDir, 0o755); err != nil {
		return err
	}
	success("Directory created")

	action("Rendering preset (" + s.opts.Preset + ")")
	if err := renderDir(presetRoot, s.opts.OutDir, s.opts); err != nil {
		return err
	}
	success("Preset rendered")

	// kalau kamu masih mau overlay selalu jalan, ya biarin.
	// Tapi kalau mau base bersih, taruh if opts.Preset == "full"
	action("Rendering database driver (" + s.opts.DB + ")")
	if err := renderDirTo(dbRoot, s.opts.OutDir, dstBase, s.opts); err != nil {
		return err
	}
	success("Database applied")

	action("Running: go mod tidy")
	if err := runVerbose(s.opts.OutDir, "go", "mod", "tidy"); err != nil {
		return err
	}
	success("Dependencies installed")

	action("Initializing git repository")
	_ = runQuiet(s.opts.OutDir, "git", "init")
	success("Git initialized")
	
	return nil
}
