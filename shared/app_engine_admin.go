package shared

import (
	"context"
	"log"

	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/option"
)

type AppEngineAdminService struct {
	srv *appengine.APIService
}

func NewAppEngineAdminService(ctx context.Context) (*AppEngineAdminService, error) {
	srv, err := appengine.NewService(ctx, option.WithScopes(appengine.CloudPlatformScope))
	if err != nil {
		log.Printf("[ERROR] app engine admin api new service error: %v", err)
		return nil, err
	}

	return &AppEngineAdminService{srv}, nil
}

func (ae *AppEngineAdminService) VersionList(ctx context.Context, projectID, serviceID string) ([]*appengine.Version, error) {
	resp, err := ae.srv.Apps.Services.Versions.List(projectID, serviceID).Context(ctx).View("BASIC").PageSize(210).Do()
	if err != nil {
		log.Printf("[ERROR] app engine admin api versions list error(service=%s): %v", serviceID, err)
		return nil, err
	}

	return resp.Versions, nil
}

func (ae *AppEngineAdminService) DeleteVersion(ctx context.Context, projectID, serviceID, versionName string) error {
	if _, err := ae.srv.Apps.Services.Versions.Delete(projectID, serviceID, versionName).Context(ctx).Do(); err != nil {
		log.Printf("[ERROR] app engine admin api versions delete error(service=%s, version=%s): %v", serviceID, versionName, err)
		return err
	}
	return nil
}
