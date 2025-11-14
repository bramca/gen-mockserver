package genmock

import (
	"fmt"
	"net"
	"slices"
	"testing"
	"time"

	"go.yaml.in/yaml/v4"

	"github.com/google/uuid"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/stretchr/testify/assert"
)

func Test_RandStringBytesRmndr_ReturnsRandomString(t *testing.T) {
	t.Parallel()

	// Arrange
	n := 10

	// Act
	result := RandStringBytesRmndr(n)

	// Assert
	assert.NotEmpty(t, result)
	assert.Len(t, result, n)
}

func Test_generateExampleData_ReturnsCorrectStringFormat(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		propSchema    *base.Schema
		formatChecker func(string) bool
	}{
		"enum": {
			propSchema: &base.Schema{
				Enum: []*yaml.Node{
					{
						Value: "test 1",
					},
					{
						Value: "test 2",
					},
				},
			},
			formatChecker: func(s string) bool {
				return slices.Contains([]string{"test 1", "test 2"}, s)
			},
		},
		"date-time": {
			propSchema: &base.Schema{
				Format: "date-time",
			},
			formatChecker: func(s string) bool {
				_, err := time.Parse("2006-01-02 15:04:05.000000000 +0000 UTC", s)

				return assert.NoError(t, err)
			},
		},
		"uuid": {
			propSchema: &base.Schema{
				Format: "uuid",
			},
			formatChecker: func(s string) bool {
				_, err := uuid.Parse(s)

				return assert.NoError(t, err)
			},
		},
		"ip": {
			propSchema: &base.Schema{
				Format: "ip",
			},
			formatChecker: func(s string) bool {
				ip := net.ParseIP(s)

				return assert.NotNil(t, ip)
			},
		},
		"ip-cidr-block": {
			propSchema: &base.Schema{
				Format: "ip-cidr-block",
			},
			formatChecker: func(s string) bool {
				_, _, err := net.ParseCIDR(s)

				return assert.NoError(t, err)
			},
		},
		"mac-address": {
			propSchema: &base.Schema{
				Format: "mac-address",
			},
			formatChecker: func(s string) bool {
				_, err := net.ParseMAC(s)

				return assert.NoError(t, err)
			},
		},
		"address-or-block-or-range": {
			propSchema: &base.Schema{
				Format: "address-or-block-or-range",
			},
			formatChecker: func(s string) bool {
				ip := net.ParseIP(s)

				return assert.NotNil(t, ip)
			},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Act
			result := generateExampleData(data.propSchema)

			// Assert
			assert.NotEmpty(t, result)
			assert.True(t, data.formatChecker(result))
		})
	}
}

func Test_SpecV3toRequestStructureMap_ReturnsResponseBody(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		maxRecursion int
		expectedMap  map[string]map[string][]RequestStructure
	}{
		"0 recursion levels": {
			maxRecursion: 0,
			expectedMap: map[string]map[string][]RequestStructure{
				"get": {
					"/addresses": []RequestStructure{
						{
							Path:          "/addresses",
							Method:        "get",
							Body:          "",
							DbEntry:       "addresses",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
					},
					"/cart": []RequestStructure{
						{
							Path:          "/cart",
							Method:        "get",
							Body:          "",
							DbEntry:       "cart",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
					},
					"/orders": []RequestStructure{
						{
							Path:          "/orders",
							Method:        "get",
							Body:          "",
							DbEntry:       "orders",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
					},
					"/orders/:orderId": []RequestStructure{
						{
							Path:         "/orders/:orderId",
							Method:       "get",
							Body:         "",
							DbEntry:      "orders",
							ResponseCode: "200",
							ResponseBody: map[string]any{
								"created_at":   "",
								"id":           "",
								"items":        []any{},
								"status":       "",
								"total_amount": nil,
							},
							RequestParams: []string{"orderId"},
							RequestBody:   nil,
						},
					},
					"/products": []RequestStructure{
						{
							Path:          "/products",
							Method:        "get",
							Body:          "",
							DbEntry:       "products",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
						{
							Path:          "/products?category=",
							Method:        "get",
							Body:          "",
							DbEntry:       "products",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
						{
							Path:          "/products?search=",
							Method:        "get",
							Body:          "",
							DbEntry:       "products",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
						{
							Path:          "/products?min_price=",
							Method:        "get",
							Body:          "",
							DbEntry:       "products",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
						{
							Path:          "/products?max_price=",
							Method:        "get",
							Body:          "",
							DbEntry:       "products",
							ResponseCode:  "200",
							ResponseBody:  []any{},
							RequestParams: []string{},
							RequestBody:   nil,
						},
					},
					"/products/:id": []RequestStructure{
						{
							Path:         "/products/:id",
							Method:       "get",
							Body:         "",
							DbEntry:      "products",
							ResponseCode: "200",
							ResponseBody: map[string]any{
								"category":    "",
								"created_at":  "",
								"description": "",
								"id":          "",
								"image_url":   "",
								"name":        "",
								"price":       nil,
								"stock":       0,
								"updated_at":  "",
							},
							RequestParams: []string{"id"},
							RequestBody:   nil,
						},
					},
				},
				"post": {
					"/addresses": []RequestStructure{
						{
							Path:          "/addresses",
							Method:        "post",
							Body:          "",
							DbEntry:       "addresses",
							ResponseCode:  "201",
							ResponseBody:  nil,
							RequestParams: []string{},
							RequestBody: map[string]any{
								"city":        "",
								"country":     "",
								"line1":       "",
								"line2":       "",
								"postal_code": "",
								"state":       "",
							},
						},
					},
					"/auth/login": []RequestStructure{
						{
							Path:          "/auth/login",
							Method:        "post",
							Body:          "",
							DbEntry:       "auth-login",
							ResponseCode:  "200",
							ResponseBody:  nil,
							RequestParams: []string{},
							RequestBody: map[string]any{
								"email":    "",
								"password": "",
							},
						},
					},
					"/auth/register": []RequestStructure{
						{
							Path:          "/auth/register",
							Method:        "post",
							Body:          "",
							DbEntry:       "auth-register",
							ResponseCode:  "201",
							ResponseBody:  nil,
							RequestParams: []string{},
							RequestBody: map[string]any{
								"email":    "",
								"name":     "",
								"password": "",
							},
						},
					},
					"/cart/items": []RequestStructure{
						{
							Path:          "/cart/items",
							Method:        "post",
							Body:          "",
							DbEntry:       "cart-items",
							ResponseCode:  "200",
							ResponseBody:  nil,
							RequestParams: []string{},
							RequestBody: map[string]any{
								"product_id": "", "quantity": 1,
							},
						},
					},
					"/checkout": []RequestStructure{
						{
							Path:         "/checkout",
							Method:       "post",
							Body:         "",
							DbEntry:      "checkout",
							ResponseCode: "201",
							ResponseBody: map[string]any{
								"created_at":   "",
								"id":           "",
								"items":        []any{},
								"status":       "",
								"total_amount": nil,
							},
							RequestParams: []string{},
							RequestBody: map[string]any{
								"address_id":        "",
								"payment_method_id": "",
							},
						},
					},
				},
			},
		},
		"1 recursion level": {
			maxRecursion: 1,

			expectedMap: map[string]map[string][]RequestStructure{
				"get": {
					"/addresses": {
						{
							Path: "/addresses", Method: "get", Body: "", DbEntry: "addresses", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"city": "", "country": "", "line1": "", "line2": "", "postal_code": "", "state": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						},
					}, "/cart": {
						{
							Path: "/cart", Method: "get", Body: "", DbEntry: "cart", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"product_id": "", "quantity": 1,
								},
							}, RequestParams: []string{}, RequestBody: nil,
						},
					}, "/orders": {
						{
							Path: "/orders", Method: "get", Body: "", DbEntry: "orders", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"created_at": "", "id": "", "items": []any{}, "status": "", "total_amount": nil,
								},
							}, RequestParams: []string{}, RequestBody: nil,
						},
					}, "/orders/:orderId": {
						{
							Path: "/orders/:orderId", Method: "get", Body: "", DbEntry: "orders", ResponseCode: "200", ResponseBody: map[string]any{
								"created_at": "", "id": "", "items": []map[string]any{
									{
										"product_id": "", "quantity": 1,
									},
								}, "status": "", "total_amount": nil,
							}, RequestParams: []string{
								"orderId",
							}, RequestBody: nil,
						},
					}, "/products": {
						{
							Path: "/products", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						}, {
							Path: "/products?category=", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						}, {
							Path: "/products?search=", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						}, {
							Path: "/products?min_price=", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						}, {
							Path: "/products?max_price=", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: []map[string]any{
								{
									"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
								},
							}, RequestParams: []string{}, RequestBody: nil,
						},
					}, "/products/:id": {
						{
							Path: "/products/:id", Method: "get", Body: "", DbEntry: "products", ResponseCode: "200", ResponseBody: map[string]any{
								"category": "", "created_at": "", "description": "", "id": "", "image_url": "", "name": "", "price": nil, "stock": 0, "updated_at": "",
							}, RequestParams: []string{
								"id",
							}, RequestBody: nil,
						},
					},
				}, "post": {
					"/addresses": {
						{
							Path: "/addresses", Method: "post", Body: "", DbEntry: "addresses", ResponseCode: "201", ResponseBody: nil, RequestParams: []string{}, RequestBody: map[string]any{
								"city": "", "country": "", "line1": "", "line2": "", "postal_code": "", "state": "",
							},
						},
					}, "/auth/login": {
						{
							Path: "/auth/login", Method: "post", Body: "", DbEntry: "auth-login", ResponseCode: "200", ResponseBody: nil, RequestParams: []string{}, RequestBody: map[string]any{
								"email": "", "password": "",
							},
						},
					}, "/auth/register": {
						{
							Path: "/auth/register", Method: "post", Body: "", DbEntry: "auth-register", ResponseCode: "201", ResponseBody: nil, RequestParams: []string{}, RequestBody: map[string]any{
								"email": "", "name": "", "password": "",
							},
						},
					}, "/cart/items": {
						{
							Path: "/cart/items", Method: "post", Body: "", DbEntry: "cart-items", ResponseCode: "200", ResponseBody: nil, RequestParams: []string{}, RequestBody: map[string]any{
								"product_id": "", "quantity": 1,
							},
						},
					}, "/checkout": {
						{
							Path: "/checkout", Method: "post", Body: "", DbEntry: "checkout", ResponseCode: "201", ResponseBody: map[string]any{
								"created_at": "", "id": "", "items": []map[string]any{
									{
										"product_id": "", "quantity": 1,
									},
								}, "status": "", "total_amount": nil,
							}, RequestParams: []string{}, RequestBody: map[string]any{
								"address_id": "", "payment_method_id": "",
							},
						},
					},
				},
			},
		},
	}
	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Act
			resultMap := SpecV3toRequestStructureMap("./testdata/examplev3.yaml", data.maxRecursion, false)

			// Assert
			assert.Equal(t, data.expectedMap, resultMap)
		})
	}
}

func Test_SpecV2toRequestStructureMap_ReturnsResponseBody(t *testing.T) {

	t.Parallel()
	// Assert
	expectedMap := map[string]map[string][]RequestStructure{
		"delete": {
			"/v1/products/:productId": {
				{
					Path:          "/v1/products/:productId",
					Method:        "delete",
					Body:          "",
					DbEntry:       "products",
					ResponseCode:  "204",
					ResponseBody:  nil,
					RequestParams: []string{"productId"},
					RequestBody:   nil,
				},
			},
		},

		"get": {
			"/v1/products": {
				{
					Path:         "/v1/products",
					Method:       "get",
					Body:         "",
					DbEntry:      "products",
					ResponseCode: "200",
					ResponseBody: map[string]any{
						"items": []map[string]any{
							{
								"description": "",
								"id":          "",
								"metadata":    nil,
								"name":        "",
								"price":       nil,
								"tags":        []any{},
							},
						},
						"page":       0,
						"pageSize":   0,
						"totalItems": 0,
					},
					RequestParams: []string{},
					RequestBody:   nil,
				},
				{
					Path:         "/v1/products?page=",
					Method:       "get",
					Body:         "",
					DbEntry:      "products",
					ResponseCode: "200",
					ResponseBody: map[string]any{
						"items": []map[string]any{
							{
								"description": "",
								"id":          "",
								"metadata":    nil,
								"name":        "",
								"price":       nil,
								"tags":        []any{},
							},
						},
						"page":       0,
						"pageSize":   0,
						"totalItems": 0,
					},
					RequestParams: []string{},
					RequestBody:   nil,
				},
				{
					Path:         "/v1/products?pageSize=",
					Method:       "get",
					Body:         "",
					DbEntry:      "products",
					ResponseCode: "200",
					ResponseBody: map[string]any{
						"items": []map[string]any{
							{
								"description": "",
								"id":          "",
								"metadata":    nil,
								"name":        "",
								"price":       nil,
								"tags":        []any{},
							},
						},
						"page":       0,
						"pageSize":   0,
						"totalItems": 0,
					},
					RequestParams: []string{},
					RequestBody:   nil,
				},
			},

			"/v1/products/:productId": {
				{
					Path:         "/v1/products/:productId",
					Method:       "get",
					Body:         "",
					DbEntry:      "products",
					ResponseCode: "200",
					ResponseBody: map[string]any{
						"description": "",
						"id":          "",
						"metadata":    nil,
						"name":        "",
						"price":       nil,
						"tags":        []any{},
					},
					RequestParams: []string{"productId"},
					RequestBody:   nil,
				},
			},

			"/v1/users/:userId/orders": {
				{
					Path:         "/v1/users/:userId/orders",
					Method:       "get",
					Body:         "",
					DbEntry:      "users-orders",
					ResponseCode: "200",
					ResponseBody: []map[string]any{
						{
							"id":         "",
							"items":      []map[string]any{},
							"status":     "",
							"totalPrice": nil,
							"userId":     "",
						},
					},
					RequestParams: []string{"userId"},
					RequestBody:   nil,
				},
				{
					Path:         "/v1/users/:userId/orders?status=",
					Method:       "get",
					Body:         "",
					DbEntry:      "users-orders",
					ResponseCode: "200",
					ResponseBody: []map[string]any{
						{
							"id":         "",
							"items":      []map[string]any{},
							"status":     "",
							"totalPrice": nil,
							"userId":     "",
						},
					},
					RequestParams: []string{"userId"},
					RequestBody:   nil,
				},
			},
		},

		"post": {
			"/v1/products": {
				{
					Path:         "/v1/products",
					Method:       "post",
					Body:         "",
					DbEntry:      "products",
					ResponseCode: "201",
					ResponseBody: map[string]any{
						"description": "",
						"id":          "",
						"metadata":    nil,
						"name":        "",
						"price":       nil,
						"tags":        []any{},
					},
					RequestParams: []string{},
					RequestBody:   map[string]any{},
				},
			},
		},
	}

	// Act
	resultMap := SpecV2toRequestStructureMap("./testdata/examplev2v2.yaml", 1, false)

	fmt.Printf("resultMap: %+v\n", resultMap)

	// Assert
	assert.Equal(t, expectedMap, resultMap)
}
