# github.com/pago-sh/pago-go

Developer-friendly & type-safe Go SDK specifically catered to leverage *github.com/pago-sh/pago-go* API.

<div align="left">
    <a href="https://www.speakeasy.com/?utm_source=github-com/pago-sh/pago-go&utm_campaign=go"><img src="https://custom-icon-badges.demolab.com/badge/-Built%20By%20Speakeasy-212015?style=for-the-badge&logoColor=FBE331&logo=speakeasy&labelColor=545454" /></a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/License-MIT-blue.svg" style="width: 100px; height: 28px;" />
    </a>
</div>

<!-- Start Summary [summary] -->
## Summary

Pago API: Pago HTTP and Webhooks API

Read the docs at https://pago.sh/docs/api-reference
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [github.com/pago-sh/pago-go](#githubcompago-shpago-go)
  * [SDK Installation](#sdk-installation)
  * [SDK Example Usage](#sdk-example-usage)
  * [Authentication](#authentication)
  * [Available Resources and Operations](#available-resources-and-operations)
  * [Pagination](#pagination)
  * [Retries](#retries)
  * [Error Handling](#error-handling)
  * [Server Selection](#server-selection)
  * [Custom HTTP Client](#custom-http-client)
  * [Special Types](#special-types)
* [Development](#development)
  * [Maturity](#maturity)
  * [Contributions](#contributions)

<!-- End Table of Contents [toc] -->

<!-- Start SDK Installation [installation] -->
## SDK Installation

To add the SDK as a dependency to your project:
```bash
go get github.com/pago-sh/pago-go
```
<!-- End SDK Installation [installation] -->

<!-- Start SDK Example Usage [usage] -->
## SDK Example Usage

### Example

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

<!-- Start Authentication [security] -->
## Authentication

### Per-Client Security Schemes

This SDK supports the following security scheme globally:

| Name          | Type | Scheme      | Environment Variable |
| ------------- | ---- | ----------- | -------------------- |
| `AccessToken` | http | HTTP Bearer | `PAGO_ACCESS_TOKEN` |

You can configure it using the `WithSecurity` option when initializing the SDK client instance. For example:
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

### Per-Operation Security Schemes

Some operations in this SDK require the security scheme to be specified at the request level. For example:
```go
package main

import (
	"context"
	pagogo "github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/models/operations"
	"log"
)

func main() {
	ctx := context.Background()

	s := pagogo.New()

	res, err := s.CustomerPortal.BenefitGrants.List(ctx, operations.CustomerPortalBenefitGrantsListRequest{}, operations.CustomerPortalBenefitGrantsListSecurity{})
	if err != nil {
		log.Fatal(err)
	}
	if res.ListResourceCustomerBenefitGrant != nil {
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
<!-- End Authentication [security] -->

<!-- Start Available Resources and Operations [operations] -->
## Available Resources and Operations

<details open>
<summary>Available methods</summary>

### [BenefitGrants](docs/sdks/benefitgrants/README.md)

* [List](docs/sdks/benefitgrants/README.md#list) - List Benefit Grants

### [Benefits](docs/sdks/benefits/README.md)

* [List](docs/sdks/benefits/README.md#list) - List Benefits
* [Create](docs/sdks/benefits/README.md#create) - Create Benefit
* [Get](docs/sdks/benefits/README.md#get) - Get Benefit
* [Delete](docs/sdks/benefits/README.md#delete) - Delete Benefit
* [Update](docs/sdks/benefits/README.md#update) - Update Benefit
* [Grants](docs/sdks/benefits/README.md#grants) - List Benefit Grants

### [CheckoutLinks](docs/sdks/checkoutlinks/README.md)

* [List](docs/sdks/checkoutlinks/README.md#list) - List Checkout Links
* [Create](docs/sdks/checkoutlinks/README.md#create) - Create Checkout Link
* [Get](docs/sdks/checkoutlinks/README.md#get) - Get Checkout Link
* [Delete](docs/sdks/checkoutlinks/README.md#delete) - Delete Checkout Link
* [Update](docs/sdks/checkoutlinks/README.md#update) - Update Checkout Link

### [Checkouts](docs/sdks/checkouts/README.md)

* [List](docs/sdks/checkouts/README.md#list) - List Checkout Sessions
* [Create](docs/sdks/checkouts/README.md#create) - Create Checkout Session
* [Get](docs/sdks/checkouts/README.md#get) - Get Checkout Session
* [Update](docs/sdks/checkouts/README.md#update) - Update Checkout Session
* [ClientGet](docs/sdks/checkouts/README.md#clientget) - Get Checkout Session from Client
* [ClientUpdate](docs/sdks/checkouts/README.md#clientupdate) - Update Checkout Session from Client
* [ClientConfirm](docs/sdks/checkouts/README.md#clientconfirm) - Confirm Checkout Session from Client

### [CustomFields](docs/sdks/customfields/README.md)

* [List](docs/sdks/customfields/README.md#list) - List Custom Fields
* [Create](docs/sdks/customfields/README.md#create) - Create Custom Field
* [Get](docs/sdks/customfields/README.md#get) - Get Custom Field
* [Delete](docs/sdks/customfields/README.md#delete) - Delete Custom Field
* [Update](docs/sdks/customfields/README.md#update) - Update Custom Field

### [CustomerMeters](docs/sdks/customermeters/README.md)

* [List](docs/sdks/customermeters/README.md#list) - List Customer Meters
* [Get](docs/sdks/customermeters/README.md#get) - Get Customer Meter

### [CustomerPortal.BenefitGrants](docs/sdks/pagobenefitgrants/README.md)

* [List](docs/sdks/pagobenefitgrants/README.md#list) - List Benefit Grants
* [Get](docs/sdks/pagobenefitgrants/README.md#get) - Get Benefit Grant
* [Update](docs/sdks/pagobenefitgrants/README.md#update) - Update Benefit Grant

### [CustomerPortal.CustomerMeters](docs/sdks/pagocustomermeters/README.md)

* [List](docs/sdks/pagocustomermeters/README.md#list) - List Meters
* [Get](docs/sdks/pagocustomermeters/README.md#get) - Get Customer Meter

### [CustomerPortal.CustomerSession](docs/sdks/customersession/README.md)

* [Introspect](docs/sdks/customersession/README.md#introspect) - Introspect Customer Session
* [GetAuthenticatedUser](docs/sdks/customersession/README.md#getauthenticateduser) - Get Authenticated Portal User

### [CustomerPortal.Customers](docs/sdks/pagocustomers/README.md)

* [Get](docs/sdks/pagocustomers/README.md#get) - Get Customer
* [Update](docs/sdks/pagocustomers/README.md#update) - Update Customer
* [ListPaymentMethods](docs/sdks/pagocustomers/README.md#listpaymentmethods) - List Customer Payment Methods
* [AddPaymentMethod](docs/sdks/pagocustomers/README.md#addpaymentmethod) - Add Customer Payment Method
* [ConfirmPaymentMethod](docs/sdks/pagocustomers/README.md#confirmpaymentmethod) - Confirm Customer Payment Method
* [DeletePaymentMethod](docs/sdks/pagocustomers/README.md#deletepaymentmethod) - Delete Customer Payment Method

### [CustomerPortal.Downloadables](docs/sdks/downloadables/README.md)

* [List](docs/sdks/downloadables/README.md#list) - List Downloadables

### [CustomerPortal.LicenseKeys](docs/sdks/pagolicensekeys/README.md)

* [List](docs/sdks/pagolicensekeys/README.md#list) - List License Keys
* [Get](docs/sdks/pagolicensekeys/README.md#get) - Get License Key
* [Validate](docs/sdks/pagolicensekeys/README.md#validate) - Validate License Key
* [Activate](docs/sdks/pagolicensekeys/README.md#activate) - Activate License Key
* [Deactivate](docs/sdks/pagolicensekeys/README.md#deactivate) - Deactivate License Key

### [CustomerPortal.Members](docs/sdks/pagomembers/README.md)

* [ListMembers](docs/sdks/pagomembers/README.md#listmembers) - List Members
* [AddMember](docs/sdks/pagomembers/README.md#addmember) - Add Member
* [RemoveMember](docs/sdks/pagomembers/README.md#removemember) - Remove Member
* [UpdateMember](docs/sdks/pagomembers/README.md#updatemember) - Update Member

### [CustomerPortal.Orders](docs/sdks/pagoorders/README.md)

* [List](docs/sdks/pagoorders/README.md#list) - List Orders
* [Get](docs/sdks/pagoorders/README.md#get) - Get Order
* [Update](docs/sdks/pagoorders/README.md#update) - Update Order
* [Invoice](docs/sdks/pagoorders/README.md#invoice) - Get Order Invoice
* [GenerateInvoice](docs/sdks/pagoorders/README.md#generateinvoice) - Generate Order Invoice
* [GetPaymentStatus](docs/sdks/pagoorders/README.md#getpaymentstatus) - Get Order Payment Status
* [ConfirmRetryPayment](docs/sdks/pagoorders/README.md#confirmretrypayment) - Confirm Retry Payment

### [CustomerPortal.Organizations](docs/sdks/pagoorganizations/README.md)

* [Get](docs/sdks/pagoorganizations/README.md#get) - Get Organization

### [CustomerPortal.Seats](docs/sdks/seats/README.md)

* [ListSeats](docs/sdks/seats/README.md#listseats) - List Seats
* [AssignSeat](docs/sdks/seats/README.md#assignseat) - Assign Seat
* [RevokeSeat](docs/sdks/seats/README.md#revokeseat) - Revoke Seat
* [ResendInvitation](docs/sdks/seats/README.md#resendinvitation) - Resend Invitation
* [ListClaimedSubscriptions](docs/sdks/seats/README.md#listclaimedsubscriptions) - List Claimed Subscriptions

### [CustomerPortal.Subscriptions](docs/sdks/pagosubscriptions/README.md)

* [List](docs/sdks/pagosubscriptions/README.md#list) - List Subscriptions
* [Get](docs/sdks/pagosubscriptions/README.md#get) - Get Subscription
* [Cancel](docs/sdks/pagosubscriptions/README.md#cancel) - Cancel Subscription
* [Update](docs/sdks/pagosubscriptions/README.md#update) - Update Subscription

### [CustomerPortal.Wallets](docs/sdks/wallets/README.md)

* [List](docs/sdks/wallets/README.md#list) - List Wallets
* [Get](docs/sdks/wallets/README.md#get) - Get Wallet

### [CustomerSeats](docs/sdks/customerseats/README.md)

* [ListSeats](docs/sdks/customerseats/README.md#listseats) - List Seats
* [AssignSeat](docs/sdks/customerseats/README.md#assignseat) - Assign Seat
* [RevokeSeat](docs/sdks/customerseats/README.md#revokeseat) - Revoke Seat
* [ResendInvitation](docs/sdks/customerseats/README.md#resendinvitation) - Resend Invitation
* [GetClaimInfo](docs/sdks/customerseats/README.md#getclaiminfo) - Get Claim Info
* [ClaimSeat](docs/sdks/customerseats/README.md#claimseat) - Claim Seat

### [CustomerSessions](docs/sdks/customersessions/README.md)

* [Create](docs/sdks/customersessions/README.md#create) - Create Customer Session

### [Customers](docs/sdks/customers/README.md)

* [List](docs/sdks/customers/README.md#list) - List Customers
* [Create](docs/sdks/customers/README.md#create) - Create Customer
* [Export](docs/sdks/customers/README.md#export) - Export Customers
* [Get](docs/sdks/customers/README.md#get) - Get Customer
* [Delete](docs/sdks/customers/README.md#delete) - Delete Customer
* [Update](docs/sdks/customers/README.md#update) - Update Customer
* [GetExternal](docs/sdks/customers/README.md#getexternal) - Get Customer by External ID
* [DeleteExternal](docs/sdks/customers/README.md#deleteexternal) - Delete Customer by External ID
* [UpdateExternal](docs/sdks/customers/README.md#updateexternal) - Update Customer by External ID
* [GetState](docs/sdks/customers/README.md#getstate) - Get Customer State
* [GetStateExternal](docs/sdks/customers/README.md#getstateexternal) - Get Customer State by External ID

### [Discounts](docs/sdks/discounts/README.md)

* [List](docs/sdks/discounts/README.md#list) - List Discounts
* [Create](docs/sdks/discounts/README.md#create) - Create Discount
* [Get](docs/sdks/discounts/README.md#get) - Get Discount
* [Delete](docs/sdks/discounts/README.md#delete) - Delete Discount
* [Update](docs/sdks/discounts/README.md#update) - Update Discount

### [Disputes](docs/sdks/disputes/README.md)

* [List](docs/sdks/disputes/README.md#list) - List Disputes
* [Get](docs/sdks/disputes/README.md#get) - Get Dispute

### [EventTypes](docs/sdks/eventtypes/README.md)

* [List](docs/sdks/eventtypes/README.md#list) - List Event Types
* [Update](docs/sdks/eventtypes/README.md#update) - Update Event Type

### [Events](docs/sdks/events/README.md)

* [List](docs/sdks/events/README.md#list) - List Events
* [ListNames](docs/sdks/events/README.md#listnames) - List Event Names
* [Get](docs/sdks/events/README.md#get) - Get Event
* [Ingest](docs/sdks/events/README.md#ingest) - Ingest Events

### [Files](docs/sdks/files/README.md)

* [List](docs/sdks/files/README.md#list) - List Files
* [Create](docs/sdks/files/README.md#create) - Create File
* [Uploaded](docs/sdks/files/README.md#uploaded) - Complete File Upload
* [Delete](docs/sdks/files/README.md#delete) - Delete File
* [Update](docs/sdks/files/README.md#update) - Update File

### [LicenseKeys](docs/sdks/licensekeys/README.md)

* [List](docs/sdks/licensekeys/README.md#list) - List License Keys
* [Get](docs/sdks/licensekeys/README.md#get) - Get License Key
* [Update](docs/sdks/licensekeys/README.md#update) - Update License Key
* [GetActivation](docs/sdks/licensekeys/README.md#getactivation) - Get Activation
* [Validate](docs/sdks/licensekeys/README.md#validate) - Validate License Key
* [Activate](docs/sdks/licensekeys/README.md#activate) - Activate License Key
* [Deactivate](docs/sdks/licensekeys/README.md#deactivate) - Deactivate License Key

### [Members](docs/sdks/members/README.md)

* [ListMembers](docs/sdks/members/README.md#listmembers) - List Members
* [CreateMember](docs/sdks/members/README.md#createmember) - Create Member
* [GetMember](docs/sdks/members/README.md#getmember) - Get Member
* [DeleteMember](docs/sdks/members/README.md#deletemember) - Delete Member
* [UpdateMember](docs/sdks/members/README.md#updatemember) - Update Member

### [Meters](docs/sdks/meters/README.md)

* [List](docs/sdks/meters/README.md#list) - List Meters
* [Create](docs/sdks/meters/README.md#create) - Create Meter
* [Get](docs/sdks/meters/README.md#get) - Get Meter
* [Update](docs/sdks/meters/README.md#update) - Update Meter
* [Quantities](docs/sdks/meters/README.md#quantities) - Get Meter Quantities

### [Metrics](docs/sdks/metrics/README.md)

* [Get](docs/sdks/metrics/README.md#get) - Get Metrics
* [Limits](docs/sdks/metrics/README.md#limits) - Get Metrics Limits
* [ListDashboards](docs/sdks/metrics/README.md#listdashboards) - List Metric Dashboards
* [CreateDashboard](docs/sdks/metrics/README.md#createdashboard) - Create Metric Dashboard
* [GetDashboard](docs/sdks/metrics/README.md#getdashboard) - Get Metric Dashboard
* [DeleteDashboard](docs/sdks/metrics/README.md#deletedashboard) - Delete Metric Dashboard
* [UpdateDashboard](docs/sdks/metrics/README.md#updatedashboard) - Update Metric Dashboard

### [Oauth2](docs/sdks/oauth2/README.md)

* [Authorize](docs/sdks/oauth2/README.md#authorize) - Authorize
* [Token](docs/sdks/oauth2/README.md#token) - Request Token
* [Revoke](docs/sdks/oauth2/README.md#revoke) - Revoke Token
* [Introspect](docs/sdks/oauth2/README.md#introspect) - Introspect Token
* [Userinfo](docs/sdks/oauth2/README.md#userinfo) - Get User Info

#### [Oauth2.Clients](docs/sdks/clients/README.md)

* [Create](docs/sdks/clients/README.md#create) - Create Client
* [Get](docs/sdks/clients/README.md#get) - Get Client
* [Update](docs/sdks/clients/README.md#update) - Update Client
* [Delete](docs/sdks/clients/README.md#delete) - Delete Client

### [Orders](docs/sdks/orders/README.md)

* [List](docs/sdks/orders/README.md#list) - List Orders
* [Export](docs/sdks/orders/README.md#export) - Export Orders
* [Get](docs/sdks/orders/README.md#get) - Get Order
* [Update](docs/sdks/orders/README.md#update) - Update Order
* [Invoice](docs/sdks/orders/README.md#invoice) - Get Order Invoice
* [GenerateInvoice](docs/sdks/orders/README.md#generateinvoice) - Generate Order Invoice

### [OrganizationAccessTokens](docs/sdks/organizationaccesstokens/README.md)

* [List](docs/sdks/organizationaccesstokens/README.md#list) - List
* [Create](docs/sdks/organizationaccesstokens/README.md#create) - Create
* [Delete](docs/sdks/organizationaccesstokens/README.md#delete) - Delete
* [Update](docs/sdks/organizationaccesstokens/README.md#update) - Update

### [Organizations](docs/sdks/organizations/README.md)

* [List](docs/sdks/organizations/README.md#list) - List Organizations
* [Create](docs/sdks/organizations/README.md#create) - Create Organization
* [Get](docs/sdks/organizations/README.md#get) - Get Organization
* [Update](docs/sdks/organizations/README.md#update) - Update Organization

### [Payments](docs/sdks/payments/README.md)

* [List](docs/sdks/payments/README.md#list) - List Payments
* [Get](docs/sdks/payments/README.md#get) - Get Payment

### [Products](docs/sdks/products/README.md)

* [List](docs/sdks/products/README.md#list) - List Products
* [Create](docs/sdks/products/README.md#create) - Create Product
* [Get](docs/sdks/products/README.md#get) - Get Product
* [Update](docs/sdks/products/README.md#update) - Update Product
* [UpdateBenefits](docs/sdks/products/README.md#updatebenefits) - Update Product Benefits

### [Refunds](docs/sdks/refunds/README.md)

* [List](docs/sdks/refunds/README.md#list) - List Refunds
* [Create](docs/sdks/refunds/README.md#create) - Create Refund

### [Subscriptions](docs/sdks/subscriptions/README.md)

* [List](docs/sdks/subscriptions/README.md#list) - List Subscriptions
* [Create](docs/sdks/subscriptions/README.md#create) - Create Subscription
* [Export](docs/sdks/subscriptions/README.md#export) - Export Subscriptions
* [Get](docs/sdks/subscriptions/README.md#get) - Get Subscription
* [Revoke](docs/sdks/subscriptions/README.md#revoke) - Revoke Subscription
* [Update](docs/sdks/subscriptions/README.md#update) - Update Subscription

### [Webhooks](docs/sdks/webhooks/README.md)

* [ListWebhookEndpoints](docs/sdks/webhooks/README.md#listwebhookendpoints) - List Webhook Endpoints
* [CreateWebhookEndpoint](docs/sdks/webhooks/README.md#createwebhookendpoint) - Create Webhook Endpoint
* [GetWebhookEndpoint](docs/sdks/webhooks/README.md#getwebhookendpoint) - Get Webhook Endpoint
* [DeleteWebhookEndpoint](docs/sdks/webhooks/README.md#deletewebhookendpoint) - Delete Webhook Endpoint
* [UpdateWebhookEndpoint](docs/sdks/webhooks/README.md#updatewebhookendpoint) - Update Webhook Endpoint
* [ResetWebhookEndpointSecret](docs/sdks/webhooks/README.md#resetwebhookendpointsecret) - Reset Webhook Endpoint Secret
* [ListWebhookDeliveries](docs/sdks/webhooks/README.md#listwebhookdeliveries) - List Webhook Deliveries
* [RedeliverWebhookEvent](docs/sdks/webhooks/README.md#redeliverwebhookevent) - Redeliver Webhook Event

</details>
<!-- End Available Resources and Operations [operations] -->

<!-- Start Pagination [pagination] -->
## Pagination

Some of the endpoints in this SDK support pagination. To use pagination, you make your SDK calls as usual, but the
returned response object will have a `Next` method that can be called to pull down the next group of results. If the
return value of `Next` is `nil`, then there are no more pages to be fetched.

Here's an example of one such pagination call:
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
<!-- End Pagination [pagination] -->

<!-- Start Retries [retries] -->
## Retries

Some of the endpoints in this SDK support retries. If you use the SDK without any configuration, it will fall back to the default retry strategy provided by the API. However, the default retry strategy can be overridden on a per-operation basis, or across the entire SDK.

To change the default retry strategy for a single API call, simply provide a `retry.Config` object to the call by using the `WithRetries` option:
```go
package main

import (
	"context"
	pagogo "github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/retry"
	"log"
	"models/operations"
	"os"
)

func main() {
	ctx := context.Background()

	s := pagogo.New(
		pagogo.WithSecurity(os.Getenv("PAGO_ACCESS_TOKEN")),
	)

	res, err := s.Organizations.List(ctx, nil, pagogo.Pointer[int64](1), pagogo.Pointer[int64](10), nil, operations.WithRetries(
		retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: 1,
				MaxInterval:     50,
				Exponent:        1.1,
				MaxElapsedTime:  100,
			},
			RetryConnectionErrors: false,
		}))
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

If you'd like to override the default retry strategy for all operations that support retries, you can use the `WithRetryConfig` option at SDK initialization:
```go
package main

import (
	"context"
	pagogo "github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/retry"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	s := pagogo.New(
		pagogo.WithRetryConfig(
			retry.Config{
				Strategy: "backoff",
				Backoff: &retry.BackoffStrategy{
					InitialInterval: 1,
					MaxInterval:     50,
					Exponent:        1.1,
					MaxElapsedTime:  100,
				},
				RetryConnectionErrors: false,
			}),
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
<!-- End Retries [retries] -->

<!-- Start Error Handling [errors] -->
## Error Handling

Handling errors in this SDK should largely match your expectations. All operations return a response object or an error, they will never return both.

By Default, an API error will return `apierrors.APIError`. When custom error responses are specified for an operation, the SDK may also return their associated error. You can refer to respective *Errors* tables in SDK docs for more details on possible error types for each operation.

For example, the `List` function may return the following errors:

| Error Type                    | Status Code | Content Type     |
| ----------------------------- | ----------- | ---------------- |
| apierrors.HTTPValidationError | 422         | application/json |
| apierrors.APIError            | 4XX, 5XX    | \*/\*            |

### Example

```go
package main

import (
	"context"
	"errors"
	pagogo "github.com/pago-sh/pago-go"
	"github.com/pago-sh/pago-go/models/apierrors"
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

		var e *apierrors.HTTPValidationError
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}

		var e *apierrors.APIError
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}
	}
}

```
<!-- End Error Handling [errors] -->

<!-- Start Server Selection [server] -->
## Server Selection

### Select Server by Name

You can override the default server globally using the `WithServer(server string)` option when initializing the SDK client instance. The selected server will then be used as the default on the operations that use it. This table lists the names associated with the available servers:

| Name         | Server                         | Description            |
| ------------ | ------------------------------ | ---------------------- |
| `production` | `https://api.pago.sh`         | Production environment |
| `sandbox`    | `https://sandbox-api.pago.sh` | Sandbox environment    |

#### Example

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
		pagogo.WithServer("production"),
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

### Override Server URL Per-Client

The default server can also be overridden globally using the `WithServerURL(serverURL string)` option when initializing the SDK client instance. For example:
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
		pagogo.WithServerURL("https://api.pago.sh"),
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
<!-- End Server Selection [server] -->

<!-- Start Custom HTTP Client [http-client] -->
## Custom HTTP Client

The Go SDK makes API calls that wrap an internal HTTP client. The requirements for the HTTP client are very simple. It must match this interface:

```go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
```

The built-in `net/http` client satisfies this interface and a default client based on the built-in is provided by default. To replace this default with a client of your own, you can implement this interface yourself or provide your own client configured as desired. Here's a simple example, which adds a client with a 30 second timeout.

```go
import (
	"net/http"
	"time"

	"github.com/pago-sh/pago-go"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	sdkClient  = pagogo.New(pagogo.WithClient(httpClient))
)
```

This can be a convenient way to configure timeouts, cookies, proxies, custom headers, and other low-level configuration.
<!-- End Custom HTTP Client [http-client] -->

<!-- Start Special Types [types] -->
## Special Types

This SDK defines the following custom types to assist with marshalling and unmarshalling data.

### Date

`types.Date` is a wrapper around time.Time that allows for JSON marshaling a date string formatted as "2006-01-02".

#### Usage

```go
d1 := types.NewDate(time.Now()) // returns *types.Date

d2 := types.DateFromTime(time.Now()) // returns types.Date

d3, err := types.NewDateFromString("2019-01-01") // returns *types.Date, error

d4, err := types.DateFromString("2019-01-01") // returns types.Date, error

d5 := types.MustNewDateFromString("2019-01-01") // returns *types.Date and panics on error

d6 := types.MustDateFromString("2019-01-01") // returns types.Date and panics on error
```
<!-- End Special Types [types] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Maturity

This SDK is in beta, and there may be breaking changes between versions without a major version update. Therefore, we recommend pinning usage
to a specific package version. This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

## Contributions

While we value open-source contributions to this SDK, this library is generated programmatically. Any manual changes added to internal files will be overwritten on the next generation. 
We look forward to hearing your feedback. Feel free to open a PR or an issue with a proof of concept and we'll do our best to include it in a future release. 

### SDK Created by [Speakeasy](https://www.speakeasy.com/?utm_source=github-com/pago-sh/pago-go&utm_campaign=go)
