package media_test

import (
	"bytes"
	"elcomware/app/actions"
	"elcomware/app/http/requests"
	"elcomware/app/models"
	"elcomware/app/services"
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"image"
	"image/png"
	"strings"
	"testing"
)

func picture(w, h int) string {
	var b bytes.Buffer
	_ = png.Encode(&b, image.NewRGBA(image.Rect(0, 0, w, h)))
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(b.Bytes())
}
func TestUploadReencodesAndChecksDimensions(t *testing.T) {
	p := services.DefaultMediaPolicy()
	result, err := services.ValidateMediaImage(picture(800, 600), p, false)
	require.NoError(t, err)
	require.Equal(t, 800, result.Width)
	require.Equal(t, "image/png", result.Mime)
	_, err = services.ValidateMediaImage(picture(80, 60), p, false)
	require.Error(t, err)
	_, err = services.ValidateMediaImage("data:image/svg+xml;base64,"+base64.StdEncoding.EncodeToString([]byte("<svg onload='alert(1)'/>")), p, false)
	require.Error(t, err)
}
func TestMockupsEscapeTextAndRejectArbitraryMarkup(t *testing.T) {
	opts := services.MockupOptions{Template: "software", Title: "<script>alert(1)</script>", Brand: "A & B", Primary: "#14283D", Accent: "#D5A64D", Background: "#F4F5F7"}
	out, err := services.RenderProductMockup(opts, services.DefaultMediaPolicy())
	require.NoError(t, err)
	require.NotContains(t, string(out), "<script>")
	require.Contains(t, string(out), "Illustrative mockup")
	opts.Primary = "url(https://evil.test)"
	_, err = services.RenderProductMockup(opts, services.DefaultMediaPolicy())
	require.Error(t, err)
	opts.Primary = "#14283D"
	opts.Template = "unknown"
	_, err = services.RenderProductMockup(opts, services.DefaultMediaPolicy())
	require.Error(t, err)
}
func TestPolicyBoundsAndURLSafety(t *testing.T) {
	p := services.DefaultMediaPolicy()
	require.NoError(t, services.ValidateMediaPolicy(p))
	p.MaxBytes = 1 << 30
	require.Error(t, services.ValidateMediaPolicy(p))
	require.NoError(t, services.ValidateProductImageURL("https://example.com/photo.jpg"))
	for _, v := range []string{"javascript:alert(1)", "data:image/png;base64,test", "//evil.test/x", "/../secret", "/public/../../secret"} {
		require.Error(t, services.ValidateProductImageURL(v), v)
	}
	require.Error(t, services.ValidateProductImageURL("https://user:secret@example.com/x"))
	require.Error(t, services.ValidateProductImageURL(strings.Repeat("x", 3000)))
}
func TestProductMediaAndValueFieldsPreserveOmittedValues(t *testing.T) {
	p := models.Product{ShortDesc: "Save time", GalleryUrls: []string{"/public/old.png"}, Featured: true}
	req := requests.CreateProductRequest{Title: "Valid product", Description: "Long description", Type: "physical", Format: "retail"}
	actions.ApplyProductInput(&p, &req)
	require.Equal(t, "Save time", p.ShortDesc)
	require.True(t, p.Featured)
	require.Equal(t, []string{"/public/old.png"}, p.GalleryUrls)
	value := "Make work easier"
	featured := false
	req.ShortDesc = &value
	req.Featured = &featured
	req.GalleryUrls = []string{"https://example.com/photo.png"}
	actions.ApplyProductInput(&p, &req)
	require.Equal(t, value, p.ShortDesc)
	require.False(t, p.Featured)
	require.Equal(t, req.GalleryUrls, p.GalleryUrls)
}
