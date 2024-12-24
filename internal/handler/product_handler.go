package handler

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rubenkristian/backend/internal/models"
	"github.com/rubenkristian/backend/internal/services"
	"github.com/rubenkristian/backend/pkg"
)

type ProductHandler struct {
	productService *services.ProductService
}

func InitializeProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (productHandler *ProductHandler) GetProduct(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad Request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	product, err := productHandler.productService.GetProduct(uint(productId))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Not found",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success get product",
		"data":    product,
	})
}

func (productHandler *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	take := c.QueryInt("take", 10)
	skip := c.QueryInt("skip", 0)
	search := c.Query("search", "")

	products, err := productHandler.productService.GetAllProduct(take, skip, search)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Bad Request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success fetch products",
		"data":    products,
	})
}

func (productHandler *ProductHandler) PostCreateProduct(c *fiber.Ctx) error {
	name := c.FormValue("name")
	desc := c.FormValue("description")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadGateway,
			"message": "Price must be a number",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	image, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadGateway,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if !pkg.IsImage(image) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadGateway,
			"message": "Bad request",
			"data": fiber.Map{
				"error": "File is not support, image only.",
			},
		})
	}

	os.MkdirAll("./images/product", os.ModePerm)

	savePath := filepath.Join("./images/product", image.Filename)

	if err := c.SaveFile(image, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	var product *models.Product = &models.Product{
		Name:        name,
		Description: desc,
		Price:       price,
		Image:       savePath,
	}

	productHandler.productService.CreateProduct(product)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "Success create product",
		"data":    product,
	})
}

func (productHandler *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	var product models.Product
	imageAvailable := true
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Bad request",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	name := c.FormValue("name")
	desc := c.FormValue("description")
	price, err := strconv.ParseFloat(c.FormValue("price"), 2)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadGateway,
			"message": "Price must be a number",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	product.Name = name
	product.Description = desc
	product.Price = price

	image, err := c.FormFile("image")

	if err != nil {
		imageAvailable = false
		product.Image = ""
	}

	if imageAvailable {
		if !pkg.IsImage(image) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    fiber.StatusBadGateway,
				"message": "Only accept image",
				"data": fiber.Map{
					"error": err.Error(),
				},
			})
		}

		os.MkdirAll("./images/product", os.ModePerm)

		savePath := filepath.Join("./images/product", image.Filename)

		if err := c.SaveFile(image, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    fiber.StatusInternalServerError,
				"message": "Something went wrong",
				"data": fiber.Map{
					"error": err.Error(),
				},
			})
		}

		product.Image = savePath
	}

	updatedProduct, err := productHandler.productService.UpdateProduct(uint(productId), &product)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err,
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    fiber.StatusCreated,
		"message": "Success update product",
		"data":    updatedProduct,
	})
}

func (productHandler *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	productId, err := c.ParamsInt("product_id")

	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err.Error(),
			},
		})
	}

	if err := productHandler.productService.DeleteProduct(uint(productId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Something went wrong",
			"data": fiber.Map{
				"error": err,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "Success delete product",
		"data":    nil,
	})
}
