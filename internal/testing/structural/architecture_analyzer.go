package structural

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Layer represents a clean architecture layer
type Layer string

const (
	LayerDomain         Layer = "domain"
	LayerUseCase        Layer = "usecase"
	LayerInfrastructure Layer = "infrastructure"
	LayerDelivery       Layer = "delivery"
)

// DependencyRule defines allowed dependencies between layers
type DependencyRule struct {
	From    Layer
	Allowed []Layer
}

// DefaultDependencyRules returns the standard clean architecture dependency rules
func DefaultDependencyRules() []DependencyRule {
	return []DependencyRule{
		{From: LayerDomain, Allowed: []Layer{}}, // Domain depends on nothing
		{From: LayerUseCase, Allowed: []Layer{LayerDomain}},
		{From: LayerInfrastructure, Allowed: []Layer{LayerDomain}},
		{From: LayerDelivery, Allowed: []Layer{LayerDomain, LayerUseCase}},
	}
}

// ImportViolation represents a violation of dependency rules
type ImportViolation struct {
	File          string
	Layer         Layer
	ImportedLayer Layer
	ImportPath    string
}

// ArchitectureAnalyzer analyzes Go code for architecture compliance
type ArchitectureAnalyzer struct {
	rootPath string
	rules    []DependencyRule
}

// NewArchitectureAnalyzer creates a new analyzer
func NewArchitectureAnalyzer(rootPath string) *ArchitectureAnalyzer {
	return &ArchitectureAnalyzer{
		rootPath: rootPath,
		rules:    DefaultDependencyRules(),
	}
}

// CheckLayerDependencies verifies that layer dependencies follow the rules
func (a *ArchitectureAnalyzer) CheckLayerDependencies() ([]ImportViolation, error) {
	violations := []ImportViolation{}

	err := filepath.Walk(filepath.Join(a.rootPath, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			fileViolations, err := a.checkFile(path)
			if err != nil {
				return err
			}
			violations = append(violations, fileViolations...)
		}

		return nil
	})

	return violations, err
}

// checkFile checks a single file for dependency violations
func (a *ArchitectureAnalyzer) checkFile(filePath string) ([]ImportViolation, error) {
	violations := []ImportViolation{}

	// Determine the layer of this file
	layer := a.determineLayer(filePath)
	if layer == "" {
		return violations, nil // Not in a recognized layer
	}

	// Parse the file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	// Check each import
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		importedLayer := a.determineLayerFromImport(importPath)

		if importedLayer != "" && !a.isAllowedDependency(layer, importedLayer) {
			violations = append(violations, ImportViolation{
				File:          filePath,
				Layer:         layer,
				ImportedLayer: importedLayer,
				ImportPath:    importPath,
			})
		}
	}

	return violations, nil
}

// determineLayer determines which layer a file belongs to
func (a *ArchitectureAnalyzer) determineLayer(filePath string) Layer {
	relPath := strings.TrimPrefix(filePath, a.rootPath)
	relPath = strings.TrimPrefix(relPath, "/")

	if strings.Contains(relPath, "internal/domain/") {
		return LayerDomain
	}
	if strings.Contains(relPath, "internal/usecase/") {
		return LayerUseCase
	}
	if strings.Contains(relPath, "internal/infrastructure/") {
		return LayerInfrastructure
	}
	if strings.Contains(relPath, "internal/delivery/") {
		return LayerDelivery
	}

	return ""
}

// determineLayerFromImport determines which layer an import path belongs to
func (a *ArchitectureAnalyzer) determineLayerFromImport(importPath string) Layer {
	if strings.Contains(importPath, "/internal/domain/") {
		return LayerDomain
	}
	if strings.Contains(importPath, "/internal/usecase/") {
		return LayerUseCase
	}
	if strings.Contains(importPath, "/internal/infrastructure/") {
		return LayerInfrastructure
	}
	if strings.Contains(importPath, "/internal/delivery/") {
		return LayerDelivery
	}

	return ""
}

// isAllowedDependency checks if a dependency is allowed by the rules
func (a *ArchitectureAnalyzer) isAllowedDependency(from, to Layer) bool {
	for _, rule := range a.rules {
		if rule.From == from {
			for _, allowed := range rule.Allowed {
				if allowed == to {
					return true
				}
			}
			return false
		}
	}
	return true // No rule found, allow by default
}

// CheckResourcePackageImports verifies no files import the resource package
func (a *ArchitectureAnalyzer) CheckResourcePackageImports() ([]string, error) {
	filesWithResourceImport := []string{}

	err := filepath.Walk(a.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			hasImport, err := a.fileImportsPackage(path, "internal/app/resource")
			if err != nil {
				return err
			}
			if hasImport {
				filesWithResourceImport = append(filesWithResourceImport, path)
			}
		}

		return nil
	})

	return filesWithResourceImport, err
}

// fileImportsPackage checks if a file imports a specific package
func (a *ArchitectureAnalyzer) fileImportsPackage(filePath, packagePath string) (bool, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return false, err
	}

	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		if strings.Contains(importPath, packagePath) {
			return true, nil
		}
	}

	return false, nil
}

// StructInfo represents information about a struct
type StructInfo struct {
	Name        string
	File        string
	HasConstructor bool
	Fields      []FieldInfo
}

// FieldInfo represents information about a struct field
type FieldInfo struct {
	Name string
	Type string
}

// CheckDependencyInjectionPattern verifies structs follow DI pattern
func (a *ArchitectureAnalyzer) CheckDependencyInjectionPattern() ([]StructInfo, error) {
	structsWithoutConstructors := []StructInfo{}

	err := filepath.Walk(filepath.Join(a.rootPath, "internal"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			structs, err := a.analyzeStructsInFile(path)
			if err != nil {
				return err
			}

			for _, s := range structs {
				if len(s.Fields) > 0 && !s.HasConstructor {
					structsWithoutConstructors = append(structsWithoutConstructors, s)
				}
			}
		}

		return nil
	})

	return structsWithoutConstructors, err
}

// analyzeStructsInFile analyzes structs in a file
func (a *ArchitectureAnalyzer) analyzeStructsInFile(filePath string) ([]StructInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structs := []StructInfo{}
	constructors := make(map[string]bool)

	// First pass: find all constructors
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if strings.HasPrefix(fn.Name.Name, "New") {
				// Extract struct name from constructor (e.g., NewAuthService -> AuthService)
				structName := strings.TrimPrefix(fn.Name.Name, "New")
				constructors[structName] = true
				constructors[strings.ToLower(structName[:1])+structName[1:]] = true // camelCase version
			}
		}
		return true
	})

	// Second pass: find all structs
	ast.Inspect(node, func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				structInfo := StructInfo{
					Name:        ts.Name.Name,
					File:        filePath,
					HasConstructor: constructors[ts.Name.Name],
					Fields:      []FieldInfo{},
				}

				// Analyze fields
				if st.Fields != nil {
					for _, field := range st.Fields.List {
						if len(field.Names) > 0 {
							fieldType := fmt.Sprintf("%v", field.Type)
							structInfo.Fields = append(structInfo.Fields, FieldInfo{
								Name: field.Names[0].Name,
								Type: fieldType,
							})
						}
					}
				}

				structs = append(structs, structInfo)
			}
		}
		return true
	})

	return structs, nil
}

// VerifyFolderStructure checks if required directories exist
func (a *ArchitectureAnalyzer) VerifyFolderStructure() ([]string, error) {
	requiredDirs := []string{
		"internal/domain",
		"internal/usecase",
		"internal/infrastructure",
		"internal/delivery",
	}

	missingDirs := []string{}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(a.rootPath, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			missingDirs = append(missingDirs, dir)
		}
	}

	return missingDirs, nil
}
