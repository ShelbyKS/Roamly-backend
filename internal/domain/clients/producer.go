package clients

import (
	"github.com/ShelbyKS/Roamly-backend/internal/domain/model"
)

type IMessageProdcuer interface {
	SendMessage(msg model.Message) error
}
