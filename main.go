package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
)

//type Server struct {
//	store      *db.Store
//	tokenMaker token.Maker
//	router     *gin.Engine
//}
//
//func NewServer(store *db.Store) *Server {
//	tokenMaker, _ := token.NewJWTMaker("12345678901234567890123456789012")
//	server := &Server{
//		store:      store,
//		tokenMaker: tokenMaker,
//	}
//	router := gin.Default()
//	router.POST("/users", server.createUser)
//	router.POST("/users/login", server.loginUser)
//	authRoutes := router.Group("/").Use(api.AuthMiddleware(server.tokenMaker))
//
//	authRoutes.POST("/accounts", server.createAccount)
//	authRoutes.GET("/accounts/:id", server.getAccount)
//	authRoutes.GET("/accounts", server.listAccounts)
//	authRoutes.PUT("/accounts", server.updateAccount)
//	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
//	authRoutes.POST("/transfers", server.createTransfer)
//
//	server.router = router
//	return server
//}
//
//func (server *Server) Start(address string) error {
//	return server.router.Run(address)
//}
//
//// createAccountRequest определяет структуру входящего JSON-запроса для создания счета.
//type createAccountRequest struct {
//	Currency string `json:"currency" binding:"required,oneof=USD EUR RUB"`
//}
//
//// createAccount - это хендлер для эндпоинта POST /accounts.
//func (server *Server) createAccount(ctx *gin.Context) {
//	var req createAccountRequest
//	// Проверяем, соответствует ли входящий JSON нашей структуре и валидации.
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	authPayload := ctx.MustGet(api.AuthorizationPayloadKey).(*token.Payload)
//
//	// Готовим аргументы для вызова функции из db.
//	arg := db.CreateAccountParams{
//		Owner:    authPayload.Username,
//		Currency: req.Currency,
//		Balance:  0, // Баланс при создании счета всегда 0
//	}
//
//	// Вызываем сгенерированную sqlc функцию CreateAccount.
//	account, err := server.store.CreateAccount(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// Если все прошло успешно, возвращаем созданный аккаунт со статусом 200 OK.
//	ctx.JSON(http.StatusOK, account)
//}
//
//type getAccountRequest struct {
//	ID int64 `uri:"id" binding:"required,min=1"`
//}
//
//func (server *Server) getAccount(ctx *gin.Context) {
//	var req getAccountRequest
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	account, err := server.store.GetAccount(ctx, req.ID)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, account)
//}
//
//type listAccountsRequest struct {
//	PageID   int32 `form:"page_id" binding:"required,min=1"`
//	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
//}
//
//func (server *Server) listAccounts(ctx *gin.Context) {
//	var req listAccountsRequest
//	if err := ctx.ShouldBindQuery(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	arg := db.ListAccountsParams{
//		Limit:  req.PageSize,
//		Offset: (req.PageID - 1) * req.PageSize,
//	}
//
//	accounts, err := server.store.ListAccounts(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, accounts)
//}
//
//type updateAccountRequest struct {
//	ID      int64 `json:"id" binding:"required,min=1"`
//	Balance int64 `json:"balance" binding:"required,min=0"`
//}
//
//func (server *Server) updateAccount(ctx *gin.Context) {
//	var req updateAccountRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	arg := db.UpdateAccountParams{
//		ID:      req.ID,
//		Balance: req.Balance,
//	}
//
//	account, err := server.store.UpdateAccount(ctx, arg)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	ctx.JSON(http.StatusOK, account)
//}
//
//type deleteAccountRequest struct {
//	ID int64 `uri:"id" binding:"required,min=1"`
//}
//
//func (server *Server) deleteAccount(ctx *gin.Context) {
//	var req deleteAccountRequest
//	if err := ctx.ShouldBindUri(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	err := server.store.DeleteAccount(ctx, req.ID)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// При успешном удалении принято возвращать статус 204 No Content
//	ctx.Status(http.StatusNoContent)
//}
//
//type transferRequest struct {
//	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
//	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
//	Amount        int64 `json:"amount" binding:"required,gt=0"`
//}
//
//func (server *Server) createTransfer(ctx *gin.Context) {
//	var req transferRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 2. Проверяем, что аккаунт-отправитель существует и на нём достаточно денег.
//	fromAccount, err := server.store.GetAccount(ctx, req.FromAccountID)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": "from_account_not_found"})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 2. Получаем из БД счет получателя
//	toAccount, err := server.store.GetAccount(ctx, req.ToAccountID)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": "to_account_not_found"})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	// 3. ПРОВЕРЯЕМ ВАЛЮТУ!
//	if fromAccount.Currency != toAccount.Currency {
//		err := fmt.Errorf("cannot transfer between different currencies: %s vs %s", fromAccount.Currency, toAccount.Currency)
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	// Важная бизнес-логика: проверяем баланс
//	if fromAccount.Balance < req.Amount {
//		err := fmt.Errorf("insufficient balance on account %d", fromAccount.ID)
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	arg := db.TransferTxParams{
//		FromAccountID: req.FromAccountID,
//		ToAccountID:   req.ToAccountID,
//		Amount:        req.Amount,
//	}
//
//	result, err := server.store.TransferTx(ctx, arg)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, result)
//}
//
//type createUserRequest struct {
//	Username string `json:"username" binding:"required,alphanum"`
//	Password string `json:"password" binding:"required,min=6"`
//	FullName string `json:"full_name" binding:"required"`
//	Email    string `json:"email" binding:"required,email"`
//}
//
//type createUserResponse struct {
//	Username          string    `json:"username"`
//	FullName          string    `json:"full_name"`
//	Email             string    `json:"email"`
//	PasswordChangedAt time.Time `json:"password_changed_at"`
//	CreatedAt         time.Time `json:"created_at"`
//}
//
//func (server *Server) createUser(ctx *gin.Context) {
//	var req createUserRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	hashedPassword, err := util.HashPassword(req.Password)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	arg := db.CreateUserParams{
//		Username:       req.Username,
//		HashedPassword: hashedPassword,
//		FullName:       req.FullName,
//		Email:          req.Email,
//	}
//
//	user, err := server.store.CreateUser(ctx, arg)
//	if err != nil {
//		if pgErr, ok := err.(*pgconn.PgError); ok {
//			switch pgErr.Code {
//			case "23505":
//				ctx.JSON(http.StatusForbidden, gin.H{"error": "username or email already exists"})
//				return
//			}
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//	}
//
//	rsp := createUserResponse{
//		Username:          user.Username,
//		FullName:          user.FullName,
//		Email:             user.Email,
//		PasswordChangedAt: user.PasswordChangedAt.Time,
//		CreatedAt:         user.CreatedAt.Time,
//	}
//	ctx.JSON(http.StatusOK, rsp)
//}
//
//type loginUserRequest struct {
//	Username string `json:"username" binding:"required,alphanum"`
//	Password string `json:"password" binding:"required,min=6"`
//}
//
//type loginUserResponse struct {
//	AccessToken string             `json:"access_token"`
//	User        createUserResponse `json:"user"`
//}
//
//func (server *Server) loginUser(ctx *gin.Context) {
//	var req loginUserRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	user, err := server.store.GetUser(ctx, req.Username)
//	if err != nil {
//		if err == pgx.ErrNoRows {
//			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
//			return
//		}
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	err = util.CheckPassword(req.Password, user.HashedPassword)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	accessToken, err := server.tokenMaker.CreateToken(
//		user.Username,
//		time.Minute*15,
//	)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	rsp := loginUserResponse{
//		AccessToken: accessToken,
//		User: createUserResponse{
//			Username:          user.Username,
//			FullName:          user.FullName,
//			Email:             user.Email,
//			PasswordChangedAt: user.PasswordChangedAt.Time,
//			CreatedAt:         user.CreatedAt.Time,
//		},
//	}
//	ctx.JSON(http.StatusOK, rsp)
//}

func main() {
	// Строка подключения к БД
	config, err := util.LoadConfig(".") // "." means current folder
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Create a connection pool using the config value
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)
	server := api.NewServer(store) // Assuming NewServer is in the 'api' package now

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
