package appointmentapp

import "time"

// calculateEndTime calcula la hora de fin dado un start "HH:MM" y duración en minutos.
func calculateEndTime(startTime string, durationMinutes int) string {
	t, err := time.Parse("15:04", startTime)
	if err != nil {
		return startTime
	}
	return t.Add(time.Duration(durationMinutes) * time.Minute).Format("15:04")
}
