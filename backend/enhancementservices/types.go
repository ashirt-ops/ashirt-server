package enhancementservices

import "github.com/theparanoids/ashirt-server/backend/models"

type BasicPayload struct {
	Type string `json:"type"`
}

type Payload struct {
	BasicPayload
	EvidenceUUID  string `json:"evidenceUuid"  db:"uuid"`
	OperationSlug string `json:"operationSlug" db:"operation_slug"`
	ContentType   string `json:"contentType"   db:"content_type"`
}

type WorkerHandler = func(workerName string, evidenceID int64, configText []byte, payload *Payload) (*models.EvidenceMetadata, error)

type BasicServiceWorkerConfig struct {
	Type    string `json:"type"`
	Version int64  `json:"version"`
}

type ServiceWorker interface {
	Build(workerName string, evidenceID int64, config []byte) error
	Test() (string, error)
	Process(payload *Payload) (*models.EvidenceMetadata, error)
}

