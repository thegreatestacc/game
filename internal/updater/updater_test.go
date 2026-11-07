package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fixture struct {
	client    *Client
	data      []byte
	manifest  Manifest
	status    int
	corrupt   bool
	truncated bool
	downloads int
}

func setup(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{data: []byte("test game binary version one"), status: 200}
	h := sha256.Sum256(f.data)
	f.manifest = Manifest{Version: "v1.0.0", Asset: GameAsset, SHA256: hex.EncodeToString(h[:]), Size: int64(len(f.data))}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("public repository request unexpectedly included a token")
		}
		switch r.URL.Path {
		case "/repos/owner/game/releases/latest":
			w.WriteHeader(f.status)
			if f.status != 200 {
				return
			}
			json.NewEncoder(w).Encode(release{Tag: "v1.0.0", Assets: []asset{{ID: 1, Name: ManifestAsset}, {ID: 2, Name: GameAsset, Size: f.manifest.Size}}})
		case "/repos/owner/game/releases/assets/1":
			if r.Header.Get("Accept") != "application/octet-stream" {
				t.Error("not an asset download")
			}
			json.NewEncoder(w).Encode(f.manifest)
		case "/repos/owner/game/releases/assets/2":
			f.downloads++
			data := append([]byte(nil), f.data...)
			if f.corrupt {
				data[0] ^= 0xff
			}
			if f.truncated {
				data = data[:len(data)/2]
			}
			w.Write(data)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	f.client = New("owner/game", "", t.TempDir())
	f.client.apiBase = srv.URL
	return f
}

func TestInstallVerifiedReleaseAndReuse(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	a, err := f.client.Latest(ctx)
	if err != nil {
		t.Fatal(err)
	}
	path, err := f.client.Install(ctx, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !verifiedFile(path, f.manifest) {
		t.Fatal("installed file mismatch")
	}
	if _, err = f.client.Install(ctx, a, nil); err != nil {
		t.Fatal(err)
	}
	if f.downloads != 1 {
		t.Fatalf("unexpected downloads: %d", f.downloads)
	}
	// Even with a verified local game, a failed fresh check cannot authorize launch.
	f.status = 503
	if _, err = f.client.Latest(ctx); err == nil {
		t.Fatal("offline check accepted cached game")
	}
}

func TestRejectCorruptAndInterruptedDownloads(t *testing.T) {
	for _, kind := range []string{"hash", "truncated"} {
		t.Run(kind, func(t *testing.T) {
			f := setup(t)
			ctx := context.Background()
			a, err := f.client.Latest(ctx)
			if err != nil {
				t.Fatal(err)
			}
			old := filepath.Join(f.client.Directory, "previous-game.exe")
			os.WriteFile(old, []byte("old working version"), 0600)
			f.corrupt = kind == "hash"
			f.truncated = kind == "truncated"
			if _, err = f.client.Install(ctx, a, nil); err == nil {
				t.Fatal("invalid download installed")
			}
			path := filepath.Join(f.client.Directory, "versions", a.Manifest.SHA256, GameAsset)
			if _, err = os.Stat(path); !os.IsNotExist(err) {
				t.Fatal("unverified executable left behind")
			}
			files, _ := filepath.Glob(filepath.Join(filepath.Dir(path), "*.tmp"))
			if len(files) > 0 {
				t.Fatal("temporary download was not cleaned")
			}
			data, _ := os.ReadFile(old)
			if string(data) != "old working version" {
				t.Fatal("old installation damaged")
			}
		})
	}
}

func TestTamperedLocalFileIsRepaired(t *testing.T) {
	f := setup(t)
	ctx := context.Background()
	a, err := f.client.Latest(ctx)
	if err != nil {
		t.Fatal(err)
	}
	path, err := f.client.Install(ctx, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, []byte("tampered"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = f.client.Install(ctx, a, nil); err != nil {
		t.Fatal(err)
	}
	if !verifiedFile(path, f.manifest) || f.downloads != 2 {
		t.Fatal("tampered install accepted")
	}
}

func TestInvalidManifestAndAccessDenied(t *testing.T) {
	for _, status := range []int{401, 403, 404, 429, 500} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			f := setup(t)
			f.status = status
			if _, err := f.client.Latest(context.Background()); err == nil {
				t.Fatal("HTTP failure accepted")
			}
		})
	}
	f := setup(t)
	f.manifest.SHA256 = "../../escape"
	if _, err := f.client.Latest(context.Background()); err == nil {
		t.Fatal("unsafe digest accepted")
	}
	f.manifest.SHA256 = strings.Repeat("0", 64)
	f.manifest.Version = "v9.9.9"
	if _, err := f.client.Latest(context.Background()); err == nil {
		t.Fatal("mismatched tag accepted")
	}
}

func TestRedirectDoesNotLeakToken(t *testing.T) {
	c := New("owner/game", "secret", t.TempDir())
	origin, _ := http.NewRequest("GET", "https://api.github.com/repos/owner/game/releases/assets/2", nil)
	next, _ := http.NewRequest("GET", "https://release-assets.githubusercontent.com/signed", nil)
	next.Header.Set("Authorization", "Bearer secret")
	if err := c.HTTP.CheckRedirect(next, []*http.Request{origin}); err != nil {
		t.Fatal(err)
	}
	if next.Header.Get("Authorization") != "" {
		t.Fatal("token leaked to storage redirect")
	}
	next.URL = &url.URL{Scheme: "http", Host: "example.com"}
	if err := c.HTTP.CheckRedirect(next, []*http.Request{origin}); err == nil {
		t.Fatal("insecure redirect accepted")
	}
}

func TestCancelledCheck(t *testing.T) {
	f := setup(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := f.client.Latest(ctx); err == nil {
		t.Fatal("cancelled check accepted")
	}
}
