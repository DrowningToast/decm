package v1

import (
	"context"
	"time"

	"decm-backend/internal/handlers"
	generated "decm-database/go/generated"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email               string `json:"email" validate:"required,email"`
	Username            string `json:"username" validate:"required,min=3,max=50"`
	FirstName           string `json:"first_name" validate:"required,min=1,max=100"`
	LastName            string `json:"last_name" validate:"required,min=1,max=100"`
	AcademicInstitution string `json:"academic_institution,omitempty"`
	AcademicEmail       string `json:"academic_email,omitempty"`
	Bio                 string `json:"bio,omitempty"`
	ProfilePictureURL   string `json:"profile_picture_url,omitempty"`
	WalletAddress       string `json:"wallet_address,omitempty"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	FirstName           *string `json:"first_name,omitempty"`
	LastName            *string `json:"last_name,omitempty"`
	Bio                 *string `json:"bio,omitempty"`
	ProfilePictureURL   *string `json:"profile_picture_url,omitempty"`
	AcademicInstitution *string `json:"academic_institution,omitempty"`
	AcademicEmail       *string `json:"academic_email,omitempty"`
	WalletAddress       *string `json:"wallet_address,omitempty"`
}

// SearchUsersRequest represents query parameters for searching users
type SearchUsersRequest struct {
	Query  string `query:"q"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

// ListUsersRequest represents query parameters for listing users
type ListUsersRequest struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

// UserRoutes sets up user-related routes
func UserRoutes(app fiber.Router, h *handlers.Handlers) {
	users := app.Group("/users")

	users.Post("/", createUser(h))
	users.Get("/:id", getUserByID(h))
	users.Put("/:id", updateUser(h))
	users.Delete("/:id", deleteUser(h))
	users.Get("/", listUsers(h))
	users.Get("/search", searchUsers(h))
	users.Get("/email/:email", getUserByEmail(h))
	users.Get("/username/:username", getUserByUsername(h))
}

// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User data"
// @Success 201 {object} interface{} "User created successfully"
// @Failure 400 {object} interface{} "Invalid request"
// @Failure 409 {object} interface{} "User already exists"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users [post]
func createUser(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		// Create user using generated queries
		user, err := h.DB.DB.Queries.CreateUser(ctx, generated.CreateUserParams{
			Email:               req.Email,
			Username:            req.Username,
			FirstName:           req.FirstName,
			LastName:            req.LastName,
			AcademicInstitution: pgtype.Text{String: req.AcademicInstitution, Valid: req.AcademicInstitution != ""},
			AcademicEmail:       pgtype.Text{String: req.AcademicEmail, Valid: req.AcademicEmail != ""},
			Bio:                 pgtype.Text{String: req.Bio, Valid: req.Bio != ""},
			ProfilePictureUrl:   pgtype.Text{String: req.ProfilePictureURL, Valid: req.ProfilePictureURL != ""},
			WalletAddress:       pgtype.Text{String: req.WalletAddress, Valid: req.WalletAddress != ""},
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to create user",
			})
		}

		return c.Status(201).JSON(fiber.Map{
			"message": "User created successfully",
			"user":    user,
		})
	}
}

// @Summary Get user by ID
// @Description Get user information by user ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} interface{} "User information"
// @Failure 400 {object} interface{} "Invalid user ID"
// @Failure 404 {object} interface{} "User not found"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/{id} [get]
func getUserByID(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		user, err := h.DB.DB.Queries.GetUserByID(ctx, id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"user": user,
		})
	}
}

// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserRequest true "Updated user data"
// @Success 200 {object} interface{} "User updated successfully"
// @Failure 400 {object} interface{} "Invalid request"
// @Failure 404 {object} interface{} "User not found"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/{id} [put]
func updateUser(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		// Convert optional strings - FirstName and LastName are required strings in UpdateUser
		var firstName, lastName string
		var bio, profilePictureUrl, academicInstitution, academicEmail, walletAddress pgtype.Text

		// Get current user first to use existing values if not provided
		currentUser, err := h.DB.DB.Queries.GetUserByID(ctx, id)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Set firstName and lastName (required strings)
		firstName = currentUser.FirstName
		lastName = currentUser.LastName

		if req.FirstName != nil {
			firstName = *req.FirstName
		}
		if req.LastName != nil {
			lastName = *req.LastName
		}

		if req.Bio != nil {
			bio = pgtype.Text{String: *req.Bio, Valid: true}
		}
		if req.ProfilePictureURL != nil {
			profilePictureUrl = pgtype.Text{String: *req.ProfilePictureURL, Valid: true}
		}
		if req.AcademicInstitution != nil {
			academicInstitution = pgtype.Text{String: *req.AcademicInstitution, Valid: true}
		}
		if req.AcademicEmail != nil {
			academicEmail = pgtype.Text{String: *req.AcademicEmail, Valid: true}
		}
		if req.WalletAddress != nil {
			walletAddress = pgtype.Text{String: *req.WalletAddress, Valid: true}
		}

		user, err := h.DB.DB.Queries.UpdateUser(ctx, generated.UpdateUserParams{
			ID:                  id,
			FirstName:           firstName,
			LastName:            lastName,
			Bio:                 bio,
			ProfilePictureUrl:   profilePictureUrl,
			AcademicInstitution: academicInstitution,
			AcademicEmail:       academicEmail,
			WalletAddress:       walletAddress,
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update user",
			})
		}

		return c.JSON(fiber.Map{
			"message": "User updated successfully",
			"user":    user,
		})
	}
}

// @Summary Delete user
// @Description Delete a user account
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {object} interface{} "User deleted successfully"
// @Failure 400 {object} interface{} "Invalid user ID"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/{id} [delete]
func deleteUser(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		err = h.DB.DB.Queries.DeleteUser(ctx, id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to delete user",
			})
		}

		return c.JSON(fiber.Map{
			"message": "User deleted successfully",
		})
	}
}

// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Produce json
// @Param limit query int false "Number of users to return (default: 10, max: 100)"
// @Param offset query int false "Number of users to skip (default: 0)"
// @Success 200 {object} interface{} "List of users"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users [get]
func listUsers(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req ListUsersRequest
		if err := c.QueryParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid query parameters",
			})
		}

		// Set defaults
		if req.Limit <= 0 || req.Limit > 100 {
			req.Limit = 10
		}
		if req.Offset < 0 {
			req.Offset = 0
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		users, err := h.DB.DB.Queries.ListUsers(ctx, generated.ListUsersParams{
			Limit:  int32(req.Limit),
			Offset: int32(req.Offset),
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to list users",
			})
		}

		return c.JSON(fiber.Map{
			"users":  users,
			"limit":  req.Limit,
			"offset": req.Offset,
		})
	}
}

// @Summary Search users
// @Description Search users by name, username, or email
// @Tags users
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Number of users to return (default: 10, max: 100)"
// @Param offset query int false "Number of users to skip (default: 0)"
// @Success 200 {object} interface{} "Search results"
// @Failure 400 {object} interface{} "Missing search query"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/search [get]
func searchUsers(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SearchUsersRequest
		if err := c.QueryParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid query parameters",
			})
		}

		if req.Query == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Search query is required",
			})
		}

		// Set defaults
		if req.Limit <= 0 || req.Limit > 100 {
			req.Limit = 10
		}
		if req.Offset < 0 {
			req.Offset = 0
		}

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		users, err := h.DB.DB.Queries.SearchUsers(ctx, generated.SearchUsersParams{
			PlaintoTsquery: req.Query,
			Limit:          int32(req.Limit),
			Offset:         int32(req.Offset),
		})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to search users",
			})
		}

		return c.JSON(fiber.Map{
			"users":  users,
			"query":  req.Query,
			"limit":  req.Limit,
			"offset": req.Offset,
		})
	}
}

// @Summary Get user by email
// @Description Get user information by email address
// @Tags users
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} interface{} "User information"
// @Failure 404 {object} interface{} "User not found"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/email/{email} [get]
func getUserByEmail(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Params("email")

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		user, err := h.DB.DB.Queries.GetUserByEmail(ctx, email)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"user": user,
		})
	}
}

// @Summary Get user by username
// @Description Get user information by username
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} interface{} "User information"
// @Failure 404 {object} interface{} "User not found"
// @Failure 500 {object} interface{} "Internal server error"
// @Router /api/v1/users/username/{username} [get]
func getUserByUsername(h *handlers.Handlers) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")

		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		user, err := h.DB.DB.Queries.GetUserByUsername(ctx, username)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		return c.JSON(fiber.Map{
			"user": user,
		})
	}
}
