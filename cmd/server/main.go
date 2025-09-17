package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/danielgalindoj/auth-service/internal/adapters/database/postgres"
	"github.com/danielgalindoj/auth-service/internal/config"
	"github.com/danielgalindoj/auth-service/internal/core/services"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No se pudo cargar archivo .env, usando variables del sistema")
	}

	// Cargar configuración
	cfg := config.Load()

	// Configurar logger
	setupLogger(cfg.LogLevel, cfg.LogFormat)

	logrus.WithFields(logrus.Fields{
		"port": cfg.Port,
		"env":  cfg.Environment,
	}).Info("Iniciando servidor Auth Service...")

	// Conectar a base de datos
	db, err := postgres.NewConnection(cfg)
	if err != nil {
		logrus.Fatal("Error conectando a la base de datos: ", err)
	}
	defer db.Close()

	logrus.Info("✅ Conexión a PostgreSQL establecida")

	// Inicializar repositorios
	userRepo := postgres.NewUserRepository(db)

	// Inicializar servicios
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)

	// Crear datos de prueba
	if err := createTestData(userService, authService); err != nil {
		logrus.Warn("Error creando datos de prueba: ", err)
	}

	// Configurar rutas
	// router := routes.SetupRoutes(authService, userService, cfg)

	// // CORS
	// corsHandler := handlers.CORS(
	// 	handlers.AllowedOrigins([]string{"*"}),
	// 	handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	// 	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	// )(router)

	logrus.Info("🚀 Servidor corriendo en http://localhost:", cfg.Port)
	logrus.Info("📝 Usuarios de prueba:")
	logrus.Info("   admin@example.com / admin123")
	logrus.Info("   john@example.com / john123")
	logrus.Info("   jane@example.com / jane123")

	// if err := http.ListenAndServe(":"+cfg.Port, corsHandler); err != nil {
	// 	logrus.Fatal("Error iniciando servidor: ", err)
	// }
}

func setupLogger(level, format string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	if format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func createTestData(userService *services.UserService, authService *services.AuthService) error {
	testUsers := []struct {
		Email     string
		Password  string
		FirstName string
		LastName  string
		Username  string
	}{
		{"admin@example.com", "admin123", "Admin", "User", "admin"},
		{"john@example.com", "john123", "John", "Doe", "john"},
		{"jane@example.com", "jane123", "Jane", "Smith", "jane"},
		{"bob@example.com", "bob123", "Bob", "Johnson", "bob"},
		{"alice@example.com", "alice123", "Alice", "Wilson", "alice"},
	}

	for _, userData := range testUsers {
		// Verificar si el usuario ya existe
		if _, err := userService.GetByEmail(userData.Email); err == nil {
			continue // Usuario ya existe, saltar
		}

		// Crear usuario
		user, err := authService.Register(userData.Email, userData.Password, userData.FirstName, userData.LastName, userData.Username)
		if err != nil {
			logrus.Warn("Error creando usuario de prueba ", userData.Email, ": ", err)
			continue
		}

		logrus.Info("✅ Usuario de prueba creado: ", user.Email)
	}

	return nil
}
