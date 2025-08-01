package repository

import (
	"errors"
	"slices"
	"sort"
	"sync"

	"github.com/oaswrap/fiberopenapi/examples/petstore/model"
)

type PetRepository interface {
	Create(req model.CreatePetRequest) (model.Pet, error)
	Update(req model.UpdatePetRequest) (model.Pet, error)
	UpdateFormData(req model.UpdatePetFormData) (model.Pet, error)
	FindByID(id int64) (model.Pet, bool)
	FindAll() []model.Pet
	FindAllByStatus(status string) []model.Pet
	FindAllByTags(tags []string) []model.Pet
	DeleteByID(id int64) error
}

// In-memory implementation
type petRepositoryImpl struct {
	data   map[int64]model.Pet
	nextID int64
	mutex  sync.RWMutex
}

func NewDummyPetRepository() PetRepository {
	return &petRepositoryImpl{
		data:   make(map[int64]model.Pet),
		nextID: 1,
	}
}

// Create new pet
func (r *petRepositoryImpl) Create(req model.CreatePetRequest) (model.Pet, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pet := model.Pet{
		ID:       r.nextID,
		Name:     req.Name,
		Category: req.Category,
		PhotoURL: req.PhotoURL,
		Tags:     req.Tags,
		Status:   req.Status,
	}

	r.data[r.nextID] = pet
	r.nextID++

	return pet, nil
}

// Update existing pet
func (r *petRepositoryImpl) Update(req model.UpdatePetRequest) (model.Pet, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pet, ok := r.data[req.ID]
	if !ok {
		return model.Pet{}, errors.New("pet not found")
	}

	// Update fields if provided (basic patch-like logic)
	if req.Name != "" {
		pet.Name = req.Name
	}
	pet.Category = req.Category
	if len(req.PhotoURL) > 0 {
		pet.PhotoURL = req.PhotoURL
	}
	pet.Tags = req.Tags
	if req.Status != "" {
		pet.Status = req.Status
	}

	r.data[req.ID] = pet

	return pet, nil
}

// UpdateFormData updates a pet using form data
func (r *petRepositoryImpl) UpdateFormData(req model.UpdatePetFormData) (model.Pet, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	pet, ok := r.data[req.ID]
	if !ok {
		return model.Pet{}, errors.New("pet not found")
	}

	pet.Name = req.Name
	pet.Status = req.Status
	r.data[req.ID] = pet

	return pet, nil
}

// FindByID
func (r *petRepositoryImpl) FindByID(id int64) (model.Pet, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pet, ok := r.data[id]
	return pet, ok
}

// FindAll (sorted by ID)
func (r *petRepositoryImpl) FindAll() []model.Pet {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pets := make([]model.Pet, 0, len(r.data))
	for _, pet := range r.data {
		pets = append(pets, pet)
	}

	sort.Slice(pets, func(i, j int) bool {
		return pets[i].ID < pets[j].ID
	})

	return pets
}

func (r *petRepositoryImpl) FindAllByStatus(status string) []model.Pet {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pets := make([]model.Pet, 0)
	for _, pet := range r.data {
		if pet.Status == status {
			pets = append(pets, pet)
		}
	}

	sort.Slice(pets, func(i, j int) bool {
		return pets[i].ID < pets[j].ID
	})

	return pets
}

func (r *petRepositoryImpl) FindAllByTags(tags []string) []model.Pet {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	pets := make([]model.Pet, 0)
	for _, pet := range r.data {
		for _, tag := range pet.Tags {
			if slices.Contains(tags, tag.Name) {
				pets = append(pets, pet)
				break // No need to check other tags for this pet
			}
		}
	}

	sort.Slice(pets, func(i, j int) bool {
		return pets[i].ID < pets[j].ID
	})

	return pets
}

// DeleteByID
func (r *petRepositoryImpl) DeleteByID(id int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.data[id]; !ok {
		return errors.New("pet not found")
	}

	delete(r.data, id)
	return nil
}
