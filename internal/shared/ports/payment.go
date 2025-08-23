package ports

import (
	"context"
	"modular_monolith/internal/shared/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentRepository interface {
	FindByOrderID(ctx context.Context, orderID primitive.ObjectID) (*model.Payment, error)
}
