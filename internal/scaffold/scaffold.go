package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/base/*
var templateFS embed.FS

type Options struct {
	ProjectName string
	ModulePath  string
	DB          string
}
func info(format string, args ...any) {
	fmt.Printf("• "+format+"\n", args...)
}
func Generate(outDir string, opts Options) error {
	info("Creating project: %s", outDir)
	info("Rendering templates...")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	if err := renderDir("templates/base", outDir, opts); err != nil {
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
func renderDir(srcRoot, outRoot string, opts Options) error {
	entries, err := templateFS.ReadDir(srcRoot)
	if err != nil {
		return err
	}

	for _, e := range entries {
		srcPath := filepath.ToSlash(filepath.Join(srcRoot, e.Name()))
		outPath := filepath.Join(outRoot, e.Name())

		if e.IsDir() {
			if err := os.MkdirAll(outPath, 0o755); err != nil {
				return err
			}
			if err := renderDir(srcPath, outPath, opts); err != nil {
				return err
			}
			continue
		}
 
		if strings.HasSuffix(e.Name(), ".tmpl") {
			outPath = strings.TrimSuffix(outPath, ".tmpl")
			b, err := templateFS.ReadFile(srcPath)
			if err != nil {
				return err
			}
			t, err := template.New(e.Name()).Parse(string(b))
			if err != nil {
				return fmt.Errorf("parse template %s: %w", srcPath, err)
			}
			var buf bytes.Buffer
			if err := t.Execute(&buf, map[string]any{
				"PROJECT_NAME": opts.ProjectName,
				"MODULE_PATH":  opts.ModulePath,
				"DB":           opts.DB,
			}); err != nil {
				return fmt.Errorf("execute template %s: %w", srcPath, err)
			}
			if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
				return err
			}
			continue
		}
 
		b, err := templateFS.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(outPath, b, 0o644); err != nil {
			return err
		}
	}

	return nil
}

func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
