package sellerhandler

import (
	sellerrequest "marketly-app/internal/dto/request/seller_request"
	sellerresponse "marketly-app/internal/dto/response/seller_response"
	sellerservice "marketly-app/internal/services/seller_service"
	errorresponse "marketly-app/pkg/constant/error_response"
	"marketly-app/pkg/constant/response"
	"marketly-app/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type SellerHandler struct {
	sellerService sellerservice.ISellerService
}

func NewSellerHandler(sellerService sellerservice.ISellerService) *SellerHandler {
	return &SellerHandler{sellerService: sellerService}
}

func (h *SellerHandler) RegisterSeller(c echo.Context) error {
	var req sellerrequest.RegisterSellerRequest

	req.Username = c.FormValue("username")
	req.Name = c.FormValue("name")
	req.Email = c.FormValue("email")
	req.Password = c.FormValue("password")
	req.PhoneNumber = c.FormValue("phone_number")
	req.Address = c.FormValue("address")
	req.StoreName = c.FormValue("store_name")
	req.StoreDescription = c.FormValue("store_description")

	photoFile, err := c.FormFile("photo_file")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Photo file is required", err.Error())
	}
	req.PhotoFile = photoFile

	err = h.sellerService.Register(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to create seller")
	}

	return response.Success(c, http.StatusCreated, "Seller Created Successfully", nil)
}

func (h *SellerHandler) LoginSeller(c echo.Context) error {
	var req sellerrequest.LoginSellerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "bad request", err.Error())
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Password == "" {
		return response.Error(c, http.StatusBadRequest, "Password wajib diisi", nil)
	}

	if req.Email == "" && req.Username == "" {
		return response.Error(c, http.StatusBadRequest, "Email atau Username wajib diisi", nil)
	}

	token, err := h.sellerService.Login(c.Request().Context(), req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid login seller")
	}

	return response.Success(c, http.StatusOK, "Login Successfully", map[string]interface{}{
		"access_token": token,
	})
}

func (h *SellerHandler) GetProfileSeller(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	seller, ok := sellerToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := seller.Claims.(jwt.MapClaims)
	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	me, err := h.sellerService.GetProfile(c.Request().Context(), sellerID)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "invalid to get profile")
	}

	sellerResponse := sellerresponse.ToSellerResponse(*me)
	return response.Success(c, http.StatusOK, "Get Profile Successfully", sellerResponse)
}

func (h *SellerHandler) GetAllSeller(c echo.Context) error {
	pageInt, limitInt := utils.ParsePaginationParams(c, 10)
	search := c.QueryParam("search")

	sellers, total, err := h.sellerService.GetAllSeller(c.Request().Context(), pageInt, limitInt, search)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "failed to get sellers")
	}

	meta := utils.BuildPaginationMeta(c, pageInt, limitInt, int(total))
	data := make([]sellerresponse.SellerResponse, len(sellers))
	for i, u := range sellers {
		data[i] = sellerresponse.ToSellerResponse(*u)
	}

	return response.PaginatedSuccess(c, http.StatusOK, "Get All Seller Successfully", data, meta)
}

func (h *SellerHandler) UpdateProfileSeller(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token := sellerToken.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	var req sellerrequest.UpdateSellerRequest
	req.Username = c.FormValue("username")
	req.Name = c.FormValue("name")
	req.PhoneNumber = c.FormValue("phone_number")
	req.Address = c.FormValue("address")
	req.StoreName = c.FormValue("store_name")
	req.StoreDescription = c.FormValue("store_description")

	if err := h.sellerService.UpdateProfile(c.Request().Context(), sellerID, req); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "Failed to update profile")
	}

	return response.Success(c, http.StatusOK, "Profile updated successfully", nil)
}

func (h *SellerHandler) UpdatePhotoProfile(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token := sellerToken.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	photoFile, err := c.FormFile("photo_file")
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Photo file is required", err.Error())
	}

	req := sellerrequest.UpdatePhotoSellerRequest{
		PhotoFile: photoFile,
	}

	err = h.sellerService.UpdatePhoto(c.Request().Context(), sellerID, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "Failed to update photo")
	}

	return response.Success(c, http.StatusOK, "Photo updated successfully", nil)
}

func (h *SellerHandler) ChangePassword(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token := sellerToken.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	var req sellerrequest.ChangePasswordSellerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	err := h.sellerService.ChangePassword(c.Request().Context(), sellerID, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "Failed to change password")
	}

	return response.Success(c, http.StatusOK, "Password updated successfully", nil)
}

func (h *SellerHandler) ChangeEmail(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
	}

	token := sellerToken.(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	var req sellerrequest.ChangeEmailSellerRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	err := h.sellerService.ChangeEmail(c.Request().Context(), sellerID, req)
	if err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "Failed to change email")
	}

	return response.Success(c, http.StatusOK, "Email updated successfully", nil)
}

func (a *SellerHandler) LogoutSeller(c echo.Context) error {
	sellerToken := c.Get("seller")
	if sellerToken == nil {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	seller, ok := sellerToken.(*jwt.Token)
	if !ok {
		return response.Error(c, http.StatusUnauthorized, "Unauthorized: token invalid or expired", nil)
	}

	claims := seller.Claims.(jwt.MapClaims)
	sellerIDStr := claims["user_id"].(string)
	sellerID, _ := strconv.Atoi(sellerIDStr)

	if err := a.sellerService.Logout(c.Request().Context(), sellerID); err != nil {
		if customErr, ok := errorresponse.AsCustomErr(err); ok {
			return response.Error(c, customErr.Status, customErr.Msg, customErr.Err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, err.Error(), "Failed to logout seller")
	}

	return response.Success(c, http.StatusOK, "Logout Successfully", nil)
}

