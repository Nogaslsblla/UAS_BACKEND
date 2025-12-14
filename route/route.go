API := app.Group("/api/v1")

	API.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "API v1 Running"})
	})

authRepo := repository.NewAuthRepository(db)
authService := service.NewAuthService(authRepo)
	


       auth := API.Group("/auth")
	auth.Post("/login", authService.LoginEndpoint)
	auth.Post("/refresh", middleware.AuthMiddleware(""), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Refresh endpoint - to be implemented"})
	})
	auth.Post("/logout", middleware.AuthMiddleware(""), authService.LogoutEndpoint)
	auth.Get("/profile", middleware.AuthMiddleware(""), authService.ProfileEndpoint)


userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo)
	
	// Users Routes (Admin only)
		users := API.Group("/users")
		users.Use(middleware.AuthMiddleware("manage_users"))
		users.Get("/", userService.GetUsersEndpoint)
		users.Get("/:id", userService.GetUserByIDEndpoint)
		users.Post("/", userService.CreateUserEndpoint)
		users.Put("/:id", userService.UpdateUserEndpoint)
		users.Delete("/:id", userService.DeleteUserEndpoint)
		users.Put("/:id/role", userService.UpdateUserRoleEndpoint)