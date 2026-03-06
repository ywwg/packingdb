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
	// Copy dirty trip names so we can release lock while saving
	dirtyNames := make([]string, 0, len(s.dirtyTrips))
	for name := range s.dirtyTrips {
		dirtyNames = append(dirtyNames, name)
	}
	s.mu.Unlock()

	if len(dirtyNames) > 0 {
		logger.Debug("Persisting dirty trips", "count", len(dirtyNames))
	}

	for _, name := range dirtyNames {
		// Hold the lock during the save so that concurrent trip mutations
		// (which also hold s.mu) cannot corrupt the trip data while it is
		// being serialized to disk, and so that clearing the dirty flag is
		// atomic with the save (preventing a dirty signal set after the save
		// started from being silently dropped).
		s.mu.Lock()
		trip, ok := s.trips[name]
		filename := s.nameToFile[name]
		if !ok || filename == "" {
			s.mu.Unlock()
			logger.Warn("Trip not found for persist", "trip", name)
			continue
		}

		if err := trip.SaveToFile(filename); err != nil {
			s.mu.Unlock()
			logger.Error("Failed to save trip", "trip", name, "file", filename, "error", err)
			continue
		}

		delete(s.dirtyTrips, name)
		s.mu.Unlock()

		logger.Info("Persisted trip to disk", "trip", name, "file", filename)
	}
}

