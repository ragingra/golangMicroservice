package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/model"
	"main/repository/order"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Order struct {
	Repo order.OrderRepository
}

func (o *Order) Create(c *gin.Context) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode request body:", err)
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	order := model.Order{
		OrderID:    rand.Uint64(),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		CreatedAt:  &now,
	}

	err := o.Repo.Insert(c, order)
	if err != nil {
		fmt.Println("failed to insert order:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal order:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write(res)
}

func (o *Order) List(c *gin.Context) {
	cursorStr := c.Query("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	fmt.Printf("customer cursor before: %s\n", cursorStr)

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("customer cursor: %d\n", cursor)

	const size = 50
	res, err := o.Repo.FindAll(c, order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all orders:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}

	response.Items = res.Orders
	response.Next = res.Cursor

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Writer.Write(data)
}

func (o *Order) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	foundOrder, err := o.Repo.FindByID(c, orderID)
	if errors.Is(err, order.ErrNotExist) {
		fmt.Println("failed to find order by ID:", err)
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find order by ID:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(c.Writer).Encode(foundOrder); err != nil {
		fmt.Println("failed to encode order:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) UpdateByID(c *gin.Context) {
	var body struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := c.Param("id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	foundOrder, err := o.Repo.FindByID(c, orderID)
	if errors.Is(err, order.ErrNotExist) {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case shippedStatus:
		if foundOrder.ShippedAt != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		foundOrder.ShippedAt = &now
	case completedStatus:
		if foundOrder.CompletedAt != nil || foundOrder.ShippedAt == nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		foundOrder.CompletedAt = &now
	default:
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.Update(c, foundOrder)
	if err != nil {
		fmt.Println("failed to update order:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(c.Writer).Encode(foundOrder); err != nil {
		fmt.Println("failed to encode order:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Order) DeleteByID(c *gin.Context) {
	idParam := c.Param("id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(c, orderID)
	if errors.Is(err, order.ErrNotExist) {
		fmt.Println("failed to find order by ID:", err)
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find order by ID:", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
