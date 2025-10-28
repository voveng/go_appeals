package handlers

import (
	"fmt"

	"go_appeals/internal/models"
	"go_appeals/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	Service *services.AppealService
}

func (h *Handlers) GetStartedAppeals(c *fiber.Ctx) error {
	appeals, err := h.Service.GetStartedAppeals()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"appeals": appeals,
	})
}

func (h *Handlers) GetAllAppeals(c *fiber.Ctx) error {
	appeals, err := h.Service.GetAllAppeals()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"appeals": appeals,
	})
}

func (h *Handlers) GetAppealByID(c *fiber.Ctx) error {
	id := c.Params("id")
	appeal, err := h.Service.GetAppealByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"appeal": appeal,
	})
}

func (h *Handlers) CreateAppeal(c *fiber.Ctx) error {
	var req models.CreateAppealRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	appeal, err := h.Service.CreateAppeal(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Appeal created successfully",
		"appeal":  appeal,
	})
}

func (h *Handlers) StartProcessing(c *fiber.Ctx) error {
	id := c.Params("id")

	appeal, err := h.Service.StartProcessing(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("appeal with ID %s not found", id) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Appeal started processing successfully",
		"appeal":  appeal,
	})
}

func (h *Handlers) CompleteAppeal(c *fiber.Ctx) error {
	id := c.Params("id")

	var req models.UpdateAppealSolutionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	appeal, err := h.Service.CompleteAppeal(id, req)
	if err != nil {
		if err.Error() == fmt.Sprintf("appeal with ID %s not found", id) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Appeal completed successfully",
		"appeal":  appeal,
	})
}

func (h *Handlers) CancelAppeal(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.Service.CancelAppeal(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("appeal with ID %s not found", id) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Appeal canceled successfully",
		"id":      id,
	})
}

func (h *Handlers) CancelAllInProgress(c *fiber.Ctx) error {
	if err := h.Service.CancelAllInProgress(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "All in progress appeals canceled successfully",
	})
}

func (h *Handlers) GetAppealsByDates(c *fiber.Ctx) error {
	layout := "2006-01-02"

	startStr := c.Query("startDate")
	endStr := c.Query("endDate")

	if startStr == "" || endStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "startDate and endDate are required",
		})
	}

	start, err := time.Parse(layout, startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid startDate format, use YYYY-MM-DD",
		})
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid endDate format, use YYYY-MM-DD",
		})
	}

	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	appeals, err := h.Service.GetAppealsByDates(start, end)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"appeals": appeals,
	})
}
