package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	user struct {
		ID       int            `json:"id"`
		Name     string         `json:"name"`
		Services []userServices `json:"services"`
	}
	userServices struct {
		Name string `json:"name"`
	}
	inputUnsort struct {
		Unsorted []int `json:"unsorted"`
	}
	outputSort struct {
		Sorted []int `json:"sorted"`
	}
)

var (
	users = map[int]*user{}
)

func sortArray(c echo.Context) error {
	u := new(inputUnsort)
	if err := c.Bind(u); err != nil {
		return err
	}
	sort.Ints(u.Unsorted)
	s := &outputSort{
		Sorted: u.Unsorted,
	}
	return c.JSON(http.StatusOK, s)
}

func sortArray2(c echo.Context) error {
	u := new(inputUnsort)
	if err := c.Bind(u); err != nil {
		return err
	}

	max := u.Unsorted[0]
	for _, value := range u.Unsorted {
		if value > max {
			max = value
		}
	}

	for i := range u.Unsorted {
		for i2, val2 := range u.Unsorted {
			if (i + 1) == val2 {
				u.Unsorted = move(i2, i, u.Unsorted)
				break
			}
		}
	}

	s := &outputSort{
		Sorted: u.Unsorted,
	}
	return c.JSON(http.StatusOK, s)
}

func move(indexToRemove int, indexWhereToInsert int, slice []int) []int {

	val := slice[indexToRemove]

	slice = append(slice[:indexToRemove], slice[indexToRemove+1:]...)

	newSlice := make([]int, indexWhereToInsert+1)
	copy(newSlice, slice[:indexWhereToInsert])
	newSlice[indexWhereToInsert] = val

	slice = append(newSlice, slice[indexWhereToInsert:]...)
	return slice
}

func createUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.ID == 0 {
		return c.JSON(http.StatusBadRequest, "Missing id")
	}
	users[u.ID] = u
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if users[id] == nil {
		return c.JSON(http.StatusInternalServerError, "User not found")
	}
	return c.JSON(http.StatusOK, users[id])
}

func updateUser(c echo.Context) error {
	u := new(user)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if users[id] == nil {
		return c.JSON(http.StatusInternalServerError, "User not found")
	}
	users[id] = u
	return c.JSON(http.StatusOK, users[id])
}

func deleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if users[id] == nil {
		return c.JSON(http.StatusInternalServerError, "User not found")
	}
	delete(users, id)
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/users", createUser)
	e.POST("/sort", sortArray)
	e.POST("/sort2", sortArray2)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
