package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"image"
	_ "image/jpeg"
	"image/png"
	"net/url"
	"regexp"
	"strings"
)

type MediaPolicy struct {
	Version      string `json:"version"`
	MaxBytes     int    `json:"max_bytes"`
	MinWidth     int    `json:"min_width"`
	MinHeight    int    `json:"min_height"`
	MaxDimension int    `json:"max_dimension"`
	AllowRemote  bool   `json:"allow_remote"`
	AllowMockups bool   `json:"allow_mockups"`
}

func DefaultMediaPolicy() MediaPolicy {
	return MediaPolicy{MaxBytes: 5 * 1024 * 1024, MinWidth: 600, MinHeight: 600, MaxDimension: 4096, AllowRemote: true, AllowMockups: true}
}
func ValidateMediaPolicy(p MediaPolicy) error {
	if p.MaxBytes < 1024*100 || p.MaxBytes > 10*1024*1024 || p.MinWidth < 100 || p.MinHeight < 100 || p.MinWidth > 2048 || p.MinHeight > 2048 || p.MaxDimension < p.MinWidth || p.MaxDimension < p.MinHeight || p.MaxDimension > 6000 {
		return errors.New("image policy must use 100 KB–10 MB, minimum dimensions 100–2048 px and maximum dimensions up to 6000 px")
	}
	return nil
}

type MediaImage struct {
	Bytes         []byte
	Width, Height int
	Mime          string
}

func ValidateMediaImage(raw string, p MediaPolicy, logo bool) (MediaImage, error) {
	var out MediaImage
	if err := ValidateMediaPolicy(p); err != nil {
		return out, err
	}
	comma := strings.IndexByte(raw, ',')
	if comma < 0 || comma > 30 || (!strings.HasPrefix(raw, "data:image/png;base64,") && !strings.HasPrefix(raw, "data:image/jpeg;base64,")) {
		return out, errors.New("upload a PNG or JPEG image")
	}
	if len(raw)-comma-1 > base64.StdEncoding.EncodedLen(p.MaxBytes) {
		return out, errors.New("image exceeds the upload size limit")
	}
	data, err := base64.StdEncoding.DecodeString(raw[comma+1:])
	if err != nil || len(data) > p.MaxBytes {
		return out, errors.New("image data is invalid or too large")
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "png" && format != "jpeg") {
		return out, errors.New("image is not a valid PNG or JPEG")
	}
	minW, minH := p.MinWidth, p.MinHeight
	if logo {
		minW, minH = 64, 64
	}
	if config.Width < minW || config.Height < minH || config.Width > p.MaxDimension || config.Height > p.MaxDimension || int64(config.Width)*int64(config.Height) > 24000000 {
		return out, fmt.Errorf("image must be at least %d × %d px and no larger than %d px on either side", minW, minH, p.MaxDimension)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return out, errors.New("image could not be decoded")
	}
	var normalized bytes.Buffer
	if err = png.Encode(&normalized, img); err != nil {
		return out, err
	}
	if normalized.Len() > 20*1024*1024 {
		return out, errors.New("normalized image is too large")
	}
	return MediaImage{normalized.Bytes(), config.Width, config.Height, "image/png"}, nil
}
func ValidateProductImageURL(raw string) error {
	if raw == "" {
		return nil
	}
	if len(raw) > 2048 || strings.ContainsAny(raw, "\r\n\\") {
		return errors.New("invalid image URL")
	}
	u, err := url.Parse(raw)
	if err != nil || u.User != nil {
		return errors.New("invalid image URL")
	}
	if u.Scheme == "https" && u.Host != "" {
		return nil
	}
	if u.Scheme == "" && u.Host == "" && strings.HasPrefix(raw, "/") && !strings.HasPrefix(raw, "//") && !strings.Contains(u.Path, "..") {
		return nil
	}
	return errors.New("image URL must use HTTPS or a local absolute path")
}

type MockupOptions struct {
	Template   string `json:"template"`
	Title      string `json:"title"`
	Brand      string `json:"brand"`
	Primary    string `json:"primary"`
	Accent     string `json:"accent"`
	Background string `json:"background"`
	LogoData   string `json:"logo_data,omitempty"`
}

var colorHex = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func RenderProductMockup(o MockupOptions, p MediaPolicy) ([]byte, error) {
	if !p.AllowMockups {
		return nil, errors.New("mockups are disabled by the media policy")
	}
	if len(o.Title) > 80 || len(o.Brand) > 40 {
		return nil, errors.New("use a title up to 80 characters and brand up to 40 characters")
	}
	if !colorHex.MatchString(o.Primary) || !colorHex.MatchString(o.Accent) || !colorHex.MatchString(o.Background) {
		return nil, errors.New("mockup colors must be six-digit hex colors")
	}
	templates := map[string]bool{"software": true, "box": true, "bottle": true, "apparel": true, "book": true, "equipment": true}
	if !templates[o.Template] {
		return nil, errors.New("choose a supported mockup template")
	}
	logo := ""
	if o.LogoData != "" {
		im, err := ValidateMediaImage(o.LogoData, p, true)
		if err != nil {
			return nil, err
		}
		logo = `<image x="480" y="260" width="240" height="110" preserveAspectRatio="xMidYMid meet" href="data:image/png;base64,` + base64.StdEncoding.EncodeToString(im.Bytes) + `"/>`
	}
	// Every shape and attribute is server-controlled. Text is escaped, and arbitrary SVG is never accepted.
	shape := ""
	switch o.Template {
	case "software":
		shape = `<rect x="205" y="160" width="790" height="470" rx="24" fill="PRIMARY"/><rect x="230" y="185" width="740" height="405" rx="12" fill="white"/><path d="M160 650H1040L990 690H210Z" fill="PRIMARY"/><rect x="255" y="210" width="145" height="355" rx="8" fill="ACCENT"/><rect x="425" y="395" width="515" height="140" rx="12" fill="BACKGROUND"/>`
	case "box":
		shape = `<path d="M350 190L780 150L870 230L440 280Z" fill="ACCENT"/><path d="M780 150L870 230V670L780 610Z" fill="PRIMARY"/><path d="M350 190L780 150V610L350 660Z" fill="white"/><path d="M350 190V660L440 710V280Z" fill="PRIMARY"/>`
	case "bottle":
		shape = `<rect x="490" y="120" width="220" height="100" rx="12" fill="PRIMARY"/><rect x="430" y="200" width="340" height="510" rx="95" fill="white"/><rect x="430" y="365" width="340" height="240" fill="ACCENT"/>`
	case "apparel":
		shape = `<path d="M420 160L300 215L180 390L325 470L385 360V710H815V360L875 470L1020 390L900 215L780 160L710 210H490Z" fill="PRIMARY"/><path d="M490 160Q600 265 710 160" fill="BACKGROUND"/><path d="M440 580H760" stroke="ACCENT" stroke-width="10"/>`
	case "book":
		shape = `<path d="M340 185L385 155H825V680H385L340 710Z" fill="PRIMARY"/><path d="M385 165H815V660H385Z" fill="white"/><path d="M385 660H825V690H385Q355 675 385 660Z" fill="ACCENT"/><path d="M385 165V660" stroke="ACCENT" stroke-width="12"/>`
	case "equipment":
		shape = `<rect x="290" y="220" width="620" height="420" rx="38" fill="PRIMARY"/><rect x="330" y="260" width="380" height="290" rx="14" fill="white"/><rect x="745" y="265" width="120" height="120" rx="10" fill="ACCENT"/><circle cx="780" cy="445" r="18" fill="white"/><circle cx="835" cy="445" r="18" fill="white"/><rect x="320" y="640" width="100" height="55" rx="10" fill="PRIMARY"/><rect x="780" y="640" width="100" height="55" rx="10" fill="PRIMARY"/>`
	}
	shape = strings.NewReplacer("PRIMARY", o.Primary, "ACCENT", o.Accent, "BACKGROUND", o.Background).Replace(shape)
	title := html.EscapeString(o.Title)
	brand := html.EscapeString(o.Brand)
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="900" viewBox="0 0 1200 900" role="img" aria-label="Illustrative product mockup"><rect width="1200" height="900" fill="%s"/><ellipse cx="600" cy="735" rx="360" ry="35" fill="%s" opacity="0.09"/>%s%s<rect x="400" y="390" width="400" height="120" rx="10" fill="white" opacity="0.95"/><text x="600" y="440" text-anchor="middle" font-family="sans-serif" font-size="28" font-weight="700" fill="%s">%s</text><text x="600" y="480" text-anchor="middle" font-family="sans-serif" font-size="16" fill="%s">%s</text><text x="600" y="822" text-anchor="middle" font-family="sans-serif" font-size="22" fill="%s">Illustrative mockup · Actual product may differ</text></svg>`, o.Background, o.Primary, shape, logo, o.Primary, brand, o.Primary, title, o.Primary)
	return []byte(svg), nil
}
