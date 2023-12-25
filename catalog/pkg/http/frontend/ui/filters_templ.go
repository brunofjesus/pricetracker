// Code generated by templ - DO NOT EDIT.

// templ: version: 0.2.476
package ui

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import product "github.com/brunofjesus/pricetracker/catalog/pkg/product"
import db_store "github.com/brunofjesus/pricetracker/catalog/internal/repository/store"
import "github.com/brunofjesus/pricetracker/catalog/pkg/pagination"
import "github.com/brunofjesus/pricetracker/catalog/util/nulltype"
import "github.com/brunofjesus/pricetracker/catalog/pkg/http/frontend/util"
import "fmt"

func FiltersComponent(page pagination.PaginatedQuery, filters product.FinderFilters, stores []db_store.Store) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"d-flex gap-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if filters.StoreId > 0 {
			templ_7745c5c3_Err = FilterBadgeComponent("Store", findStoreName(filters.StoreId, stores), generateUrl(page, filters, "StoreId")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Price", float64ToString(filters.MinPrice/100), generateUrl(page, filters, "MinPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Price", float64ToString(filters.MaxPrice/100), generateUrl(page, filters, "MaxPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(filters.NameLike) > 0 {
			templ_7745c5c3_Err = FilterBadgeComponent("Name", filters.NameLike, generateUrl(page, filters, "NameLike")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(filters.BrandLike) > 0 {
			templ_7745c5c3_Err = FilterBadgeComponent("Brand", filters.BrandLike, generateUrl(page, filters, "BrandLike")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if nulltype.IsUndefined(filters.Available) == false {
			templ_7745c5c3_Err = FilterBadgeComponent("Available", nulltype.ToString(filters.Available), generateUrl(page, filters, "Available")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if len(filters.ProductUrl) > 0 {
			templ_7745c5c3_Err = FilterBadgeComponent("Product URL", "Yes", generateUrl(page, filters, "ProductUrl")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinDifference > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Difference", float64ToString(filters.MinDifference/100), generateUrl(page, filters, "MinDifference")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxDifference > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Difference", float64ToString(filters.MaxDifference/100), generateUrl(page, filters, "MaxDifference")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinDiscountPercent > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Discount (%)", float64ToString(filters.MinDiscountPercent), generateUrl(page, filters, "MinDiscountPercent")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxDiscountPercent > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Discount (%)", float64ToString(filters.MaxDiscountPercent), generateUrl(page, filters, "MaxDiscountPercent")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinAveragePrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Average Price", float64ToString(filters.MinAveragePrice/100), generateUrl(page, filters, "MinAveragePrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxAveragePrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Average Price", float64ToString(filters.MaxAveragePrice/100), generateUrl(page, filters, "MaxAveragePrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinMinimumPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Lower Price", float64ToString(filters.MinMinimumPrice/100), generateUrl(page, filters, "MinMinimumPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxMinimumPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Lower Price", float64ToString(filters.MaxMinimumPrice/100), generateUrl(page, filters, "MaxMinimumPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MinMaximumPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Min. Highest Price", float64ToString(filters.MinMaximumPrice/100), generateUrl(page, filters, "MinMaximumPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		if filters.MaxMaximumPrice > -1 {
			templ_7745c5c3_Err = FilterBadgeComponent("Max. Highest Price", float64ToString(filters.MaxMaximumPrice/100), generateUrl(page, filters, "MaxMaximumPrice")).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func generateUrl(page pagination.PaginatedQuery, filters product.FinderFilters, toRemove string) string {
	switch toRemove {
	case "StoreId":
		filters.StoreId = -1
	case "MinPrice":
		filters.MinPrice = -1
	case "MaxPrice":
		filters.MaxPrice = -1
	case "NameLike":
		filters.NameLike = ""
	case "BrandLike":
		filters.BrandLike = ""
	case "Available":
		filters.Available = -1
	case "ProductUrl":
		filters.ProductUrl = ""
	case "MinDifference":
		filters.MinDifference = -1
	case "MaxDifference":
		filters.MaxDifference = -1
	case "MinDiscountPercent":
		filters.MinDiscountPercent = -1
	case "MaxDiscountPercent":
		filters.MaxDiscountPercent = -1
	case "MinAveragePrice":
		filters.MinAveragePrice = -1
	case "MaxAveragePrice":
		filters.MaxAveragePrice = -1
	case "MinMinimumPrice":
		filters.MinMinimumPrice = -1
	case "MaxMinimumPrice":
		filters.MaxMinimumPrice = -1
	case "MinMaximumPrice":
		filters.MinMaximumPrice = -1
	case "MaxMaximumPrice":
		filters.MaxMaximumPrice = -1
	}
	return util.QueryParamString(page, filters)
}

func findStoreName(id int, stores []db_store.Store) string {
	for _, store := range stores {
		if store.StoreId == int64(id) {
			return store.Name
		}
	}

	return ""
}

func float64ToString(n float64) string {
	return fmt.Sprintf("%.2f", n)
}
