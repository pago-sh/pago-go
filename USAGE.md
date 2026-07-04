<!-- Start SDK Example Usage [usage] -->
```go
package main

import (
	"context"
	pagogo "github.com/pago-sh/pago-go"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	s := pagogo.New(
		pagogo.WithSecurity(os.Getenv("PAGO_ACCESS_TOKEN")),
	)

	res, err := s.Organizations.List(ctx, nil, pagogo.Pointer[int64](1), pagogo.Pointer[int64](10), nil)
	if err != nil {
		log.Fatal(err)
	}
	if res.ListResourceOrganization != nil {
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
<!-- End SDK Example Usage [usage] -->