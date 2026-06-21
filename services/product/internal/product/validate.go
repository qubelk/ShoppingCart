package product

import (
	"golang.org/x/sync/errgroup"
)

func validateTitle(title string) error {
	if len(title) < 3 {
		return ErrInvalidTitle
	}

	return nil
}

func validateDescription(description string) error {
	if description == "" {
		return ErrInvalidDescription
	}

	return nil
}

func validatePrice(price float64) error {
	if price < 0 {
		return ErrInvalidPrice
	}

	return nil
}

func validateCount(count int) error {
	if count <= 0 {
		return ErrInvalidCount
	}

	return nil
}

func (r *CreateProductRequest) Validate() error {
	var g errgroup.Group

	g.Go(func() error {
		return validateTitle(r.Title)
	})

	g.Go(func() error {
		return validateDescription(r.Description)
	})

	g.Go(func() error {
		return validatePrice(r.Price)
	})

	g.Go(func() error {
		return validateCount(r.Count)
	})

	return g.Wait()
}
