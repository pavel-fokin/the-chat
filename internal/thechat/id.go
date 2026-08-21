// Package thechat holds the domain entities for The Chat.
package thechat

import "github.com/google/uuid"

// ID is a phantom-typed UUID: T pins it to one entity type so IDs from
// different entities can't be mixed up.
type ID[T any] struct{ uuid.UUID }
