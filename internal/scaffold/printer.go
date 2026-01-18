package scaffold

import (
	"fmt"
	"path/filepath"
)

type Printer interface {
	PrintNextSteps()
}

type printerImpl struct {
	opts ScaffoldOptions
}

func NewPrinter(opts ScaffoldOptions) Printer {
	return &printerImpl{opts}
}

func (p *printerImpl) PrintNextSteps() {
	fmt.Println()

	cOk.Printf("✓ ")
	cOk.Println("Project generated successfully")

	cTip.Print("⚑ ")
	cTip.Println("Review your environment variables before running")

	fmt.Println()
	cHeader.Println("Next steps")
	cBullet.Println("────────────────────────────────────────")

	projectDir := filepath.Base(filepath.Clean(p.opts.OutDir))

	printStep("Go to project directory", fmt.Sprintf("cd %s", projectDir))
	printStep("Setup environment", "cp .env.example .env.local")
	cNote.Println("    Configure database and app settings inside .env.local before running the project.")

	if p.opts.Preset == "full" {
		p.printFullPresetSteps()
	} else {
		p.printBasePresetSteps()
	}

	fmt.Println()
	cNote.Println("  • Server: http://localhost:8080")
	cNote.Println("  • Edit .env.local if config changes")
	fmt.Println()
}

func (p *printerImpl) printFullPresetSteps() {
	fmt.Println()
	cStepTitle.Println("  FULL preset detected")

	cNote.Println("    This preset uses Makefile and Air for development.")
	cNote.Println("    Make sure you have the following installed:")

	cBullet.Println("      - make")
	cBullet.Println("      - air (live reload)")
	cBullet.Println()

	printStep("Install Air (if not installed)", "go install github.com/air-verse/air@latest")
	printStep("Start dependencies", "make docker")
	printStep("Run migration", "make migrate")
	printStep("Run seed", "make seed")
	printStep("Run development server", "make dev")
}

func (p *printerImpl) printBasePresetSteps() {
	fmt.Println()
	cStepTitle.Println("  BASE preset detected")

	printStep("Run migration", "go run ./cmd/migrate")
	printStep("Run seed", "go run ./cmd/seed")
	printStep("Run the app", "go run ./cmd/api")
}
