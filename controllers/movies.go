package controllers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Succsess bool   `json:"success"`
	Message  string `json:"message"`
	Results  any    `json:"results,omitempty"`
}

type Movie struct {
	Id          int    `json:"id"`
	Title       string `json:"title" form:"title"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
}

type ListMovies []Movie

var Data = ListMovies{
	{
		Id:          1,
		Title:       "Lovely Runner",
		Image:       "example/lovelyrunner.com",
		Description: "Lorem ipsum dolor sit amet",
	},
	{
		Id:          2,
		Title:       "Love Next Door",
		Image:       "example/lovenextdoor.com",
		Description: "Lorem ipsum dolor sit amet",
	},
	{
		Id:          3,
		Title:       "Uncontrollably Fond",
		Image:       "example/uncontrollablyfond.com",
		Description: "Lorem ipsum dolor sit amet",
	},
	{
		Id:          4,
		Title:       "Hometown Cha-cha-cha",
		Image:       "example/hometown.com",
		Description: "Lorem ipsum dolor sit amet",
	},
	{
		Id:          5,
		Title:       "Start Up",
		Image:       "example/startup.com",
		Description: "Lorem ipsum dolor sit amet",
	},
}

var sequence int = len(Data)

func GetAllMovies(ctx *gin.Context) {
	search := ctx.Query("search")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "2"))
	data := Data

	sort.Slice(data, func(i, j int) bool {
		return data[i].Title < data[j].Title
	})

	if search == "" {
		if page*limit > len(data) {
			ctx.JSON(http.StatusOK, Response{
				Succsess: true,
				Message:  "List all movies",
				Results:  data[(page-1)*limit : len(Data)],
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Succsess: true,
			Message:  "List all movies",
			Results:  data[(page-1)*limit : page*limit],
		})
	} else {
		var resMov Movie
		var listDetails []Movie
		for _, dataSearch := range data {
			if strings.Contains(strings.ToLower(dataSearch.Title), search) || strings.Contains(dataSearch.Title, search) {
				resMov = dataSearch
				listDetails = append(listDetails, resMov)
			}
		}
		if resMov == (Movie{}) {
			ctx.JSON(http.StatusNotFound, Response{
				Succsess: false,
				Message:  "Movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, Response{
			Succsess: true,
			Message:  "Details movie",
			Results:  listDetails,
		})
	}
}

func GetMovieById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	data := Data
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Succsess: false,
			Message:  "Wrong movie id format",
		})
		return
	}

	var resMov Movie

	for _, data := range data {
		if data.Id == id {
			resMov = data
		}
	}
	if resMov == (Movie{}) {
		ctx.JSON(http.StatusNotFound, Response{
			Succsess: false,
			Message:  "Movie not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, Response{
		Succsess: true,
		Message:  "Details movie",
		Results:  resMov,
	})
}

func AddMovie(ctx *gin.Context) {
	var form Movie

	ctx.ShouldBind(&form)
	sequence++
	form.Id = sequence
	Data = append(Data, form)

	ctx.JSON(http.StatusOK, Response{
		Succsess: true,
		Message:  "Movie added",
		Results:  form,
	})
}

func EditMovie(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var form Movie
	ctx.ShouldBind(&form)
	for i, data := range Data {
		if data.Id == id {
			if form.Title != "" {
				Data[i].Title = form.Title
			}
			if form.Image != "" {
				Data[i].Image = form.Image
			}
			if form.Description != "" {
				Data[i].Description = form.Description
			}
			ctx.JSON(http.StatusOK, Response{
				Succsess: true,
				Message:  "movie detail has modify",
				Results:  Data[i],
			})
			return
		}
	}
}

func DeleteMovie(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	for i, data := range Data {
		if data.Id == id {
			Data = append(Data[:i], Data[i+1:]...)
			ctx.JSON(http.StatusOK, Response{
				Succsess: true,
				Message:  "movie deleted",
				Results:  data,
			})
			return
		}
	}
}
