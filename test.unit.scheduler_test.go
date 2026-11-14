package scheduler

import "testing"

func TestScheduleSubjects(t *testing.T) {
	if SubjectCreate != "scheduler.schedules.create" || SubjectUpdate != "scheduler.schedules.update" || SubjectDelete != "scheduler.schedules.delete" {
		t.Fatal("unexpected schedule subjects")
	}
}
