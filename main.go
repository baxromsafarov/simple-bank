package main

import (
	db "github.com/baxromsafarov/simple-bank/sqlc" // Замени на свой путь
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// createAccountRequest определяет структуру входящего JSON-запроса для создания счета.
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR RUB"`
}

// createAccount - это хендлер для эндпоинта POST /accounts.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	// Проверяем, соответствует ли входящий JSON нашей структуре и валидации.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Готовим аргументы для вызова функции из db.
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0, // Баланс при создании счета всегда 0
	}

	// Вызываем сгенерированную sqlc функцию CreateAccount.
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Если все прошло успешно, возвращаем созданный аккаунт со статусом 200 OK.
	ctx.JSON(http.StatusOK, account)
}

func main() {
	// Строка подключения к БД
	connStr := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"

	// Устанавливаем соединение с базой данных
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn)
	server := NewServer(store)

	// Запускаем сервер на порту 8080
	err = server.Start("0.0.0.0:8080")
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
