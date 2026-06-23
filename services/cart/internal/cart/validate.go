package cart

import (
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidID
	}

	return uuid.Validate(id.String())
}

func validateTitle(title string) error {
	if len(title) < 3 || len(title) > 255 {
		return ErrInvalidTitle
	}

	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return ErrInvalidPrice
	}

	return nil
}

func validateQuantity(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	return nil
}

func (r *AddItemRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateID(r.UserID)
	})

	g.Go(func() error {
		return validateID(r.Item.ProductID)
	})

	g.Go(func() error {
		return validateTitle(r.Item.Title)
	})

	g.Go(func() error {
		return validatePrice(r.Item.Price)
	})

	g.Go(func() error {
		return validateQuantity(r.Item.Quantity)
	})

	return g.Wait()
}

func (r *UpdateQuantityRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateID(r.UserID)
	})

	g.Go(func() error {
		return validateID(r.ProductID)
	})

	g.Go(func() error {
		return validateQuantity(r.Quantity)
	})

	return g.Wait()
}

func (r *RemoveItemRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateID(r.UserID)
	})

	g.Go(func() error {
		return validateID(r.ProductID)
	})

	return g.Wait()
}

func (r *ClearCartRequest) Validate() error {
	return validateID(r.UserID)
}

func (r *GetCartRequest) Validate() error {
	return validateID(r.UserID)
}

func (r *GetCartTTLRequest) Validate() error {
	return validateID(r.UserID)
}
