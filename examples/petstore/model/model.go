package model

type Pet struct {
	ID       int64    `json:"id" example:"1"`
	Name     string   `json:"name" binding:"required" example:"doggie"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls" binding:"required"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold" example:"available"`
}

type Category struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Dogs"`
}

type Tag struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"friendly"`
}

type CreatePetRequest struct {
	Name     string   `json:"name" binding:"required"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls" binding:"required"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold"`
}

type UpdatePetRequest struct {
	ID       int64    `json:"id" binding:"required"`
	Name     string   `json:"name"`
	Category Category `json:"category"`
	PhotoURL []string `json:"photoUrls"`
	Tags     []Tag    `json:"tags"`
	Status   string   `json:"status" enums:"available,pending,sold"`
}

type UpdatePetFormData struct {
	ID     int64  `path:"petId" binding:"required"`
	Name   string `form:"name" binding:"required"`
	Status string `form:"status" binding:"required" enums:"available,pending,sold"`
}

type FindPetsByStatusRequest struct {
	Status string `form:"status" binding:"required" enums:"available,pending,sold"`
}

type FindPetsByTagsRequest struct {
	Tags []string `form:"tags" binding:"required"`
}

type User struct {
	ID         int64  `json:"id" example:"1"`
	Username   string `json:"username" example:"john_doe"`
	FirstName  string `json:"firstName" example:"John"`
	LastName   string `json:"lastName" example:"Doe"`
	Email      string `json:"email" example:"john_doe@example.com"`
	Password   string `json:"password" example:"password123"`
	Phone      string `json:"phone" example:"123-456-7890"`
	UserStatus int64  `json:"userStatus" example:"1"`
}

type ApiResponse struct {
	Code    int64  `json:"code" example:"200"`
	Type    string `json:"type" example:"success"`
	Message string `json:"message" example:"Pet created successfully"`
}
