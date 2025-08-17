package models

// AdminDetailResponse defines the detailed admin response structure
type AdminDetailResponse struct {
	AdminResponse
	Restaurants []RestaurantResponse `json:"restaurants"`
}

// RestaurantDetailResponse defines the detailed restaurant response structure
type RestaurantDetailResponse struct {
	RestaurantResponse
	MenuItems []MenuResponse `json:"menu_items"`
}

// PageInfo defines the pagination information
type PageInfo struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
	PerPage     int `json:"per_page"`
}

// PaginatedResponse defines the paginated response structure
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	PageInfo PageInfo    `json:"page_info"`
}

// UpdateAdminStatusRequest defines the request to update admin status
type UpdateAdminStatusRequest struct {
	Status string `json:"status" example:"active" enums:"active,inactive"`
}

// UpdateRestaurantStatusRequest defines the request to update restaurant status
type UpdateRestaurantStatusRequest struct {
	Status string `json:"status" example:"active" enums:"active,inactive,pending"`
}

// UpdateProfileRequest defines the request to update profile
type UpdateProfileRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}
