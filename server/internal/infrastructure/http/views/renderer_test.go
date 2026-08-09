package views

import (
	"net/url"
	"strings"
	"testing"

	domainapp "github.com/isdenmois/appdroid/server/internal/domain/app"
)

func TestEncodeURI(t *testing.T) {
	// arrange
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"reserved pass through", `!'()*-._~;,/?:@&=+$#`, `!'()*-._~;,/?:@&=+$#`},
		{"alphanumerics pass through", "abcXYZ0129", "abcXYZ0129"},
		{"space is encoded", "a b", "a%20b"},
		{"utf-8 bytes are encoded", "привет", "%D0%BF%D1%80%D0%B8%D0%B2%D0%B5%D1%82"},
		{"backslash is encoded", `a\b`, "a%5Cb"},
		{"quotes and brackets are encoded", `"a[b]`, "%22a%5Bb%5D"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// act
			got := encodeURI(tc.in)

			// assert
			if got != tc.want {
				t.Errorf("encodeURI(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestEncodeURIMatchesJSEmbedding(t *testing.T) {
	// arrange
	// The link payload must survive being embedded into an HTML attribute and
	// decoded by a browser, so encodeURI output is a valid URL-escaped value.
	in := "https://example.com/app/a&b=c?d=1#frag"

	// act
	enc := encodeURI(in)
	dec, err := url.QueryUnescape(enc)

	// assert
	if err != nil {
		t.Fatalf("unescape: %v", err)
	}
	if dec != in {
		t.Errorf("round trip %q -> %q, want %q", in, enc, dec)
	}
}

func TestObtainiumLink(t *testing.T) {
	// arrange
	app := domainapp.App{
		ID:          "a1b2",
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "2.1",
		VersionName: "2.1",
		Type:        domainapp.TypeStatic,
		Apk:         "com.example.app.apk",
	}

	// act
	link := obtainiumLink("http://192.168.1.5:3999", app)

	// assert
	if !strings.HasPrefix(link, "obtainium://app/") {
		t.Fatalf("expected obtainium prefix, got %q", link)
	}

	// The payload is the JSON params object percent-encoded with encodeURI, so
	// it must decode back to a JSON document with the right fields.
	dec, err := url.QueryUnescape(strings.TrimPrefix(link, "obtainium://app/"))
	if err != nil {
		t.Fatalf("unescape payload: %v", err)
	}
	for _, want := range []string{
		`"id":"a1b2"`,
		`"url":"http://192.168.1.5:3999/app/a1b2"`,
		`"author":"Appdroid"`,
		`"name":"Example App"`,
		`"versionExtractionRegEx":"-(.*)\\.apk$"`,
		`"matchGroupToUse":"1"`,
		`"versionDetection":true`,
	} {
		if !strings.Contains(dec, want) {
			t.Errorf("payload missing %s, got %q", want, dec)
		}
	}
}

func TestRenderAppList(t *testing.T) {
	// arrange
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("renderer: %v", err)
	}
	app := domainapp.App{
		ID:          "a1b2",
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "2.1",
		VersionName: "2.1",
		Type:        domainapp.TypeStatic,
		Apk:         "com.example.app.apk",
	}

	// act
	out, err := r.AppList("https://appdroid.example.com", []domainapp.App{app})

	// assert
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	body := string(out)
	if !strings.Contains(body, "obtainium://app/") {
		t.Errorf("expected obtainium link in page, got %q", body)
	}
	if strings.Contains(body, "#ZgotmplZ") {
		t.Errorf("obtainium scheme must not be sanitized, got %q", body)
	}
	if !strings.Contains(body, "https://appdroid.example.com/app/a1b2") {
		t.Errorf("expected app url in link, got %q", body)
	}
}

func TestRenderAppPage(t *testing.T) {
	// arrange
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("renderer: %v", err)
	}
	app := domainapp.App{
		ID:          "a1b2",
		AppID:       "com.example.app",
		Name:        "Example App",
		Version:     "2.1",
		VersionName: "2.1",
		Type:        domainapp.TypeStatic,
		Apk:         "com.example.app.apk",
	}

	// act
	out, err := r.AppPage([]domainapp.App{app})

	// assert
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	body := string(out)
	if !strings.Contains(body, "/apk/com.example.app.apk/a1b2-2.1.apk") {
		t.Errorf("expected download link in page, got %q", body)
	}
}
