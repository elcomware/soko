package controllers

import (
	"elcomware/app/facades"
	appmiddleware "elcomware/app/http/middleware"
	"elcomware/app/models"
	"elcomware/app/services"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	inertia "github.com/goravel/inertia/facades"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MediaController struct{}

func NewMediaController() *MediaController { return &MediaController{} }
func mediaActor(ctx http.Context) *models.User {
	u := appmiddleware.CurrentUser(ctx)
	if u == nil || !u.Can("media.create", u.ID) {
		return nil
	}
	return u
}
func mediaError(ctx http.Context, status int, msg string) http.Response {
	return ctx.Response().Json(status, map[string]any{"error": msg})
}
func mediaDecode(ctx http.Context, out any) error {
	const limit = 15 * 1024 * 1024
	req := ctx.Request().Origin()
	if req.ContentLength > limit {
		return errors.New("request is too large")
	}
	d := json.NewDecoder(io.LimitReader(req.Body, limit+1))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return errors.New("invalid media request")
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return errors.New("invalid trailing request data")
	}
	return nil
}
func mediaWorkspaceActor(ctx http.Context) *models.User {
	u := appmiddleware.CurrentUser(ctx)
	if u == nil || (!u.Can("media.create", u.ID) && !u.Can("media.settings.manage", "")) {
		return nil
	}
	return u
}
func (c *MediaController) Workspace(ctx http.Context) http.Response {
	u := mediaWorkspaceActor(ctx)
	if u == nil {
		return ctx.Response().String(403, "Media access is not permitted.")
	}
	return inertia.Inertia().Render(ctx, "Dashboard/Media", map[string]any{"currentUser": u})
}
func (c *MediaController) Index(ctx http.Context) http.Response {
	u := mediaWorkspaceActor(ctx)
	if u == nil {
		return mediaError(ctx, 403, "Media access is not permitted.")
	}
	page, _ := strconv.Atoi(ctx.Request().Query("page", "1"))
	page = max(1, min(page, 10000))
	const perPage = 24
	assets := []models.ProductMedia{}
	var total int64
	canCreate := u.Can("media.create", u.ID)
	if canCreate {
		q := facades.Orm().Query().Model(&models.ProductMedia{}).Where("owner_id = ?", u.ID)
		search := strings.TrimSpace(ctx.Request().Query("q", ""))
		if len(search) > 100 {
			return mediaError(ctx, 422, "Search must be 100 characters or fewer.")
		}
		if search != "" {
			q = q.Where("(id LIKE ? OR options_json LIKE ?)", "%"+search+"%", "%"+search+"%")
		}
		kind := ctx.Request().Query("kind", "")
		if kind != "" && kind != "mockup" && kind != "upload" {
			return mediaError(ctx, 422, "Choose photos or mockups.")
		}
		if kind != "" {
			q = q.Where("kind = ?", kind)
		}
		var err error
		total, err = q.Count()
		if err != nil {
			return mediaError(ctx, 503, "Media storage is not ready.")
		}
		page = min(page, max(1, int((total+perPage-1)/perPage)))
		if err = q.OrderByDesc("created_at").OrderByDesc("id").Offset((page - 1) * perPage).Limit(perPage).Get(&assets); err != nil {
			return mediaError(ctx, 503, "Media storage is not ready.")
		}
	}
	return ctx.Response().Json(200, map[string]any{"assets": assets, "page": page, "per_page": perPage, "total": total, "policy": services.LoadMediaPolicy(), "can_create": canCreate, "can_manage_policy": u.Can("media.settings.manage", "")})
}

// Recipe returns only editable settings from the caller's own mockup. It never
// exposes embedded raster data or modifies the original published asset.
func (c *MediaController) Recipe(ctx http.Context) http.Response {
	u := mediaActor(ctx)
	if u == nil {
		return mediaError(ctx, 403, "Mockup access is not permitted.")
	}
	var asset models.ProductMedia
	if err := facades.Orm().Query().Where("id = ? AND owner_id = ? AND kind = ?", ctx.Request().Route("id"), u.ID, "mockup").First(&asset); err != nil || asset.ID == "" {
		return mediaError(ctx, 404, "Mockup not found.")
	}
	var options services.MockupOptions
	if err := json.Unmarshal([]byte(asset.OptionsJSON), &options); err != nil {
		return mediaError(ctx, 422, "This mockup has no reusable settings.")
	}
	return ctx.Response().Header("Cache-Control", "private, no-store").Json(200, map[string]any{"recipe": map[string]string{"template": options.Template, "title": options.Title, "brand": options.Brand, "primary": options.Primary, "accent": options.Accent, "background": options.Background}})
}
func (c *MediaController) Upload(ctx http.Context) http.Response {
	u := mediaActor(ctx)
	if u == nil {
		return mediaError(ctx, 403, "Media upload is not permitted.")
	}
	var req struct {
		Data string `json:"data"`
	}
	if err := mediaDecode(ctx, &req); err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	im, err := services.ValidateMediaImage(req.Data, services.LoadMediaPolicy(), false)
	if err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	asset, err := saveMedia(u.ID, "upload", "png", im.Bytes, im.Width, im.Height, "")
	if err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	return ctx.Response().Json(201, map[string]any{"asset": asset})
}
func (c *MediaController) Mockup(ctx http.Context) http.Response {
	u := mediaActor(ctx)
	if u == nil {
		return mediaError(ctx, 403, "Mockup creation is not permitted.")
	}
	var req services.MockupOptions
	if err := mediaDecode(ctx, &req); err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	b, err := services.RenderProductMockup(req, services.LoadMediaPolicy())
	if err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	req.LogoData = ""
	options, _ := json.Marshal(req)
	asset, err := saveMedia(u.ID, "mockup", "svg", b, 1200, 900, string(options))
	if err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	return ctx.Response().Json(201, map[string]any{"asset": asset})
}
func saveMedia(owner, kind, ext string, data []byte, w, h int, options string) (*models.ProductMedia, error) {
	count, err := facades.Orm().Query().Model(&models.ProductMedia{}).Where("owner_id = ?", owner).Count()
	if err != nil {
		return nil, errors.New("media storage is not ready")
	}
	if count >= 500 {
		return nil, errors.New("media library limit reached; contact the administrator")
	}
	id := uuid.NewString()
	dir := filepath.Join("storage", "app", "product-media")
	if err = os.MkdirAll(dir, 0700); err != nil {
		return nil, errors.New("media storage is unavailable")
	}
	file := filepath.Join(dir, id+"."+ext)
	if err = os.WriteFile(file, data, 0600); err != nil {
		return nil, errors.New("image could not be saved")
	}
	asset := &models.ProductMedia{ID: id, OwnerID: owner, URL: "/media/" + id, Kind: kind, Width: w, Height: h, OptionsJSON: options, CreatedAt: time.Now()}
	if err = facades.Orm().Query().Create(asset); err != nil {
		_ = os.Remove(file)
		return nil, errors.New("image could not be registered")
	}
	return asset, nil
}
func (c *MediaController) Policy(ctx http.Context) http.Response {
	u := appmiddleware.CurrentUser(ctx)
	if u == nil || !u.Can("media.settings.manage", "") {
		return mediaError(ctx, 403, "Media policy administration is not permitted.")
	}
	var req services.MediaPolicy
	if err := mediaDecode(ctx, &req); err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	if err := services.ValidateMediaPolicy(req); err != nil {
		return mediaError(ctx, 422, err.Error())
	}
	saved, err := services.SaveMediaPolicy(req)
	if errors.Is(err, services.ErrMediaPolicyConflict) {
		return mediaError(ctx, 409, err.Error())
	}
	if err != nil {
		return mediaError(ctx, 503, "Media policy could not be saved.")
	}
	return ctx.Response().Json(200, map[string]any{"policy": saved})
}

// Show keeps unreferenced drafts private. Product references are the publication boundary.
func (c *MediaController) Show(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")
	if _, err := uuid.Parse(id); err != nil {
		return mediaError(ctx, 404, "Image not found.")
	}
	var asset models.ProductMedia
	if err := facades.Orm().Query().Where("id = ?", id).First(&asset); err != nil || asset.ID == "" {
		return mediaError(ctx, 404, "Image not found.")
	}
	allowed := false
	if u := appmiddleware.CurrentUser(ctx); u != nil && (u.ID == asset.OwnerID || u.Can("media.create", asset.OwnerID)) {
		allowed = true
	}
	if !allowed {
		var products []models.Product
		if err := facades.Orm().Query().Where("image_url = ? OR gallery_urls LIKE ?", asset.URL, "%"+asset.URL+"%").Get(&products); err != nil {
			return mediaError(ctx, 404, "Image not found.")
		}
		for _, p := range products {
			if p.ImageUrl == asset.URL {
				allowed = true
				break
			}
			p.Hydrate()
			for _, u := range p.GalleryUrls {
				if u == asset.URL {
					allowed = true
					break
				}
			}
		}
	}
	if !allowed {
		return mediaError(ctx, 404, "Image not found.")
	}
	ext, mime := "png", "image/png"
	if asset.Kind == "mockup" {
		ext, mime = "svg", "image/svg+xml"
	}
	if strings.ContainsAny(asset.ID, "/\\") {
		return mediaError(ctx, 404, "Image not found.")
	}
	data, err := os.ReadFile(filepath.Join("storage", "app", "product-media", asset.ID+"."+ext))
	if err != nil {
		return mediaError(ctx, 404, "Image not found.")
	}
	return ctx.Response().Header("Cache-Control", "private, no-store").Header("X-Content-Type-Options", "nosniff").Header("Content-Security-Policy", "default-src 'none'; img-src data:; style-src 'unsafe-inline'; sandbox").Data(200, mime, data)
}
