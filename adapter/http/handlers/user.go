package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/fiber-micro/adapter/database/entities"
	"github.com/javiorfo/fiber-micro/application/port"
	"github.com/javiorfo/go-microservice-lib/pagination"
	"github.com/javiorfo/go-microservice-lib/response"
	"github.com/javiorfo/go-microservice-lib/tracing"
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
// @Security		OAuth2Password
func FindAllUsers(ds port.UserService) fiber.Handler {
	tracer := otel.Tracer(tracing.Name())
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()

		page, err := pagination.NewPage(
			c.Query("page", "1"),
			c.Query("size", "10"),
			c.Query("sortBy", "id"),
			c.Query("sortOrder", "asc"),
		)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(response.NewResponseError(span, response.Error{
					Code:    "FIBER-MICRO-005",
					Message: response.Message(err.Error()),
				}))
		}

		users, err := ds.FindAll(ctx, entities.NewUserFilter(*page,
			c.Get("username"),
			c.Get("permissions.name"),
			c.Get("createDate"),
		))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(response.InternalServerError(span, response.Message(err.Error())))
		}

		return c.JSON(response.NewPaginationResponse(pagination.Paginator(*page, len(users)), users))
	}
}
