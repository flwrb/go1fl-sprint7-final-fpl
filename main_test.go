package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, min(len(cafeList["moscow"]), 100)},
	}
	for _, v := range requests {
		vStr := "/cafe?count=" + strconv.Itoa(v.count) + "&city=moscow"
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", vStr, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		resBody := response.Body.String()
		if resBody == "" {
			assert.Equal(t, 0, v.want)
		} else {
			resBodySplit := strings.Split(resBody, ",")
			assert.Equal(t, v.want, len(resBodySplit))
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		vReq := "/cafe?city=moscow&search=" + v.search
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", vReq, nil)
		handler.ServeHTTP(response, req)
		require.Equal(t, http.StatusOK, response.Code)
		resBody := response.Body.String()

		names := []string{}
		if resBody != "" {
			names = strings.Split(resBody, ",")
		}
		assert.Equal(t, v.wantCount, len(names))

		for _, name := range names {
			assert.Contains(t, strings.ToLower(name), strings.ToLower(v.search))
		}
	}
}
