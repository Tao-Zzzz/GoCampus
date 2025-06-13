package handler

import (
	"context"

	"github.com/campus-trading/user-service/internal/service"
	"github.com/campus-trading/user-service/pkg/proto"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *proto.RegisterRequest, rsp *proto.RegisterResponse) error {
	userID, err := h.svc.Register(ctx, req.Email, req.Phone, req.Password, req.Nickname)
	if err != nil {
		return err
	}
	rsp.UserId = userID
	rsp.Message = "Registration successful"
	return nil
}

func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest, rsp *proto.LoginResponse) error {
	token, userID, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		return err
	}
	rsp.Token = token
	rsp.UserId = userID
	return nil
}

func (h *UserHandler) UpdateProfile(ctx context.Context, req *proto.UpdateProfileRequest, rsp *proto.UpdateProfileResponse) error {
	if err := h.svc.UpdateProfile(ctx, req.UserId, req.Nickname, req.AvatarUrl, req.ContactInfo); err != nil {
		return err
	}
	rsp.Message = "Profile updated successfully"
	return nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *proto.GetProfileRequest, rsp *proto.GetProfileResponse) error {
	user, err := h.svc.GetProfile(ctx, req.UserId)
	if err != nil {
		return err
	}
	rsp.UserId = user.UserID
	rsp.Email = user.Email
	rsp.Phone = user.Phone
	rsp.Nickname = user.Nickname
	rsp.AvatarUrl = user.AvatarURL
	rsp.ContactInfo = user.ContactInfo
	rsp.CreatedAt = user.CreatedAt.String()
	return nil
}