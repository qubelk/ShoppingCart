package cart

import (
	"pkg/testutils"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func longString(length uint) string {
	var b strings.Builder

	for range length {
		b.WriteString("a")
	}

	return b.String()
}

func TestValidateID(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: uuid.New(), WantErr: false},
		{Name: "invalid: nil", Data: uuid.Nil, WantErr: true},
		{Name: "invalid: empty", Data: uuid.UUID{}, WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validateID(a.(uuid.UUID))
	})
}

func TestValidateTitle(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: "valid title", WantErr: false},
		{Name: "invalid: short", Data: "in", WantErr: true},
		{Name: "invalid: long", Data: longString(256), WantErr: true},
		{Name: "invalid: empty", Data: "", WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validateTitle(a.(string))
	})
}

func TestValidatePrice(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: 9.99, WantErr: false},
		{Name: "invalid: zero", Data: 0.0, WantErr: true},
		{Name: "invalid: negative", Data: -1.1, WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validatePrice(a.(float64))
	})
}

func TestValidateQuantity(t *testing.T) {
	tests := []testutils.TestCase{
		{Name: "valid", Data: 1, WantErr: false},
		{Name: "invalid: zero", Data: 0, WantErr: true},
		{Name: "invalid: negative", Data: -1, WantErr: true},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return validateQuantity(a.(int))
	})
}

func TestAddItemValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil user id",
			Data: &AddItemRequest{
				UserID: uuid.Nil,
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: nil product id",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.Nil,
					Title:     "item title",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty user id",
			Data: &AddItemRequest{
				UserID: uuid.UUID{},
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty product id",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.UUID{},
					Title:     "item title",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: short title",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "tt",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: long title",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     longString(256),
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty title",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "",
					Price:     9.99,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: negative price",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     -9.99,
					Quantity:  1,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: zero price",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     0,
					Quantity:  2,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: negative quantity",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     9.99,
					Quantity:  -1,
				},
			},
			WantErr: true,
		},
		{
			Name: "invalid: zero quantity",
			Data: &AddItemRequest{
				UserID: uuid.New(),
				Item: CartItem{
					ProductID: uuid.New(),
					Title:     "item title",
					Price:     9.99,
					Quantity:  0,
				},
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*AddItemRequest).Validate()
	})
}

func TestUpdateQuantityValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.New(),
				ProductID: uuid.New(),
				Quantity:  1,
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil user id",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.Nil,
				ProductID: uuid.New(),
				Quantity:  1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty user id",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.UUID{},
				ProductID: uuid.New(),
				Quantity:  1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: nil product id",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.New(),
				ProductID: uuid.Nil,
				Quantity:  1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty product id",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.New(),
				ProductID: uuid.UUID{},
				Quantity:  1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: negative quantity",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.New(),
				ProductID: uuid.New(),
				Quantity:  -1,
			},
			WantErr: true,
		},
		{
			Name: "invalid: zero quantity",
			Data: &UpdateQuantityRequest{
				UserID:    uuid.New(),
				ProductID: uuid.New(),
				Quantity:  0,
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*UpdateQuantityRequest).Validate()
	})
}

func TestRemoveItemValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &RemoveItemRequest{
				UserID:    uuid.New(),
				ProductID: uuid.New(),
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil user id",
			Data: &RemoveItemRequest{
				UserID:    uuid.Nil,
				ProductID: uuid.New(),
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty user id",
			Data: &RemoveItemRequest{
				UserID:    uuid.UUID{},
				ProductID: uuid.New(),
			},
			WantErr: true,
		},
		{
			Name: "invalid: nil product id",
			Data: &RemoveItemRequest{
				UserID:    uuid.New(),
				ProductID: uuid.Nil,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty product id",
			Data: &RemoveItemRequest{
				UserID:    uuid.New(),
				ProductID: uuid.UUID{},
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*RemoveItemRequest).Validate()
	})
}

func TestClearCartValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &ClearCartRequest{
				UserID: uuid.New(),
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil uuid",
			Data: &ClearCartRequest{
				UserID: uuid.Nil,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty uuid",
			Data: &ClearCartRequest{
				UserID: uuid.UUID{},
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*ClearCartRequest).Validate()
	})
}

func TestGetCartValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &GetCartRequest{
				UserID: uuid.New(),
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil uuid",
			Data: &GetCartRequest{
				UserID: uuid.Nil,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty uuid",
			Data: &GetCartRequest{
				UserID: uuid.UUID{},
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*GetCartRequest).Validate()
	})
}

func TestGetCartTTLValidate(t *testing.T) {
	tests := []testutils.TestCase{
		{
			Name: "valid",
			Data: &GetCartTTLRequest{
				UserID: uuid.New(),
			},
			WantErr: false,
		},
		{
			Name: "invalid: nil uuid",
			Data: &GetCartTTLRequest{
				UserID: uuid.Nil,
			},
			WantErr: true,
		},
		{
			Name: "invalid: empty uuid",
			Data: &GetCartTTLRequest{
				UserID: uuid.UUID{},
			},
			WantErr: true,
		},
	}

	testutils.UnitTestRunner(t, tests, func(a any) error {
		return a.(*GetCartTTLRequest).Validate()
	})
}
