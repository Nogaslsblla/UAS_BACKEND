package route

import (
	"database/sql"
	"uas_backend/app/repository"
	"uas_backend/app/service"
	"uas_backend/middleware"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *fiber.App, db *sql.DB, mongoColl *mongo.Collection) {
	// API v1 group
	API := app.Group("/api/v1")

	// Health check
	API.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "API v1 Running"})
	})

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	achievementRepo := repository.NewAchievementRepository(db, mongoColl)

	// Initialize services
	authService := service.NewAuthService(authRepo)
	userService := service.NewUserService(userRepo)
	achievementService := service.NewAchievementService(achievementRepo)

	// Authentication Routes
	auth := API.Group("/auth")
	auth.Post("/login", authService.LoginEndpoint)
	auth.Post("/refresh", middleware.AuthMiddleware(""), authService.RefreshTokenEndpoint)
	auth.Post("/logout", middleware.AuthMiddleware(""), authService.LogoutEndpoint)
	auth.Get("/profile", middleware.AuthMiddleware(""), authService.ProfileEndpoint)

	// Users Routes (Admin only)
	users := API.Group("/users")
	users.Use(middleware.AuthMiddleware("manage_users"))
	users.Get("/", userService.GetUsersEndpoint)
	users.Get("/:id", userService.GetUserByIDEndpoint)
	users.Post("/", userService.CreateUserEndpoint)
	users.Put("/:id", userService.UpdateUserEndpoint)
	users.Delete("/:id", userService.DeleteUserEndpoint)
	users.Put("/:id/role", userService.UpdateUserRoleEndpoint)

	// Achievements Routes
	achievements := API.Group("/achievements")
	achievements.Use(middleware.AuthMiddleware(""))
	achievements.Get("/", achievementService.GetAchievementsEndpoint)

	achievements.Get("/:id", achievementService.GetAchievementByIDEndpoint)
	achievements.Post("/", achievementService.CreateAchievementEndpoint)
	achievements.Put("/:id", achievementService.UpdateAchievementEndpoint)
	achievements.Delete("/:id", achievementService.DeleteAchievementEndpoint)

	achievements.Post("/:id/submit", achievementService.SubmitAchievementEndpoint)

	achievements.Post("/:id/verify", achievementService.VerifyAchievementEndpoint)

	achievements.Post("/:id/reject", achievementService.RejectAchievementEndpoint)

	achievements.Get("/:id/history", achievementService.GetAchievementHistoryEndpoint)

	achievements.Post("/:id/attachments", achievementService.UploadAttachmentEndpoint)
	achievements.Get("/statistics", achievementService.GetAchievementStatisticsEndpoint)

	// Students Routes
	students := API.Group("/students")
	students.Use(middleware.AuthMiddleware(""))
	students.Get("/", userService.GetStudentsEndpoint)
	students.Get("/:id", userService.GetStudentByIDEndpoint)
	students.Get("/:id/achievements", userService.GetStudentAchievementsEndpoint)
	students.Put("/:id/advisor", middleware.AuthMiddleware("manage_users"), userService.UpdateStudentAdvisorEndpoint)

	// Lecturers Routes
	lecturers := API.Group("/lecturers")
	lecturers.Use(middleware.AuthMiddleware(""))
	lecturers.Get("/", userService.GetLecturersEndpoint)
	lecturers.Get("/:id/advisees", userService.GetLecturerAdviseesEndpoint)

	// Reports & Analytics Routes
	reports := API.Group("/reports")
	reports.Use(middleware.AuthMiddleware(""))
	reports.Get("/statistics", achievementService.GetReportsStatisticsEndpoint)
	reports.Get("/student/:id", achievementService.GetStudentReportEndpoint)

	// Admin Routes
	admin := API.Group("/admin")
	admin.Use(middleware.AuthMiddleware("manage_users"))
	admin.Get("/achievements", achievementService.GetAllAchievementsForAdminEndpoint)
}