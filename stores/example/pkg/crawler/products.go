package crawler

import (
	"github.com/brunofjesus/pricetracker/stores/connector/dto"
)

var products = []dto.StoreProduct{
	{
		StoreSlug: StoreSlug,
		EAN:       []string{"1000000000000"},
		SKU:       []string{"PN101"},
		Name:      "Microscope",
		Brand:     "ExtremeZoom",
		Price:     500000,
		Available: false,
		ImageLink: "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3a/Optical_microscope_nikon_alphaphot.jpg/317px-Optical_microscope_nikon_alphaphot.jpg",
		Link:      "http://localhost/pineapple",
		Currency:  "EUR",
	},
	{
		StoreSlug: StoreSlug,
		EAN:       []string{"2000000000000"},
		SKU:       []string{"PN201"},
		Name:      "Model Glider",
		Brand:     "Sindbad",
		Price:     100000,
		Available: false,
		ImageLink: "https://upload.wikimedia.org/wikipedia/commons/7/72/KN_Modelloldtimer_Sindbad.JPG",
		Link:      "http://localhost/glider",
		Currency:  "EUR",
	},
	{
		StoreSlug: StoreSlug,
		EAN:       []string{"3000000000000"},
		SKU:       []string{"PN301"},
		Name:      "Binoculars",
		Brand:     "ExtremeZoom",
		Price:     83000,
		Available: false,
		ImageLink: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Binocular_with_8x_magnification_and_42_mm_lens_diameter.jpg/320px-Binocular_with_8x_magnification_and_42_mm_lens_diameter.jpg",
		Link:      "http://localhost/binoculars",
		Currency:  "EUR",
	},
}
