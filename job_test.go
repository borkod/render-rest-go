package render

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestJob(t *testing.T) {
	mux := chi.NewMux()

	mux.Route("/v1/services", func(r chi.Router) {

		r.Route("/{serviceID}", func(r chi.Router) {
			r.Use(ServiceCtx)
			r.Route("/jobs", func(r chi.Router) {
				r.Get("/{jobID}", getJob)
				r.Post("/", createJob)
			})
		})
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := &Client{
		HttpClient:   ts.Client(),
		ApiKey:       "test",
		Host:         ts.URL,
		ServicesBase: "/v1/services/",
	}

	t.Run("Test GetJob", func(t *testing.T) {
		job, err := client.GetJob("job-cc8de9sgqg4fs9cvbg4g", "srv-cblfck319n07sfpr4d2g")
		require.NoError(t, err)
		t.Log(job)
	})

	t.Run("Test CreateJob with no plan id", func(t *testing.T) {
		nj := NewJob{
			ServiceId:    "srv-cblfck319n07sfpr4d2g",
			StartCommand: "echo 'hello world'",
		}
		job, err := client.CreateJob(nj)
		require.NoError(t, err)
		t.Log(job)
	})

	t.Run("Test CreateJob with no plan id", func(t *testing.T) {
		nj := NewJob{
			ServiceId:    "srv-cblfck319n07sfpr4d2g",
			PlanId:       "plan-srv-001",
			StartCommand: "echo 'hello world'",
		}
		job, err := client.CreateJob(nj)
		require.NoError(t, err)
		t.Log(job)
	})
}

func ServiceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceID := chi.URLParam(r, "serviceID")
		ctx := context.WithValue(r.Context(), "serviceID", serviceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceID, ok := ctx.Value("serviceID").(string)
	jobID := chi.URLParam(r, "jobID")
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	result := `{"id": "` + jobID + `","serviceId": "` + serviceID + `","startCommand": "python terraso_backend/manage.py clean_up_deleted_files","planId": "plan-srv-006","createdAt": "2022-09-01T16:00:39.627588Z","startedAt": "2022-09-01T16:00:39Z","finishedAt": "2022-09-01T16:09:13Z","status": "succeeded"}`
	w.Write([]byte(result))
}

func createJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceID, ok := ctx.Value("serviceID").(string)
	jobID := "job-cc8de9sgqg4fs9cvbg4g"
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "{\"message\": \"Error reading body\"}", http.StatusInternalServerError)
	}

	nj := NewJob{}
	err = json.Unmarshal(bytes, &nj)
	if err != nil {
		http.Error(w, "{\"message\": \"Error unmarshalling body\"}", http.StatusInternalServerError)
	}

	if len(nj.PlanId) == 0 {
		nj.PlanId = "plan-srv-006"
	}

	result := `{"id": "` + jobID + `","serviceId": "` + serviceID + `","startCommand": "` + nj.StartCommand + `","planId": "` + nj.PlanId + `","createdAt": "` + time.Now().Format(time.RFC3339) + `","startedAt": "","finishedAt": "","status": ""}`
	w.Write([]byte(result))
}
