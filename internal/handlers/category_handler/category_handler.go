package categoryhandler

import (
	categoryrequest "marketly-app/internal/dto/request/category_request"
	categoryresponse "marketly-app/internal/dto/response/category_response"
	categoryservice "marketly-app/internal/services/category_service"
	errorresponse "marketly-app/pkg/constant/error_response"
	"marketly-app/pkg/constant/response"
	"marketly-app/pkg/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService categoryservice.ICategoryService
}

func NewCategoryHandler(categoryService categoryservice.ICategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (a *CategoryHandler) CreateCategory(c echo.Context) error {
	var req categoryrequest.CreateCategoryRequest
	req.Name = c.FormValue("name")

	err := a.categoryService.CreateCategory(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create category")
	}

	return response.Success(c, http.StatusOK, "Category Created Successfully", nil)
}

func (a *CategoryHandler) GetAllCategory(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	categories, total, err := a.categoryService.GetAllCategory(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get category")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, total)
	data := make([]categoryresponse.CategoryResponse, len(categories))
	for i, s := range categories {
		data[i] = categoryresponse.ToCategoryResponse(*s)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Category Successfully", data, meta)
}

func (a *CategoryHandler) GetByIdCategory(c echo.Context) error {
	categoryIdStr := c.Param("categoryId")

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}

	category, err := a.categoryService.GetByIdCategory(c.Request().Context(), categoryId)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get category")
	}

	res := categoryresponse.ToCategoryResponse(*category)

	return response.Success(c, http.StatusOK, "Get Category Successfully", res)
}

func (a *CategoryHandler) UpdateCategory(c echo.Context) error {
	categoryIdStr := c.Param("categoryId")

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}

	var req categoryrequest.UpdateCategoryRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	err = a.categoryService.UpdateCategory(c.Request().Context(), categoryId, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to update category")
	}

	return response.Success(c, http.StatusOK, "Category Updated Successfully", nil)
}

func (a *CategoryHandler) DeleteCategory(c echo.Context) error {
	categoryIdStr := c.Param("categoryId")

	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid category ID", err.Error())
	}

	if err := a.categoryService.DeleteCategory(c.Request().Context(), categoryId); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to delete category")
	}

	return response.Success(c, http.StatusOK, "Category Deleted Successfully", nil)
}
