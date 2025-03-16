package handler

import (
	"github.com/osmanunal/commerce-oms/product-service/internal/model"
	"github.com/osmanunal/commerce-oms/product-service/internal/service"
	"github.com/osmanunal/commerce-oms/product-service/server/viewmodel"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	products, err := h.productService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var productResponses []viewmodel.ProductResponse
	for _, product := range products {
		productResponse := viewmodel.ProductResponse{}
		productResponse.FromModel(product)
		productResponses = append(productResponses, productResponse)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"products": productResponses})
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	productRequest := viewmodel.ProductRequest{}
	if err := c.BodyParser(&productRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	product := productRequest.ToModel(model.Product{})

	err := h.productService.Create(c.Context(), &product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": fiber.Map{
			"id": product.ID.String(),
		},
	})
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	productID := c.Params("id")
	product, err := h.productService.GetByID(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ürünün alınması sırasında bir hata oluştu"})
	}

	var productResponse viewmodel.ProductResponse
	productResponse.FromModel(*product)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"product": productResponse})
}
