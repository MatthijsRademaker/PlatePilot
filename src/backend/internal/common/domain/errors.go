package domain

import "errors"

// Domain errors
var (
	ErrRecipeNotFound           = errors.New("recipe not found")
	ErrIngredientNotFound       = errors.New("ingredient not found")
	ErrCuisineNotFound          = errors.New("cuisine not found")
	ErrAllergyNotFound          = errors.New("allergy not found")
	ErrInvalidInput             = errors.New("invalid input")
	ErrDuplicateEntry           = errors.New("duplicate entry")
	ErrShoppingListNotFound     = errors.New("shopping list not found")
	ErrShoppingListItemNotFound = errors.New("shopping list item not found")
	ErrEmptyMealPlan            = errors.New("meal plan has no recipes")
)

// DomainError wraps domain-specific errors with additional context
type DomainError struct {
	Err     error
	Message string
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Err.Error()
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error with context
func NewDomainError(err error, message string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
	}
}
