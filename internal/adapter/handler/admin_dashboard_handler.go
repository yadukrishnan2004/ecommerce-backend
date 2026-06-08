package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/yadukrishnan2004/ecommerce-backend/internal/utils/response"
)

func (h *AdminHandler) GetDashboardGraphs(c *fiber.Ctx) error {
	graphs, err := h.svc.GetDashboardGraphs(c.Context())
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "Failed to fetch dashboard graphs", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "Dashboard graphs fetched successfully", graphs, nil)
}

func (h *AdminHandler) GetDashboardKPIs(c *fiber.Ctx) error {
	kpis, err := h.svc.GetDashboardKPIs(c.Context())
	if err != nil {
		return response.Response(c, http.StatusInternalServerError, "failed to fetch dashboard KPIs", nil, err.Error())
	}

	return response.Response(c, http.StatusOK, "dashboard KPIs fetched successfully", kpis, nil)
}
