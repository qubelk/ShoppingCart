package product

import (
	"pkg/testutil"
	"strings"
	"testing"
)

func longString(length int) string {
	var b strings.Builder
	for range length {
		b.WriteString("a")
	}
	return b.String()
}

func TestValidateTitle(t *testing.T) {
	tests := []testutil.TestCase{
		{Name: "valid merged", Data: "test_product", WantErr: false},
		{Name: "valid splitted", Data: "test product", WantErr: false},
		{Name: "invalid short", Data: "p", WantErr: true},
		{Name: "invalid empty", Data: "", WantErr: true},
		{Name: "invalid space", Data: " ", WantErr: true},
		{Name: "invalid tab", Data: "	", WantErr: true},
		{Name: "invalid too long", Data: longString(256), WantErr: true},
	}

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return validateTitle(tt.Data.(string))
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []testutil.TestCase{
		{Name: "valid", Data: "test", WantErr: false},
		{Name: "invalid empty", Data: "", WantErr: true},
		{Name: "invalid too long", Data: longString(1001), WantErr: true},
	}

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return validateDescription(tt.Data.(string))
		})
	}
}

func TestValidatePrice(t *testing.T) {
	tests := []testutil.TestCase{
		{Name: "valid", Data: 9.99, WantErr: false},
		{Name: "invalid zero", Data: 0.0, WantErr: true},
		{Name: "invalid negative", Data: -1.1, WantErr: true},
	}

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return validatePrice(tt.Data.(float64))
		})
	}
}

func TestValidateCount(t *testing.T) {
	tests := []testutil.TestCase{
		{Name: "valid", Data: 1, WantErr: false},
		{Name: "invalid zero", Data: 0, WantErr: true},
		{Name: "invalid negative", Data: -1, WantErr: true},
	}

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return validateCount(tt.Data.(int))
		})
	}
}

func TestCreateProductRequestValidate(t *testing.T) {
	tests := []testutil.TestCase{
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
				Title:       longString(256),
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
				Description: longString(1001),
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

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return tt.Data.(*CreateProductRequest).Validate()
		})
	}
}

func TestSearchProductRequestValidate(t *testing.T) {
	tests := []testutil.TestCase{
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
				Title: longString(256),
			},
			WantErr: true,
		},
	}

	helper := testutil.New(t)
	for _, tt := range tests {
		helper.RunTest(tt, func() error {
			return tt.Data.(*SearchProductRequest).Validate()
		})
	}
}
