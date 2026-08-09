// Package apkparser extracts app metadata from APK files.
//
// It replaces the previous `aapt` binary shell-out (aaptjs) with a pure-Go
// parser that reads AndroidManifest.xml and resolves resource references
// (the application label) from resources.arsc.
package apkparser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"github.com/avast/apkparser"

	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
)

// Parser implements application.ApkParser.
type Parser struct{}

// NewParser creates an APK metadata parser.
func NewParser() *Parser {
	return &Parser{}
}

// manifestDecoder collects the decoded XML tokens of AndroidManifest.xml.
type manifestDecoder struct {
	tokens []xml.Token
}

func (d *manifestDecoder) EncodeToken(t xml.Token) error {
	d.tokens = append(d.tokens, xml.CopyToken(t))
	return nil
}

func (d *manifestDecoder) Flush() error { return nil }

// Parse reads the APK at path and returns its metadata.
func (p *Parser) Parse(path string) (domainapp.ApkMetadata, error) {
	var md domainapp.ApkMetadata

	dec := &manifestDecoder{}
	zipErr, resErr, manErr := apkparser.ParseApk(path, dec)
	if zipErr != nil {
		return md, fmt.Errorf("open apk: %w", zipErr)
	}
	if manErr != nil {
		return md, fmt.Errorf("parse manifest: %w", manErr)
	}

	var inApplication bool
	for _, t := range dec.tokens {
		switch el := t.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "manifest":
				for _, a := range el.Attr {
					switch a.Name.Local {
					case "package":
						md.AppID = a.Value
					case "versionCode":
						md.Version = a.Value
					case "versionName":
						md.VersionName = a.Value
					}
				}
			case "application":
				inApplication = true
				for _, a := range el.Attr {
					if a.Name.Local == "label" {
						md.Name = a.Value
					}
				}
			}
		case xml.EndElement:
			if el.Name.Local == "application" {
				inApplication = false
			}
		case xml.CharData:
			// The application label may be emitted as character data of the
			// label element when it does not resolve to a plain string.
			if inApplication && md.Name == "" {
				if v := strings.TrimSpace(string(el)); v != "" {
					md.Name = v
				}
			}
		}
	}

	// resourcesErr is non-nil when resources.arsc could not be parsed; the
	// label then stays unresolved. Only fail when the label is missing.
	if md.Name == "" {
		if resErr != nil {
			return md, fmt.Errorf("parse resources: %w", resErr)
		}
		return md, errors.New("apk manifest has no application label")
	}
	if md.AppID == "" {
		return md, errors.New("apk manifest has no package id")
	}

	return md, nil
}
