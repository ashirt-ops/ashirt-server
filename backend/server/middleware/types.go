// Copyright 2023, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package middleware

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler
