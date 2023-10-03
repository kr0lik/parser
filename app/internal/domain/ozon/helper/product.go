package helper

import (
	"github.com/google/uuid"
	"parser/internal/domain/ozon/dto"
)

func GenerateProductUuid(product dto.Product) uuid.UUID {
	return uuid.NewMD5(uuid.NameSpaceURL, []byte(product.Category.Id+"-"+product.Name))
}
