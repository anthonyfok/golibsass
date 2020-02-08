// Copyright © 2020 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package libsass a SCSS transpiler to CSS using LibSASS.
package libsass

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bep/golibsass/internal/libsass"
)

type libsassTranspiler struct {
	options Options
}

// New creates a new libsass transpiler configured with the given options.
func New(options Options) (Transpiler, error) {
	return libsassTranspiler{options: options}, nil
}

// Execute transpiles the SCSS or SASS from src into dst.
func (t libsassTranspiler) Execute(src string) (Result, error) {
	var result Result

	if t.options.SassSyntax {
		// LibSass does not support this directly, so have to handle the main SASS content
		// special.
		var err error
		src, err = libsass.SassToScss(src)
		if err != nil {
			return result, err
		}
	}

	dataCtx := libsass.SassMakeDataContext(src)

	opts := libsass.SassDataContextGetOptions(dataCtx)
	{
		// Set options

		if t.options.ImportResolver != nil {
			idx := libsass.AddImportResolver(opts, t.options.ImportResolver)
			defer libsass.DeleteImportResolver(idx)
		}

		if t.options.Precision != 0 {
			libsass.SassOptionSetPrecision(opts, t.options.Precision)
		}

		if t.options.SourceMapFilename != "" {
			libsass.SassOptionSetSourceMapFile(opts, t.options.SourceMapFilename)
		}

		if t.options.SourceMapRoot != "" {
			libsass.SassOptionSetSourceMapRoot(opts, t.options.SourceMapRoot)
		}

		if t.options.OutputPath != "" {
			libsass.SassOptionSetOutputPath(opts, t.options.OutputPath)
		}
		if t.options.InputPath != "" {
			libsass.SassOptionSetInputPath(opts, t.options.InputPath)
		}

		libsass.SassOptionSetSourceMapContents(opts, t.options.SourceMapContents)
		libsass.SassOptionSetOmitSourceMapURL(opts, t.options.OmitSourceMapURL)
		libsass.SassOptionSetSourceMapEmbed(opts, t.options.EnableEmbeddedSourceMap)
		libsass.SassOptionSetIncludePath(opts, strings.Join(t.options.IncludePaths, string(os.PathListSeparator)))
		libsass.SassOptionSetOutputStyle(opts, int(t.options.OutputStyle))
		libsass.SassOptionSetSourceComments(opts, false)
		libsass.SassDataContextSetOptions(dataCtx, opts)
	}

	ctx := libsass.SassDataContextGetContext(dataCtx)
	compiler := libsass.SassMakeDataCompiler(dataCtx)
	defer libsass.SassDeleteCompiler(compiler)

	libsass.SassCompilerParse(compiler)
	libsass.SassCompilerExecute(compiler)

	result.CSS = libsass.SassContextGetOutputString(ctx)

	if status := libsass.SassContextGetErrorStatus(ctx); status != 0 {
		return result, jsonToError(libsass.SassContextGetErrorJSON(ctx))
	}

	result.SourceMapFilename = libsass.SassOptionGetSourceMapFile(opts)
	result.SourceMapContent = libsass.SassContextGetSourceMapString(ctx)

	return result, nil
}

type Result struct {
	CSS string

	// If source maps are configured.
	SourceMapFilename string
	SourceMapContent  string
}

type Transpiler interface {
	Execute(src string) (Result, error)
}

type (
	OutputStyle int
)

const (
	NestedStyle OutputStyle = iota
	ExpandedStyle
	CompactStyle
	CompressedStyle
)

func ParseOutputStyle(s string) OutputStyle {
	switch strings.ToLower(s) {
	case "nested":
		return NestedStyle
	case "expanded":
		return ExpandedStyle
	case "compact":
		return CompactStyle
	case "compressed":
		return CompressedStyle
	}
	return NestedStyle
}

type Options struct {
	// Default is nested.
	OutputStyle OutputStyle

	// Precision of floating point math.
	Precision int

	// File paths to use to resolve imports.
	IncludePaths []string

	// ImportResolver can be used to supply a custom import resolver, both to redirect
	// to another URL or to return the body.
	ImportResolver func(url string, prev string) (newURL string, body string, resolved bool)

	// Used to indicate "old style" SASS for the input stream.
	SassSyntax bool

	// Source map settings
	SourceMapFilename       string
	SourceMapRoot           string
	InputPath               string
	OutputPath              string
	SourceMapContents       bool
	OmitSourceMapURL        bool
	EnableEmbeddedSourceMap bool
}

func jsonToError(jsonstr string) (e Error) {
	if err := json.Unmarshal([]byte(jsonstr), &e); err != nil {
		e.Message = "unknown error"
	}
	return
}

type Error struct {
	Status  int    `json:"status"`
	Column  int    `json:"column"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("file %q, line %d, col %d: %s ", e.File, e.Line, e.Column, e.Message)
}
