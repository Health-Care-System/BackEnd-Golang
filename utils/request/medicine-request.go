package request

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToMedicineRequest(medicine web.MedicineRequest) *schema.Medicine {
	return &schema.Medicine{
		Code:     medicine.Code,
		Name:     medicine.Name,
		Merk:     medicine.Merk,
		Category: medicine.Category,
		Type:     medicine.Type,
		Stock:    medicine.Stock,
		Price:    medicine.Price,
		Details:  medicine.Details,
		Image:    medicine.Image,
	}
}
