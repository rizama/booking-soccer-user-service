package controllers

import (
	"fmt"
	"net/http"
	errWrap "user-service/common/error"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserController adalah struct yang menangani semua HTTP request terkait user
// Menggunakan Dependency Injection pattern dengan menerima service registry
type UserController struct {
	service services.IServiceRegistry // Interface untuk mengakses semua service
}

// IUserController adalah interface yang mendefinisikan kontrak untuk user controller
// Menggunakan Interface Segregation Principle - hanya method yang diperlukan
type IUserController interface {
	Login(*gin.Context)         // Endpoint untuk login user
	Register(*gin.Context)      // Endpoint untuk registrasi user baru
	Update(*gin.Context)        // Endpoint untuk update data user
	GetUserLogin(*gin.Context)  // Endpoint untuk mendapatkan data user yang sedang login
	GetUserByUUID(*gin.Context) // Endpoint untuk mendapatkan user berdasarkan UUID
}

// NewUserController adalah constructor function untuk membuat instance UserController
// Menggunakan Dependency Injection - service di-inject dari luar
func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{
		service: service, // Menyimpan service registry untuk digunakan di method-method
	}
}

// Login menangani HTTP POST request untuk autentikasi user
// Alur: Request Binding -> Validation -> Business Logic -> Response
func (u *UserController) Login(ctx *gin.Context) {
	// Step 1: Inisialisasi struct untuk menampung data request
	request := &dto.LoginRequest{}

	// Step 2: Binding JSON dari request body ke struct
	// ShouldBindJSON otomatis parse JSON dan mapping ke field struct
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		// Jika binding gagal (JSON invalid/format salah), return error 400
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest, // 400 - Bad Request
			Err:  err,
			Gin:  ctx,
		})
		return // Early return untuk menghentikan eksekusi
	}

	// Step 3: Validasi data menggunakan validator
	// Cek apakah field required sudah diisi, format email benar, dll
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		// Jika validasi gagal, return error 422 dengan detail error
		errMesage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err) // Format error menjadi user-friendly
		response.HttpRresponse(response.ParamHttpResp{
			Code:    http.StatusUnprocessableEntity, // 422 - Unprocessable Entity
			Message: &errMesage,
			Data:    errResponse, // Detail field mana yang error
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// fmt.Println(request)

	// Step 4: Panggil business logic di service layer
	// Controller hanya sebagai penghubung, logic ada di service
	user, err := u.service.GetUserService().Login(ctx, request)
	if err != nil {
		// Jika login gagal (username/password salah), return error 400
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest, // 400 - Bad Request
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// fmt.Println(user)

	// Step 5: Return response sukses dengan data user dan token
	response.HttpRresponse(response.ParamHttpResp{
		Code:  http.StatusOK, // 200 - Success
		Data:  user.User,     // Data user (tanpa password)
		Token: &user.Token,   // JWT token untuk autentikasi selanjutnya
		Gin:   ctx,
	})
}

// Register menangani HTTP POST request untuk registrasi user baru
// Alur sama dengan Login: Binding -> Validation -> Business Logic -> Response
func (u *UserController) Register(ctx *gin.Context) {
	// Step 1: Inisialisasi struct untuk data registrasi
	request := &dto.RegisterRequest{}

	// Step 2: Binding JSON request ke struct
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Step 3: Validasi data registrasi
	// Cek required fields, format email, password confirmation, dll
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMesage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpRresponse(response.ParamHttpResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMesage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// Step 4: Panggil service untuk proses registrasi
	// Service akan handle: hash password, cek duplikasi, save ke DB
	user, err := u.service.GetUserService().Register(ctx, request)
	if err != nil {
		// Error bisa karena: email/username sudah ada, role tidak valid, dll
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Step 5: Return data user baru (tanpa password dan token)
	response.HttpRresponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user.User, // Data user yang baru dibuat
		Gin:  ctx,
	})
}

// Update menangani HTTP PUT/PATCH request untuk update data user
// Menggunakan UUID dari URL parameter untuk identifikasi user
func (u *UserController) Update(ctx *gin.Context) {
	// Step 1: Inisialisasi struct untuk data update
	request := &dto.UpdateRequest{}
	// Step 2: Ambil UUID dari URL parameter (contoh: /users/:uuid)
	uuid := ctx.Param("uuid")

	// Step 3: Binding JSON request ke struct
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Step 4: Validasi data update
	// Validasi field yang akan diupdate
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMesage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpRresponse(response.ParamHttpResp{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMesage,
			Data:    errResponse,
			Err:     err,
			Gin:     ctx,
		})
		return
	}

	// Step 5: Panggil service untuk update user
	// Service akan cek apakah user exist, lalu update data
	user, err := u.service.GetUserService().Update(ctx, request, uuid)
	if err != nil {
		// Error bisa karena: user tidak ditemukan, email/username conflict, dll
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	// Step 6: Return data user yang sudah diupdate
	response.HttpRresponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user, // Data user terbaru setelah update
		Gin:  ctx,
	})
}

// GetUserLogin menangani HTTP GET request untuk mendapatkan data user yang sedang login
// Biasanya digunakan setelah user login untuk mendapatkan profile mereka
// Memerlukan JWT token di header Authorization
func (u *UserController) GetUserLogin(ctx *gin.Context) {
	// Panggil service untuk mendapatkan data user dari token JWT
	// Service akan extract user info dari JWT token di header
	user, err := u.service.GetUserService().GetUserLogin(ctx.Request.Context())
	if err != nil {
		// Error bisa karena: token invalid, user tidak ditemukan, token expired
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return // Tambahkan return untuk konsistensi
	}

	// Return data user yang sedang login
	response.HttpRresponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user, // Data profile user
		Gin:  ctx,
	})
}

// GetUserByUUID menangani HTTP GET request untuk mendapatkan user berdasarkan UUID
// Endpoint ini untuk admin atau untuk mendapatkan data user lain
// URL pattern: GET /users/:uuid
func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	// Ambil UUID dari URL parameter dan panggil service
	user, err := u.service.GetUserService().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		// Error bisa karena: UUID tidak valid, user tidak ditemukan
		response.HttpRresponse(response.ParamHttpResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return // Tambahkan return untuk konsistensi
	}

	fmt.Println(user, "user")

	// Return data user berdasarkan UUID
	response.HttpRresponse(response.ParamHttpResp{
		Code: http.StatusOK,
		Data: user, // Data user yang dicari
		Gin:  ctx,
	})
}
