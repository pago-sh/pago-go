# Pago Go SDK

Go client for the [Pago](https://pago.sh) API.

> This SDK is generated from the Pago OpenAPI specification. Do not edit it by hand.

## Installation

```bash
go get github.com/pago-sh/pago-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/v2026_04"
)

func main() {
	client := v2026_04.New(pago.WithAccessToken(os.Getenv("PAGO_ACCESS_TOKEN")))

	product, err := client.Products.Get(context.Background(), "PRODUCT_ID")
	if err != nil {
		panic(err)
	}
	fmt.Println(product.Name)
}
```

Each API version lives in its own package:

- `github.com/pago-sh/pago-go/v2026_04` — API version `2026-04`

### Configuration

`New` accepts the options of the runtime package:

```go
client := v2026_04.New(
	pago.WithAccessToken("..."),
	pago.WithBaseURL("https://api.pago.sh"),
	pago.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
)
```

### Pagination

Paginated endpoints expose an `...AutoPaging` iterator that fetches every page
transparently:

```go
for order, err := range client.Orders.ListAutoPaging(ctx, v2026_04.OrdersListParams{}) {
	if err != nil {
		return err
	}
	fmt.Println(order.ID)
}
```

### Errors

Every error returned by the SDK implements `pago.Error`. Documented error
responses have a generated type carrying the decoded body:

```go
product, err := client.Products.Get(ctx, "unknown")

var notFound *v2026_04.ResourceNotFoundError
switch {
case errors.As(err, &notFound):
	// notFound.Data holds the decoded body.
case errors.As(err, new(*pago.RateLimitError)):
	// Retry later.
}
```

### Unions

The API uses polymorphic schemas that Go cannot express directly. They are
generated as a struct with one pointer field per variant, with the JSON
marshalling wired up for you:

```go
switch {
case benefit.BenefitCustom != nil:
	fmt.Println(benefit.BenefitCustom.Description)
case benefit.BenefitDiscord != nil:
	fmt.Println(benefit.BenefitDiscord.Description)
}
```

### Webhooks

`ValidateEvent` verifies the [Standard Webhooks](https://www.standardwebhooks.com)
signature of a request and returns the typed event:

```go
func handler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	event, err := v2026_04.ValidateEvent(body, r.Header, os.Getenv("PAGO_WEBHOOK_SECRET"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	switch event := event.(type) {
	case *v2026_04.WebhookOrderCreatedPayload:
		fmt.Println(event.Data.ID)
	}
}
```

## License

MIT
