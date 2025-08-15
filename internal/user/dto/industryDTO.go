package dto

type ListIndustriesResponse struct {
	ID           int    `json:"id"`            // 行业ID
	IndustryCode string `json:"industry_code"` // 行业编码
	IndustryName string `json:"industry_name"` // 行业名称
}
