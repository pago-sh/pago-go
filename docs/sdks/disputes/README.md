# Disputes

## Overview

### Available Operations

* [List](#list) - List Disputes
* [Get](#get) - Get Dispute

## List

List disputes.

**Scopes**: `disputes:read`

### Example Usage

<!-- UsageSnippet language="go" operationID="disputes:list" method="get" path="/v1/disputes/" -->
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

    res, err := s.Disputes.List(ctx, operations.DisputesListRequest{
        OrganizationID: pagogo.Pointer(operations.CreateDisputesListQueryParamOrganizationIDFilterStr(
            "1dbfc517-0bbf-4301-9ba8-555ca42b9737",
        )),
    })
    if err != nil {
        log.Fatal(err)
    }
    if res.ListResourceDispute != nil {
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

| Parameter                                                                        | Type                                                                             | Required                                                                         | Description                                                                      |
| -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| `ctx`                                                                            | [context.Context](https://pkg.go.dev/context#Context)                            | :heavy_check_mark:                                                               | The context to use for the request.                                              |
| `request`                                                                        | [operations.DisputesListRequest](../../models/operations/disputeslistrequest.md) | :heavy_check_mark:                                                               | The request object to use for the request.                                       |
| `opts`                                                                           | [][operations.Option](../../models/operations/option.md)                         | :heavy_minus_sign:                                                               | The options for this request.                                                    |

### Response

**[*operations.DisputesListResponse](../../models/operations/disputeslistresponse.md), error**

### Errors

| Error Type                    | Status Code                   | Content Type                  |
| ----------------------------- | ----------------------------- | ----------------------------- |
| apierrors.HTTPValidationError | 422                           | application/json              |
| apierrors.APIError            | 4XX, 5XX                      | \*/\*                         |

## Get

Get a dispute by ID.

**Scopes**: `disputes:read`

### Example Usage

<!-- UsageSnippet language="go" operationID="disputes:get" method="get" path="/v1/disputes/{id}" -->
```go
package main

import(
	"context"
	"os"
	pagogo "github.com/pago-sh/pago-go"
	"log"
)

func main() {
    ctx := context.Background()

    s := pagogo.New(
        pagogo.WithSecurity(os.Getenv("PAGO_ACCESS_TOKEN")),
    )

    res, err := s.Disputes.Get(ctx, "<value>")
    if err != nil {
        log.Fatal(err)
    }
    if res.Dispute != nil {
        // handle response
    }
}
```

### Parameters

| Parameter                                                | Type                                                     | Required                                                 | Description                                              |
| -------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- | -------------------------------------------------------- |
| `ctx`                                                    | [context.Context](https://pkg.go.dev/context#Context)    | :heavy_check_mark:                                       | The context to use for the request.                      |
| `id`                                                     | `string`                                                 | :heavy_check_mark:                                       | The dispute ID.                                          |
| `opts`                                                   | [][operations.Option](../../models/operations/option.md) | :heavy_minus_sign:                                       | The options for this request.                            |

### Response

**[*operations.DisputesGetResponse](../../models/operations/disputesgetresponse.md), error**

### Errors

| Error Type                    | Status Code                   | Content Type                  |
| ----------------------------- | ----------------------------- | ----------------------------- |
| apierrors.ResourceNotFound    | 404                           | application/json              |
| apierrors.HTTPValidationError | 422                           | application/json              |
| apierrors.APIError            | 4XX, 5XX                      | \*/\*                         |