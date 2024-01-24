package crawler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
	"github.com/brunofjesus/pricetracker/stores/motocard/pkg/definition"
	"github.com/brunofjesus/pricetracker/stores/motocard/pkg/definition/imagegen"
	"log/slog"
	"net/url"
	"strings"
)

func mapMotocardProductToStoreProduct(
	logger *slog.Logger,
	storeSlug string,
	in definition.MotorcardProductResult,
	currency string,
) dto.StoreProduct {

	imageUrl, err := generateImageLink(in.Image)
	if err != nil {
		logger.Warn(
			"cannot generate image link",
			slog.String("product", in.URL),
			slog.Any("error", err),
		)
	}

	return dto.StoreProduct{
		StoreSlug: storeSlug,
		EAN:       collectEANs(in),
		SKU:       []string{in.References.Motocard},
		Name:      in.Title.Name,
		Brand:     in.Title.Brand,
		Price:     int64(in.RawPrice),
		Available: in.NumSizesWithStock > 0,
		ImageLink: imageUrl,
		Link:      in.URL,
		Currency:  strings.ToUpper(currency),
	}
}

func collectEANs(in definition.MotorcardProductResult) []string {
	if len(in.Sizes) == 0 {
		return nil
	}

	var result []string
	for _, size := range in.Sizes {
		if len(size.Gtin) > 0 {
			result = append(result, size.Gtin)
		}
	}

	return result
}

func generateImageLink(source string) (string, error) {
	parsedUrl, err := url.Parse(source)
	if err != nil {
		return "", fmt.Errorf("could not parse image link: %w", err)
	}

	v := parsedUrl.Query().Get("v")
	path := parsedUrl.Path

	genQuery := imagegen.Query{
		Bucket: "motocard",
		Key:    path,
		Edits: imagegen.QueryEdits{
			Webp: imagegen.QueryEditQuality{Quality: 85},
			Jpeg: imagegen.QueryEditQuality{Quality: 91},
			Resize: imagegen.QueryEditResize{
				Width:  550,
				Height: 550,
				Fit:    "cover",
			},
		},
		V: v,
	}

	b, err := json.Marshal(genQuery)
	if err != nil {
		return "", fmt.Errorf("could not marshall image generation query to json: %w", err)
	}

	enc := base64.URLEncoding.EncodeToString(b)
	return fmt.Sprintf("https://images.motocard.com/%s", enc), nil
}
