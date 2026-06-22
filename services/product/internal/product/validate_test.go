package product

import (
	"pkg/tests"
	"strings"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	var b strings.Builder
	for range 256 {
		b.WriteString("a")
	}

	test := []tests.TestCase{
		{Name: "valid merged", Data: "test_product", WantErr: false},
		{Name: "valid splitted", Data: "test product", WantErr: false},
		{Name: "invalid short", Data: "p", WantErr: true},
		{Name: "invalid empty", Data: "", WantErr: true},
		{Name: "invalid space", Data: " ", WantErr: true},
		{Name: "invalid tab", Data: "	", WantErr: true},
		{Name: "invalid to long", Data: b.String(), WantErr: true},
	}

	tests.UnitTest(t, test, func(a any) error {
		return validateTitle(a.(string))
	})
}

func TestValidateDescription(t *testing.T) {
	var b strings.Builder
	for range 1001 {
		b.WriteString("a")
	}

	test := []tests.TestCase{
		{Name: "valid", Data: "test", WantErr: false},
		{Name: "invalid empty", Data: "", WantErr: true},
		{Name: "invalid to long", Data: b.String(), WantErr: true},
	}

	tests.UnitTest(t, test, func(a any) error {
		return validateDescription(a.(string))
	})
}

func TestValidatePrice(t *testing.T) {
	test := []tests.TestCase{
		{Name: "valid", Data: 9.99, WantErr: false},
		{Name: "invalid zero", Data: 0.0, WantErr: true},
		{Name: "invalid negative", Data: -1.1, WantErr: true},
	}

	tests.UnitTest(t, test, func(a any) error {
		return validatePrice(a.(float64))
	})
}

func TestValidateCount(t *testing.T) {
	test := []tests.TestCase{
		{Name: "valid", Data: 1, WantErr: false},
		{Name: "invalid zero", Data: 0, WantErr: true},
		{Name: "invalid negative", Data: -1, WantErr: true},
	}

	tests.UnitTest(t, test, func(a any) error {
		return validateCount(a.(int))
	})
}

func TestValidateCreateProductRequest(t *testing.T) {
	var b strings.Builder
	for range 1001 {
		b.WriteString("a")
	}

	test := []tests.TestCase{
		{
			Name: "valid",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "product description",
				Price:       9.99,
				Count:       1,
			},
			WantErr: false,
		},
		{
			Name: "invalid: empty title",
			Data: &CreateProductRequest{
				Title:       "",
				Description: "product description",
				Price:       9.99,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: short title",
			Data: &CreateProductRequest{
				Title:       "sh",
				Description: "product description",
				Price:       9.99,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: too long title",
			Data: &CreateProductRequest{
				Title:       b.String(),
				Description: "product description",
				Price:       9.99,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty description",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "",
				Price:       9.99,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: too long description",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: b.String(),
				Price:       9.99,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: zero price",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "product description",
				Price:       0,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: negative price",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "product description",
				Price:       -1,
				Count:       1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: zero count",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "product description",
				Price:       9.99,
				Count:       0,
			},
			WantErr: true,
		},
		{
			Name: "invalid: negative count",
			Data: &CreateProductRequest{
				Title:       "product title",
				Description: "product description",
				Price:       9.99,
				Count:       -1,
			},
			WantErr: true,
		},
	}

	tests.UnitTest(t, test, func(a any) error {
		return a.(*CreateProductRequest).Validate()
	})
}

func TestValidateSearchProductRequest(t *testing.T) {
	var b strings.Builder
	for range 256 {
		b.WriteString("a")
	}

	test := []tests.TestCase{
		{
			Name: "valid",
			Data: &SearchProductRequest{
				Title: "product title",
			},
			WantErr: false,
		},
		{
			Name: "invalid: empty title",
			Data: &SearchProductRequest{
				Title: "",
			},
			WantErr: true,
		},
		{
			Name: "invalid: too short title",
			Data: &SearchProductRequest{
				Title: "pr",
			},
			WantErr: true,
		},
		{
			Name: "invalid: too long title",
			Data: &SearchProductRequest{
				Title: b.String(),
			},
			WantErr: true,
		},
	}

	tests.UnitTest(t, test, func(a any) error {
		return a.(*SearchProductRequest).Validate()
	})
}
