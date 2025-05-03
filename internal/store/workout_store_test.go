package store

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=workout_db port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// run migrations for test database
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations for test database: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}

	return db
}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "Leg day workout",
				Description:     "Test Description",
				DurationMinutes: 60,
				CaloriesBurned:  450,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Squats",
						Sets:         4,
						Reps:         IntPtr(12),
						Weight:       FloatPtr(90),
						Notes:        "Test Notes",
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid workout",
			workout: &Workout{
				Title:           "Full body workout",
				Description:     "complete workout",
				DurationMinutes: 30,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Plank",
						Sets:         3,
						Reps:         IntPtr(60),
						Notes:        "Test Notes",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Squats",
						Sets:            4,
						Reps:            IntPtr(12),
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(135.5),
						Notes:           "Test Notes",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			require.NoError(t, err)

			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range tt.workout.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrieved.Entries[i].Reps)
				// assert.Equal(t, tt.workout.Entries[i].Weight, retrieved.Entries[i].Weight)
				// assert.Equal(t, tt.workout.Entries[i].Notes, retrieved.Entries[i].Notes)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}
