// Package updater checks GitHub Releases and installs verified, versioned game
// executables. It has no dependency on the graphical engine.
package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const GameAsset = "SpellGame-game.exe"
const ManifestAsset = "release-manifest.json"
const maxGameSize int64 = 512 << 20

type Manifest struct {
	Version string `json:"version"`
	Asset   string `json:"asset"`
	SHA256  string `json:"sha256"`
	Size    int64  `json:"size"`
}

type asset struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type release struct {
	Tag        string  `json:"tag_name"`
	Draft      bool    `json:"draft"`
	Prerelease bool    `json:"prerelease"`
	Assets     []asset `json:"assets"`
}

type Available struct {
	Manifest Manifest
	assetID  int64
}

type Progress func(message string, fraction float64)

type Client struct {
	Repository string
	// Token is optional; public releases are fetched anonymously.
	Token     string
	Directory string
	HTTP      *http.Client
	apiBase   string
}

func New(repository, token, directory string) *Client {
	return &Client{Repository: repository, Token: strings.TrimSpace(token), Directory: directory,
		apiBase: "https://api.github.com", HTTP: &http.Client{Timeout: 5 * time.Minute, CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("слишком много перенаправлений")
			}
			if req.URL.Scheme != "https" {
				return errors.New("небезопасное перенаправление загрузки")
			}
			// A signed GitHub storage URL does not need our private repo token.
			if len(via) > 0 && req.URL.Host != via[0].URL.Host {
				req.Header.Del("Authorization")
			}
			return nil
		}}}
}

var repoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`)

func (c *Client) request(ctx context.Context, path, accept string) (*http.Response, error) {
	if !repoPattern.MatchString(c.Repository) {
		return nil, errors.New("неверный адрес репозитория GitHub")
	}
	u := c.apiBase + "/repos/" + c.Repository + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.New("не удалось создать запрос к GitHub")
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("User-Agent", "SpellGame-Launcher")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		// Do not surface signed redirect URLs or credentials in a UI error.
		if ctx.Err() != nil {
			return nil, errors.New("проверка отменена или превышено время ожидания")
		}
		return nil, errors.New("нет соединения с GitHub; проверьте интернет и повторите попытку")
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		switch resp.StatusCode {
		case 401:
			return nil, errors.New("GitHub отклонил запрос к релизу")
		case 403, 429:
			return nil, errors.New("GitHub ограничил частоту запросов; подождите и повторите попытку")
		case 404:
			return nil, errors.New("релиз не найден; проверьте публикацию первого релиза GitHub Actions")
		default:
			return nil, fmt.Errorf("GitHub ответил HTTP %d; повторите позже", resp.StatusCode)
		}
	}
	return resp, nil
}

func (c *Client) read(ctx context.Context, path, accept string, limit int64) ([]byte, error) {
	resp, err := c.request(ctx, path, accept)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, errors.New("не удалось прочитать ответ GitHub")
	}
	if int64(len(data)) > limit {
		return nil, errors.New("ответ GitHub превышает допустимый размер")
	}
	return data, nil
}

func assetPath(id int64) string { return fmt.Sprintf("/releases/assets/%d", id) }

// Latest always performs an online check. A cached manifest never authorizes
// game launch when the server is unreachable.
func (c *Client) Latest(ctx context.Context) (Available, error) {
	var result Available
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	data, err := c.read(ctx, "/releases/latest", "application/vnd.github+json", 2<<20)
	if err != nil {
		return result, err
	}
	var r release
	if json.Unmarshal(data, &r) != nil || r.Tag == "" || r.Draft || r.Prerelease {
		return result, errors.New("GitHub не вернул опубликованный стабильный релиз")
	}
	var manifestAsset, gameAsset asset
	for _, a := range r.Assets {
		if a.Name == ManifestAsset {
			manifestAsset = a
		}
		if a.Name == GameAsset {
			gameAsset = a
		}
	}
	if manifestAsset.ID <= 0 || gameAsset.ID <= 0 {
		return result, errors.New("в релизе отсутствует игра или release-manifest.json; дождитесь завершения сборки")
	}
	data, err = c.read(ctx, assetPath(manifestAsset.ID), "application/octet-stream", 64<<10)
	if err != nil {
		return result, err
	}
	var m Manifest
	if json.Unmarshal(data, &m) != nil {
		return result, errors.New("повреждено описание версии игры")
	}
	digest, err := hex.DecodeString(m.SHA256)
	if err != nil || len(digest) != sha256.Size || m.Version != r.Tag || m.Asset != GameAsset || m.Size <= 0 || m.Size > maxGameSize || gameAsset.Size != m.Size {
		return result, errors.New("описание релиза не прошло проверку версии, размера или SHA-256")
	}
	m.SHA256 = strings.ToLower(m.SHA256)
	return Available{Manifest: m, assetID: gameAsset.ID}, nil
}

func verifiedFile(path string, m Manifest) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() != m.Size {
		return false
	}
	h := sha256.New()
	n, err := io.Copy(h, io.LimitReader(f, m.Size+1))
	return err == nil && n == m.Size && hex.EncodeToString(h.Sum(nil)) == m.SHA256
}

// Install stages a download in the destination directory. Its hash must match
// before rename, and each digest gets a separate directory, so an update never
// overwrites the currently running or previously verified executable.
func (c *Client) Install(ctx context.Context, a Available, report Progress) (string, error) {
	if report == nil {
		report = func(string, float64) {}
	}
	m := a.Manifest
	// Available values normally come from Latest; validate again before using
	// a digest as a directory name, including calls from other Go code.
	digest, err := hex.DecodeString(m.SHA256)
	if err != nil || len(digest) != 32 || m.Size <= 0 || m.Size > maxGameSize || m.Asset != GameAsset || a.assetID <= 0 {
		return "", errors.New("некорректные данные обновления")
	}
	dir := filepath.Join(c.Directory, "versions", strings.ToLower(m.SHA256))
	path := filepath.Join(dir, GameAsset)
	report("Проверка установленных файлов…", 0)
	if verifiedFile(path, m) {
		report("Установлена последняя версия: "+m.Version, 1)
		return path, nil
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", errors.New("не удалось создать папку игры; проверьте права доступа")
	}
	resp, err := c.request(ctx, assetPath(a.assetID), "application/octet-stream")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.ContentLength >= 0 && resp.ContentLength != m.Size {
		return "", errors.New("размер загрузки не соответствует релизу")
	}
	f, err := os.CreateTemp(dir, "download-*.tmp")
	if err != nil {
		return "", errors.New("не удалось создать временный файл обновления")
	}
	temp := f.Name()
	defer os.Remove(temp)
	defer f.Close()
	h := sha256.New()
	reader := &progressReader{reader: io.LimitReader(resp.Body, m.Size+1), total: m.Size, report: report}
	n, err := io.Copy(io.MultiWriter(f, h), reader)
	if err != nil {
		return "", errors.New("загрузка прервана; повторите обновление")
	}
	if n != m.Size || hex.EncodeToString(h.Sum(nil)) != m.SHA256 {
		return "", errors.New("файл не прошёл проверку SHA-256; игра не будет запущена")
	}
	if err = f.Sync(); err != nil {
		return "", errors.New("не удалось сохранить обновление на диск")
	}
	if err = f.Close(); err != nil {
		return "", errors.New("не удалось закрыть файл обновления")
	}
	if err = os.Rename(temp, path); err != nil {
		// Another launcher may have installed the same valid release meanwhile.
		if verifiedFile(path, m) {
			return path, nil
		}
		return "", errors.New("не удалось установить обновление; закройте игру и повторите попытку")
	}
	report("Установлена последняя версия: "+m.Version, 1)
	return path, nil
}

type progressReader struct {
	reader      io.Reader
	done, total int64
	report      Progress
	last        time.Time
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.done += int64(n)
	if time.Since(p.last) > 100*time.Millisecond || err == io.EOF {
		p.last = time.Now()
		p.report(fmt.Sprintf("Загрузка игры: %.1f / %.1f МБ", float64(p.done)/(1<<20), float64(p.total)/(1<<20)), min(1, float64(p.done)/float64(p.total)))
	}
	return n, err
}
