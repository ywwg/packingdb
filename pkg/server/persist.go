package server

import (
	"log"
	"time"
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

	for _, name := range dirtyNames {
		s.mu.RLock()
		trip, ok := s.trips[name]
		s.mu.RUnlock()

		if !ok {
			continue
		}

		filename := s.findTripFile(name)
		if filename == "" {
			continue
		}

		if err := trip.SaveToFile(filename); err != nil {
			log.Printf("Error saving trip %s: %v", name, err)
			continue
		}

		s.mu.Lock()
		delete(s.dirtyTrips, name)
		s.mu.Unlock()
	}
}

// markDirty marks a trip as having unsaved changes
func (s *Server) markDirty(name string) {
	s.mu.Lock()
	s.dirtyTrips[name] = true
	s.mu.Unlock()
}
