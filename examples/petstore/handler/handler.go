package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/fiberopenapi/examples/petstore/model"
	"github.com/oaswrap/fiberopenapi/examples/petstore/repository"
)

type PetHandler struct {
	repo repository.PetRepository
}

func NewPetHandler(repo repository.PetRepository) *PetHandler {
	return &PetHandler{
		repo: repo,
	}
}

func (h *PetHandler) GetAllPets(c *fiber.Ctx) error {
	pets := h.repo.FindAll()
	return c.JSON(pets)
}

func (h *PetHandler) CreatePet(c *fiber.Ctx) error {
	var req model.CreatePetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	pet, err := h.repo.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create pet",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(pet)
}

func (h *PetHandler) UpdatePet(c *fiber.Ctx) error {
	var req model.UpdatePetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	pet, err := h.repo.Update(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update pet",
		})
	}

	return c.JSON(pet)
}

func (h *PetHandler) GetPetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pet ID",
		})
	}

	pet, found := h.repo.FindByID(int64(id))
	if !found {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Pet not found",
		})
	}

	return c.JSON(pet)
}

func (h *PetHandler) FindPetsByStatus(c *fiber.Ctx) error {
	var req model.FindPetsByStatusRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	pets := h.repo.FindAllByStatus(req.Status)

	return c.JSON(pets)
}

func (h *PetHandler) FindPetsByTags(c *fiber.Ctx) error {
	var req model.FindPetsByTagsRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	pets := h.repo.FindAllByTags(req.Tags)

	return c.JSON(pets)
}

func (h *PetHandler) DeletePet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid pet ID",
		})
	}

	if err := h.repo.DeleteByID(int64(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete pet",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *PetHandler) UpdatePetFormData(c *fiber.Ctx) error {
	var req model.UpdatePetFormData
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	pet, err := h.repo.UpdateFormData(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update pet",
		})
	}

	return c.JSON(pet)
}
