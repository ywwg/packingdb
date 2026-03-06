package server

import (
	"time"

	"github.com/ywwg/packingdb/pkg/logger"
)

// StartBackgroundPersist starts the background goroutine that periodically saves dirty trips
func (s *Server) StartBackgroundPersist() {
	go func() {
		ticker := time.NewTicker(s.persistInterval)
		defer ticker.Stop()
		defer close(s.persistDone)

		for {
			select {
			case <-ticker.C:
				s.persistDirtyTrips()
			case <-s.stopPersist:
				// Final save before shutdown
				s.persistDirtyTrips()
				return
			}
		}
	}()
}

// Shutdown stops the background persistence and saves all dirty trips
func (s *Server) Shutdown() {
	close(s.stopPersist)
	<-s.persistDone // wait for final persist to complete
}

// persistDirtyTrips saves all trips that have been modified
func (s *Server) persistDirtyTrips() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.dirtyTrips) > 0 {
		logger.Debug("Persisting dirty trips", "count", len(s.dirtyTrips))
	}

	for name := range s.dirtyTrips {
		trip, ok := s.trips[name]
		filename := s.nameToFile[name]
		if !ok || filename == "" {
			logger.Warn("Trip not found for persist", "trip", name)
			continue
		}

		if err := trip.SaveToFile(filename); err != nil {
			logger.Error("Failed to save trip", "trip", name, "file", filename, "error", err)
			continue
		}

		delete(s.dirtyTrips, name)
		logger.Info("Persisted trip to disk", "trip", name, "file", filename)
	}
}

