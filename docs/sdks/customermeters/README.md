# CustomerMeters

## Overview

### Available Operations

* [List](#list) - List Customer Meters
* [Get](#get) - Get Customer Meter

## List

List customer meters.

**Scopes**: `customer_meters:read`

### Example Usage

<!-- UsageSnippet language="go" operationID="customer_meters:list" method="get" path="/v1/customer-meters/" -->
```go
package main

import(
	"context"
	"os"
	pagogo "github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/models/operations"
	"log"
)

func main() {
    ctx := context.Background()

    s := pagogo.New(
        pagogo.WithSecurity(os.Getenv("PAGO_ACCESS_TOKEN")),
    )

    res, err := s.CustomerMeters.List(ctx, operations.CustomerMetersListRequest{
        OrganizationID: pagogo.Pointer(operations.CreateCustomerMetersListQueryParamOrganizationIDFilterStr(
            "1dbfc517-0bbf-4301-9ba8-555ca42b9737",
        )),
    })
    if err != nil {
        log.Fatal(err)
    }
    if res.ListResourceCustomerMeter != nil {
        for {
            // handle items

            res, err = res.Next()

            if err != nil {
                // handle error
            }

            if res == nil {
                break
            }
        }
    }
}
```

### Parameters

| Parameter                                                                                    | Type                                                                                         | Required                                                                                     | Description                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `ctx`                                                                                        | [context.Context](https://pkg.go.dev/context#Context)                                        | :heavy_check_mark:                                                                           | The context to use for the request.                                                          |
| `request`                                                                                    | [operations.CustomerMetersListRequest](../../models/operations/customermeterslistrequest.md) | :heavy_check_mark:                                                                           | The request object to use for the request.                                                   |
| `opts`                                                                                       | [][operations.Option](../../models/operations/option.md)                                     | :heavy_minus_sign:                                                                           | The options for this request.                                                                |

### Response

**[*operations.CustomerMetersListResponse](../../models/operations/customermeterslistresponse.md), error**

### Errors

| Error Type                    | Status Code                   | Content Type                  |
| ----------------------------- | ----------------------------- | ----------------------------- |
| apierrors.HTTPValidationError | 422                           | application/json              |
| apierrors.APIError            | 4XX, 5XX                      | \*/\*                         |

## Get

Get a customer meter by ID.

**Scopes**: `customer_meters:read`

### Example Usage

<!-- UsageSnippet language="go" operationID="customer_meters:get" method="get" path="/v1/customer-meters/{id}" -->
```go
package main

import(
	"context"
	"os"
	pagogo "github.com/pago-sh/pago-go"
	"log"
	"github.com/pago-sh/pago-go/models/components"
)

func main() {
    ctx := context.Background()

    s := pagogo.New(
        pagogo.WithSecurity(os.Getenv("PAGO_ACCESS_TOKEN")),
    )

    res, err := s.CustomerMeters.Get(ctx, "<value>")
    if err != nil {
        log.Fatal(err)
    }
    if res.CustomerMeter != nil {
        switch res.CustomerMeter.Customer.Type {
            case components.CustomerUnionTypeIndividual:
                // res.CustomerMeter.Customer.CustomerIndividual is populated
            case components.CustomerUnionTypeTeam:
                // res.CustomerMeter.Customer.CustomerTeam is populated
        }

    }
}
```

### Parameters

| Parameter                                                | Type                                                     | Required                                                 | Description                                              |
| -------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| `ctx`                                                    | [context.Context](https://pkg.go.dev/context#Context)    | :heavy_check_mark:                                       | The context to use for the request.                      |
| `id`                                                     | `string`                                                 | :heavy_check_mark:                                       | The customer meter ID.                                   |
| `opts`                                                   | [][operations.Option](../../models/operations/option.md) | :heavy_minus_sign:                                       | The options for this request.                            |

### Response

**[*operations.CustomerMetersGetResponse](../../models/operations/customermetersgetresponse.md), error**

### Errors

| Error Type                    | Status Code                   | Content Type                  |
| ----------------------------- | ----------------------------- | ----------------------------- |
| apierrors.ResourceNotFound    | 404                           | application/json              |
| apierrors.HTTPValidationError | 422                           | application/json              |
| apierrors.APIError            | 4XX, 5XX                      | \*/\*                         |