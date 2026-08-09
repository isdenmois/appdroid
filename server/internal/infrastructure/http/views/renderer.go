// Package views renders the two server-side HTML pages of the service.
package views

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"strings"

	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
)

//go:embed layout.html app_list.html app_page.html
var files embed.FS

// Item is a page row: an app plus the precomputed obtainium link.
//
// ObtainiumLink is typed template.URL because html/template would otherwise
// refuse to render the custom "obtainium://" scheme and output #ZgotmplZ.
// The link is generated from a fixed template, so it is safe to mark as URL.
type Item struct {
	domainapp.App
	ObtainiumLink template.URL
}

// Data is the payload shared by the page templates.
type Data struct {
	Title string
	Apps  []Item
}

// Renderer renders the SSR pages.
type Renderer struct {
	appList *template.Template
	appPage *template.Template
}

// NewRenderer parses the embedded page templates into two separate template
// sets (one per page) so their {{define "content"}} blocks never collide.
func NewRenderer() (*Renderer, error) {
	appList, err := template.New("layout.html").ParseFS(files, "layout.html", "app_list.html")
	if err != nil {
		return nil, err
	}

	appPage, err := template.New("layout.html").ParseFS(files, "layout.html", "app_page.html")
	if err != nil {
		return nil, err
	}

	return &Renderer{appList: appList, appPage: appPage}, nil
}

// AppList renders the /apps page listing all apps as Obtainium links. baseURL
// is the scheme://host prefix used inside the links.
func (r *Renderer) AppList(baseURL string, apps []domainapp.App) ([]byte, error) {
	items := make([]Item, 0, len(apps))
	for _, a := range apps {
		items = append(items, Item{App: a, ObtainiumLink: template.URL(obtainiumLink(baseURL, a))})
	}
	return render(r.appList, Data{Apps: items})
}

// AppPage renders the /app/:id page with a download link for the app.
func (r *Renderer) AppPage(apps []domainapp.App) ([]byte, error) {
	items := make([]Item, 0, len(apps))
	for _, a := range apps {
		items = append(items, Item{App: a})
	}
	return render(r.appPage, Data{Apps: items})
}

func render(tmpl *template.Template, d Data) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout.html", d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// additionalSettings mirrors the JSON the previous implementation embedded in
// every obtainium link.
type additionalSettings struct {
	VersionExtractionRegEx string `json:"versionExtractionRegEx"`
	MatchGroupToUse        string `json:"matchGroupToUse"`
	VersionDetection       bool   `json:"versionDetection"`
}

// obtainiumParams mirrors the params JSON payload of an obtainium link.
type obtainiumParams struct {
	ID                 string             `json:"id"`
	URL                string             `json:"url"`
	Author             string             `json:"author"`
	Name               string             `json:"name"`
	AdditionalSettings additionalSettings `json:"additionalSettings"`
}

// obtainiumLink builds an obtainium:// app link for the given app. The link is
// the JSON params payload percent-encoded the same way the previous
// implementation did it with encodeURI().
func obtainiumLink(baseURL string, app domainapp.App) string {
	params, err := json.Marshal(obtainiumParams{
		ID:     app.ID,
		URL:    baseURL + "/app/" + app.ID,
		Author: "Appdroid",
		Name:   app.Name,
		AdditionalSettings: additionalSettings{
			VersionExtractionRegEx: `-(.*)\.apk$`,
			MatchGroupToUse:        "1",
			VersionDetection:       true,
		},
	})
	if err != nil {
		return ""
	}

	return "obtainium://app/" + encodeURI(string(params))
}

// encodeURI encodes the string the way JavaScript encodeURI does: it escapes
// everything except the reserved characters !'()*-._~ and alphanumerics.
func encodeURI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9',
			c == '!', c == '\'', c == '(', c == ')', c == '*', c == '-',
			c == '.', c == '_', c == '~',
			c == ';', c == ',', c == '/', c == '?', c == ':', c == '@',
			c == '&', c == '=', c == '+', c == '$', c == '#':
			b.WriteByte(c)
		default:
			b.WriteString("%" + strings.ToUpper(hexByte(c)))
		}
	}
	return b.String()
}

// hexByte formats one byte as two uppercase hexadecimal digits.
func hexByte(b byte) string {
	const digits = "0123456789ABCDEF"
	return string([]byte{digits[b>>4], digits[b&0x0f]})
}
