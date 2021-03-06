// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package correlation

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts"
)

func ManageHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header.Get(contracts.CorrelationHeader)
		if hdr == "" {
			hdr = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), contracts.CorrelationHeader, hdr)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func OnResponseComplete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		next.ServeHTTP(w, r)
		correlationId := IdFromContext(r.Context())
		lc := LoggingClientFromContext(r.Context())
		if lc != nil {
			lc.Trace("Response complete", contracts.CorrelationHeader, correlationId, "duration", time.Since(begin).String())
		}
	})
}

func OnRequestBegin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationId := IdFromContext(r.Context())
		lc := LoggingClientFromContext(r.Context())
		if lc != nil {
			lc.Trace("Begin request", contracts.CorrelationHeader, correlationId, "path", r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}
