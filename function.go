package function

import (
	"context"
	"encoding/json"
	"log"
	"sort"

	"github.com/qushot/go-gae-app-version-rotate/shared"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type Body struct {
	ProjectID        string `json:"project_id"`
	ServiceName      string `json:"service_name"`
	KeepVersionCount int    `json:"keep_version_count"`
}

func GAEAppVersionRotate(ctx context.Context, m PubSubMessage) error {
	var body Body
	if err := json.Unmarshal(m.Data, &body); err != nil {
		log.Printf("request unmarshal error: %v", err)
		return err
	}

	aeAdminService, err := shared.NewAppEngineAdminService(ctx)
	if err != nil {
		return err
	}

	versions, err := aeAdminService.VersionList(ctx, body.ProjectID, body.ServiceName)
	if err != nil {
		return err
	}

	createTimes := make([]string, len(versions))
	createTimeVersionMap := make(map[string]string)
	for i, version := range versions {
		createTimes[i] = version.CreateTime
		createTimeVersionMap[version.CreateTime] = version.Id
	}

	sort.Strings(createTimes)

	deletedCount := 0
	for _, createTime := range createTimes {
		if len(versions) <= deletedCount+body.KeepVersionCount {
			break
		}

		versionName := createTimeVersionMap[createTime]
		if err := aeAdminService.DeleteVersion(ctx, body.ProjectID, body.ServiceName, versionName); err != nil {
			continue
		}
		deletedCount++
	}

	return nil
}
