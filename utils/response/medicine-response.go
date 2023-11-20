package response

import (
	"healthcare/models/schema"
	"healthcare/models/web"
)

func ConvertToAdminMedicineResponse(medicine *schema.Medicine) web.MedicineResponse {
	return web.MedicineResponse{
		ID:       medicine.ID,
		Code:     medicine.Code,
		Name:     medicine.Name,
		Merk:     medicine.Merk,
		Category: medicine.Category,
		Type:     medicine.Type,
		Stock:    medicine.Stock,
		Price:    medicine.Price,
		Details:  medicine.Details,
		Image:    medicine.Image,
		CreatedAt: medicine.CreatedAt,
	}
}

func ConvertToAdminGetAllMedicinesResponse(medicines []schema.Medicine) []web.MedicineResponse {
	var results []web.MedicineResponse
	for _, medicine := range medicines {
		medicineResponse := web.MedicineResponse{
			ID:       medicine.ID,
			Code:     medicine.Code,
			Name:     medicine.Name,
			Merk:     medicine.Merk,
			Category: medicine.Category,
			Type:     medicine.Type,
			Stock:    medicine.Stock,
			Price:    medicine.Price,
			Details:  medicine.Details,
			Image:    medicine.Image,
			CreatedAt: medicine.CreatedAt,
		}
		results = append(results, medicineResponse)
	}
	return results
}

func ConvertToUserMedicineResponse(medicine *schema.Medicine) web.MedicineResponse {
	return web.MedicineResponse{
		Name:     medicine.Name,
		Code:     medicine.Code,
		Merk:     medicine.Merk,
		Category: medicine.Category,
		Type:     medicine.Type,
		Stock:    medicine.Stock,
		Price:    medicine.Price,
		Details:  medicine.Details,
		Image:    medicine.Image,
	}
}

func ConvertToUserGetAllMedicinesResponse(medicines []schema.Medicine) []web.MedicineResponse {
	var results []web.MedicineResponse
	for _, medicine := range medicines {
		medicineResponse := web.MedicineResponse{
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
		results = append(results, medicineResponse)
	}
	return results
}
