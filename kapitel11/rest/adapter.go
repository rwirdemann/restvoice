package rest

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rwirdemann/restvoice/kapitel06/usecase"

	"github.com/gorilla/mux"
)

type Adapter struct {
	r *mux.Router
}

func NewAdapter() *Adapter {
	return &Adapter{mux.NewRouter()}
}

func (a Adapter) ListenAndServe() {
	log.Printf("Listening on http://0.0.0.0%s\n", ":8080")
	http.ListenAndServe(":8080", a.r)
}

func (a Adapter) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	return a.r.NewRoute().Path(path).HandlerFunc(f)
}

func (a Adapter) MakeGetActivitiesHandler(usecase usecase.GetActivities) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=0")
		activities := usecase.Run(extractUserIDFromHeader(r.Header))
		cacheableActivities := ActivitiesPresenter{}.Present(activities)

		// Test if client wants a full refresh
		cacheControl := r.Header.Get("Cache-Control")
		if strings.Contains(cacheControl, "no-cache") {
			w.Header().Set("Last-Modified", cacheableActivities.LastModified.Format(layout))
			w.Write(cacheableActivities.Activities)
		}

		// Cache logic: return 304 if nothing has changed since "Last-Modified-Since"
		lastModifiedSince := r.Header.Get("Last-Modified-Since")
		if lastModifiedSince != "" {
			t, err := time.Parse(layout, lastModifiedSince)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if truncateToSeconds(t).Equal(truncateToSeconds(cacheableActivities.LastModified)) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.Header().Set("Last-Modified", cacheableActivities.LastModified.Format(layout))
		w.Write(cacheableActivities.Activities)
	}
}

func extractUserIDFromHeader(h http.Header) string {
	return "1234"
}

func truncateToSeconds(t time.Time) time.Time {
	return t.Truncate(time.Duration(time.Second))
}

const layout = "Mon, _2 Jan 2006 15:04:05 GMT"
