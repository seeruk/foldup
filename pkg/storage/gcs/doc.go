// Package gcs provides interfaces and wraps calls for interacting with Google Cloud Storage. This
// package is mainly necessary because of the lack of interfaces provided by Google with their
// client library, resulting in ballooning code structures that are nowhere near as easy to follow
// (but actually might result in a simpler, more focused API where they're used).
//
// The package aims to be unit-testable without calling any external services, making storage
// gateways easier to test.
package gcs
