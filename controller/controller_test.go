package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"deltaFiJWT/dao"
	"deltaFiJWT/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func setupTest(t *testing.T) *gin.Engine {
	assert.NoError(t, dao.InitTestDB())
	assert.NoError(t, dao.DB.Migrator().DropTable(&dao.User{}))
	assert.NoError(t, dao.DB.AutoMigrate(&dao.User{}))

	r := gin.Default()
	SetUpRouter(r)

	t.Cleanup(func() {
		assert.NoError(t, dao.DB.Migrator().DropTable(&dao.User{}))
	})
	return r
}

func createUser(input types.CreateUserInput, t *testing.T, router *gin.Engine) *httptest.ResponseRecorder {
	jsonstr, err := json.Marshal(input)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	router.ServeHTTP(w, req)
	return w
}

func login(input Login, t *testing.T, router *gin.Engine) *httptest.ResponseRecorder {
	jsonstr, err := json.Marshal(input)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	router.ServeHTTP(w, req)
	return w
}

func greeting(cookie []*http.Cookie, router *gin.Engine) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	for _, c := range cookie {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)
	return w
}

func TestCreateUser(t *testing.T) {
	router := setupTest(t)

	// Test create user
	input := types.CreateUserInput{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}

	w := createUser(input, t, router)
	assert.Equal(t, http.StatusOK, w.Code)

	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, gjson.Get(string(body), "code").Int())
	assert.Equal(t, "jackson@example.com", gjson.Get(string(body), "user.email").String())
	assert.Equal(t, "jackson", gjson.Get(string(body), "user.firstName").String())
	assert.Equal(t, "wang", gjson.Get(string(body), "user.lastName").String())

	// Test creat user with same email
	w = createUser(input, t, router)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeCreateUserFail, gjson.Get(string(body), "code").Int())

	// Test invalid input
	input = types.CreateUserInput{
		Email:     "",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	w = createUser(input, t, router)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeInvalidInput, gjson.Get(string(body), "code").Int())
}

func TestLogin(t *testing.T) {
	router := setupTest(t)

	// Test login success
	input := types.CreateUserInput{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	w := createUser(input, t, router)
	assert.Equal(t, http.StatusOK, w.Code)

	loginInput := Login{
		Email:    "jackson@example.com",
		Password: "123",
	}
	w = login(loginInput, t, router)
	assert.Equal(t, http.StatusOK, w.Code)
	cookie := w.Result().Cookies()

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	for _, c := range cookie {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeSuccess, gjson.Get(string(body), "code").Int())
	assert.EqualValues(t, "Hello jackson wang", gjson.Get(string(body), "greeting").String())

	// Test login fail
	loginInput = Login{
		Email:    "jackson@example.com",
		Password: "123312",
	}
	w = login(loginInput, t, router)
	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.EqualValues(t, CodeLoginFail, gjson.Get(string(body), "code").Int())
}

func TestUpdateUser(t *testing.T) {
	router := setupTest(t)

	// Test update user success
	input := types.CreateUserInput{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	w := createUser(input, t, router)
	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	userId := gjson.Get(string(body), "user.id").Int()
	assert.Equal(t, http.StatusOK, w.Code)

	loginInput := Login{
		Email:    "jackson@example.com",
		Password: "123",
	}
	w = login(loginInput, t, router)
	assert.Equal(t, http.StatusOK, w.Code)
	cookie := w.Result().Cookies()

	updateInput := types.UpdateUserInput{
		ID:        uint(userId),
		FirstName: "amy",
		LastName:  "june",
	}
	jsonstr, err := json.Marshal(updateInput)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for _, c := range cookie {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeSuccess, gjson.Get(string(body), "code").Int())

	w = greeting(cookie, router)
	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeSuccess, gjson.Get(string(body), "code").Int())
	assert.EqualValues(t, "Hello amy june", gjson.Get(string(body), "greeting").String())
}

func TestDeleteUser(t *testing.T) {
	router := setupTest(t)

	// Test delete user success
	input := types.CreateUserInput{
		Email:     "jackson@example.com",
		FirstName: "jackson",
		LastName:  "wang",
		Password:  "123",
	}
	w := createUser(input, t, router)
	body, err := io.ReadAll(w.Body)
	assert.NoError(t, err)
	userId := gjson.Get(string(body), "user.id").Int()
	assert.Equal(t, http.StatusOK, w.Code)

	loginInput := Login{
		Email:    "jackson@example.com",
		Password: "123",
	}
	w = login(loginInput, t, router)
	assert.Equal(t, http.StatusOK, w.Code)
	cookie := w.Result().Cookies()

	updateInput := types.DeleteUserInput{
		ID: uint(userId),
	}
	jsonstr, err := json.Marshal(updateInput)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/user", bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for _, c := range cookie {
		req.AddCookie(c)
	}
	router.ServeHTTP(w, req)

	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeSuccess, gjson.Get(string(body), "code").Int())

	w = greeting(cookie, router)
	body, err = io.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.EqualValues(t, CodeUserNotFound, gjson.Get(string(body), "code").Int())
}
