package appointmentapp

import "time"

func calculateEndTime(startTime string, durationMinutes int) string {
	t, err := time.Parse("15:04", startTime)
	if err != nil {
		return startTime
	}
	return t.Add(time.Duration(durationMinutes) * time.Minute).Format("15:04")
}
