// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services

import (
	"context"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/policy"
	"github.com/theparanoids/ashirt-server/backend/server/middleware"

	sq "github.com/Masterminds/squirrel"
)

type CreateServiceWorkerInput struct {
	Name        string
	ServiceType string
	Config      string
}

type UpdateServiceWorkerInput struct {
	ID          int64
	Name        string
	ServiceType string
	Config      string
}

type DeleteServiceWorkerInput struct {
	ID       int64
	DoDelete bool
}

func ListServiceWorker(ctx context.Context, db *database.Connection) ([]*dtos.ServiceWorker, error) {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return nil, backend.WrapError("Insufficient access to view service workers", backend.UnauthorizedReadErr(err))
	}

	var services []models.ServiceWorker
	err := db.Select(&services,
		sq.Select("*").From("service_workers"),
	)

	if err != nil {
		return nil, backend.WrapError("Could not create a service worker", backend.DatabaseErr(err))
	}

	serviceWorkersDTO := make([]*dtos.ServiceWorker, len(services))
	for i, v := range services {
		serviceWorkersDTO[i] = &dtos.ServiceWorker{
			ID:      v.ID,
			Name:    v.Name,
			Config:  v.Config,
			Deleted: v.DeletedAt != nil,
		}
	}

	return serviceWorkersDTO, nil
}

func CreateServiceWorker(ctx context.Context, db *database.Connection, i CreateServiceWorkerInput) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Insufficient access to create a service worker", backend.UnauthorizedWriteErr(err))
	}

	_, err := db.Insert("service_workers", map[string]interface{}{
		"name":   i.Name,
		"config": i.Config,
	})

	if err != nil {
		return backend.WrapError("Could not create a service worker", backend.DatabaseErr(err))
	}

	return nil
}

func UpdateServiceWorker(ctx context.Context, db *database.Connection, i UpdateServiceWorkerInput) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Insufficient access to update the service worker", backend.UnauthorizedWriteErr(err))
	}

	err := db.Update(
		sq.Update("service_workers").
			SetMap(map[string]interface{}{
				"name":   i.Name,
				"config": i.Config,
			}).
			Where(sq.Eq{"id": i.ID}),
	)

	if err != nil {
		return backend.WrapError("Could not update the service worker", backend.DatabaseErr(err))
	}

	return nil
}

func DeleteServiceWorker(ctx context.Context, db *database.Connection, i DeleteServiceWorkerInput) error {
	if err := policy.Require(middleware.Policy(ctx), policy.AdminUsersOnly{}); err != nil {
		return backend.WrapError("Insufficient access to create a service worker", backend.UnauthorizedWriteErr(err))
	}
	query := sq.Update("service_workers").Where(sq.Eq{"id": i.ID})

	if i.DoDelete {
		query = query.Set("deleted_at", time.Now())
	} else {
		query = query.Set("deleted_at", nil)
	}

	err := db.Update(query)
	if err != nil {
		return backend.WrapError("Could not delete the service worker", backend.DatabaseErr(err))
	}

	return nil
}
