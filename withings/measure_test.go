package withings

import "testing"

func TestMeasureType(t *testing.T) {
	t.Run("AllValid", func(t *testing.T) {
		for _, v := range AllMeasureTypes() {
			if !v.IsValid() {
				t.Errorf("%d is supposed to be a valid MeasureType", v)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if MeasureType(0).IsValid() {
			t.Error("non existent MeasureType should not be valid")
		}
	})
}

func TestMeasureCategory(t *testing.T) {
	t.Run("AllValid", func(t *testing.T) {
		for _, v := range []MeasureCategory{MeasureCategoryRealMeasure, MeasureCategoryUserObjective} {
			if !v.IsValid() {
				t.Errorf("%d is supposed to be a valid MeasureCategory", v)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if MeasureCategory(0).IsValid() {
			t.Error("non existent MeasureCategory should not be valid")
		}
	})
}

func TestActivityField(t *testing.T) {
	t.Run("AllValid", func(t *testing.T) {
		for _, v := range AllActivityFields() {
			if !v.IsValid() {
				t.Errorf("%s is supposed to be a valid ActivityField", v)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if ActivityField("invalid").IsValid() {
			t.Error("non existent ActivityField should not be valid")
		}
	})
}

func TestIntradayActivityField(t *testing.T) {
	t.Run("AllValid", func(t *testing.T) {
		for _, v := range AllIntradayActivityFields() {
			if !v.IsValid() {
				t.Errorf("%s is supposed to be a valid IntradayActivityField", v)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if IntradayActivityField("invalid").IsValid() {
			t.Error("non existent IntradayActivityField should not be valid")
		}
	})
}

func TestWorkoutField(t *testing.T) {
	t.Run("AllValid", func(t *testing.T) {
		for _, v := range AllWorkoutFields() {
			if !v.IsValid() {
				t.Errorf("%s is supposed to be a valid WorkoutField", v)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if WorkoutField("invalid").IsValid() {
			t.Error("non existent WorkoutField should not be valid")
		}
	})
}
