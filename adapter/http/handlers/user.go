package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/adapter/http/request"
	srvResponse "github.com/javiorfo/fiber-micro/adapter/http/response"
	"github.com/javiorfo/fiber-micro/application/domain/model"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/response"
	"github.com/javiorfo/go-microservice-lib/security"
	"github.com/javiorfo/go-microservice-lib/tracing"
	"github.com/javiorfo/go-microservice-lib/validation"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"go.opentelemetry.io/otel"
)

// @Summary		List all users
// @Description	Get a list of users with pagination
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			page				query		int										false	"Page number"
// @Param			size				query		int										false	"Size per page"
// @Param			sortBy				query		string									false	"Sort by field"
// @Param			sortOrder			query		string									false	"Sort order (asc or desc)"
// @Param			username			header		string									false	"Username filter"
// @Param			permissions.name	header		string									false	"Permission Name filter"
// @Param			createDate			header		string									false	"Creation date filter"
// @Success		200					{object}	response.PaginationResponse[model.User]	"Paginated list of dummies"
// @Failure		400					{object}	response.ResponseError					"Invalid query parameters"
// @Failure		500					{object}	response.ResponseError					"Internal server error"
// @Router			/users [get]
// @Security		BearerAuth
func FindAllUsers(service port.UserService) fiber.Handler {
	tracer := otel.Tracer(tracing.Name())
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()

		pn := c.Query("page", "1")
		ps := c.Query("size", "10")
		pageRequest, err := pagination.PageRequestFrom(
			pn,
			ps,
			pagination.WithSortOrder(
				c.Query("sortBy", "id"),
				sort.DirectionFromString(c.Query("sortOrder", "asc")),
			),
			pagination.WithFilter(
				entities.NewUserFilter(
					c.Get("username"),
					c.Get("permissions.name"),
					c.Get("createDate"),
				),
			),
		)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(response.NewResponseError(span, response.Error{
					Code:    "FIBER-MICRO-005",
					Message: response.Message(err.Error()),
				}))
		}

		page, err := service.FindAll(ctx, pageRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(response.InternalServerError(span, response.Message(err.Error())))
		}

		return c.Status(fiber.StatusOK).
			JSON(response.NewPaginationResponse(pn, ps, page.Total, page.Elements))
	}
}

// @Summary		Create a new user
// @Description	Create a new user with the provided information
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			dummy	body		request.CreateUserRequest	true	"Dummy information"
// @Success		201		{object}	response.CreateUserResponse
// @Failure		400		{object}	response.ResponseError	"Invalid request body or validation errors"
// @Failure		500		{object}	response.ResponseError	"Internal server error"
// @Router			/users [post]
// @Security		BearerAuth
func CreateUser(service port.UserService) fiber.Handler {
	tracer := otel.Tracer(tracing.Name())
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()

		userRequest, errResp := validation.ValidateRequest[request.CreateUserRequest](c, span,
			"FIBER-MICRO-006",
			request.ValidateUserNotBlank("FIBER-MICRO-007"),
			model.ValidateStatus("FIBER-MICRO-008"),
		)

		if errResp != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}

		log.Info(tracing.LogInfo(span, fmt.Sprintf("Received user: %+v", userRequest)))

		user := userRequest.Into(security.GetTokenUsername(c))
		if err := service.Create(ctx, &user, userRequest.Permission); err != nil {
			return err.ToResponse(c)
		}

		return c.Status(fiber.StatusCreated).JSON(srvResponse.CreateUserResponse{User: user})
	}
}

// @Summary		Login user
// @Description	Login a user and return a JWT token
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			user	body		request.LoginRequest	true	"Username and password"
// @Success		201		{object}	response.LoginResponse
// @Failure		400		{object}	response.ResponseError	"Invalid request body or validation errors"
// @Failure		500		{object}	response.ResponseError	"Internal server error"
// @Router			/users/login [post]
// @Security		BearerAuth
func Login(service port.UserService) fiber.Handler {
	tracer := otel.Tracer(tracing.Name())
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()

		loginReq, errResp := validation.ValidateRequest[request.LoginRequest](c, span,
			"FIBER-MICRO-009",
			request.ValidateUserNotBlank("FIBER-MICRO-010"),
		)
		if errResp != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errResp)
		}

		log.Infof(tracing.LogInfo(span, fmt.Sprintf("Received credentials: %+v", loginReq)))

		token, err := service.Login(ctx, loginReq.Username, loginReq.Password)

		if err != nil {
			return err.ToResponse(c)
		}

		return c.Status(fiber.StatusOK).JSON(srvResponse.NewLoginResponse(loginReq.Username, token))
	}
}
