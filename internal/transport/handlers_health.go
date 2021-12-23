package transport

import (
	"net/http"

	"github.com/skandyla/s3-uploader/internal/models"
	"github.com/skandyla/s3-uploader/internal/version"
)

func (h *Handler) liveness(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//err := h.healthService.Ping(ctx)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pgDep, _ := h.healthService.Info(ctx)

	resp := &models.Info{
		Name: "microtester",
		Build: models.InfoBuild{
			Version:   version.BuildVersion,
			Date:      version.BuildTime,
			Branch:    version.BuildBranch,
			Commit:    version.BuildCommit,
			GoVersion: version.GoVersion,
		},
		Dependencies: map[string]models.InfoDependencyItem{
			"postgres": pgDep,
		},
	}

	statusCode := http.StatusOK

	if pgDep.Status != 200 {
		statusCode = http.StatusServiceUnavailable
	}

	respondWithJSON(w, statusCode, resp)
}
