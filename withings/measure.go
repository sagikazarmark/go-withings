package withings

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// MeasureService handles communication with the measure related
// methods of the Withings API.
//
// Withings API docs: https://developer.withings.com/api-reference#tag/measure
type MeasureService service

// MeasureGetOptions specifies parameters for various Measure related operations
// that support date based filters and/or pagination.
//
// Withings API docs: https://developer.withings.com/api-reference#tag/measure
type MeasureGetOptions struct {
	// Use StartDate and EndDate for a date range query.
	//
	// Mutually exclusive with LastUpdate!
	StartDate time.Time
	EndDate   time.Time

	// Use LastUpdate to query new values.
	//
	// Mutually exclusive with StartDate and EndDate!
	LastUpdate time.Time

	// Offset retrieves the next batch from the resultset.
	Offset int
}

// MeasureType is is a metric that Withings devices track.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type MeasureType int

// MeasureType values
const (
	MeasureTypeWeight         MeasureType = 1   // Weight (kg)
	MeasureTypeHeight         MeasureType = 4   // Height (meter)
	MeasureTypeFatFreeMass    MeasureType = 5   // Fat Free Mass (kg)
	MeasureTypeFatRatio       MeasureType = 6   // Fat Ratio (%)
	MeasureTypeFatMassWeight  MeasureType = 8   // Fat Mass Weight (kg)
	MeasureTypeDiastolicBP    MeasureType = 9   // Diastolic Blood Pressure (mmHg)
	MeasureTypeSystolicBP     MeasureType = 10  // Systolic Blood Pressure (mmHg)
	MeasureTypeHeartPulse     MeasureType = 11  // Heart Pulse (bpm)
	MeasureTypeTemp           MeasureType = 12  // Temperature (celsius)
	MeasureTypeSpO2           MeasureType = 54  // SpO2 (%)
	MeasureTypeBodyTemp       MeasureType = 71  // Body Temperature (celsius)
	MeasureTypeSkinTemp       MeasureType = 73  // Skin temperature (celsius)
	MeasureTypeMuscleMass     MeasureType = 76  // Muscle Mass (kg)
	MeasureTypeHydration      MeasureType = 77  // Hydration (kg)
	MeasureTypeBoneMass       MeasureType = 88  // Bone Mass (kg)
	MeasureTypePWaveVel       MeasureType = 91  // Pulse Wave Velocity (m/s)
	MeasureTypeVO2Max         MeasureType = 123 // VO2 max is a numerical measurement of your body’s ability to consume oxygen (ml/min/kg).
	MeasureTypeQRSInterval    MeasureType = 135 // QRS interval duration based on ECG signal
	MeasureTypePRInterval     MeasureType = 136 // PR interval duration based on ECG signal
	MeasureTypeQTInterval     MeasureType = 137 // QT interval duration based on ECG signal
	MeasureTypeCorrQTInterval MeasureType = 138 // Corrected QT interval duration based on ECG signal
	MeasureTypeAtrialFib      MeasureType = 139 // Atrial fibrillation result from PPG
)

var validMeasureTypeValues = map[MeasureType]struct{}{
	MeasureTypeWeight:         {},
	MeasureTypeHeight:         {},
	MeasureTypeFatFreeMass:    {},
	MeasureTypeFatRatio:       {},
	MeasureTypeFatMassWeight:  {},
	MeasureTypeDiastolicBP:    {},
	MeasureTypeSystolicBP:     {},
	MeasureTypeHeartPulse:     {},
	MeasureTypeTemp:           {},
	MeasureTypeSpO2:           {},
	MeasureTypeBodyTemp:       {},
	MeasureTypeSkinTemp:       {},
	MeasureTypeMuscleMass:     {},
	MeasureTypeHydration:      {},
	MeasureTypeBoneMass:       {},
	MeasureTypePWaveVel:       {},
	MeasureTypeVO2Max:         {},
	MeasureTypeQRSInterval:    {},
	MeasureTypePRInterval:     {},
	MeasureTypeQTInterval:     {},
	MeasureTypeCorrQTInterval: {},
	MeasureTypeAtrialFib:      {},
}

// IsValid checks if v is a valid MeasureType.
func (v MeasureType) IsValid() bool {
	_, ok := validMeasureTypeValues[v]

	return ok
}

// AllMeasureTypes returns the list of all MeasureType values.
func AllMeasureTypes() []MeasureType {
	return []MeasureType{
		MeasureTypeWeight,
		MeasureTypeHeight,
		MeasureTypeFatFreeMass,
		MeasureTypeFatRatio,
		MeasureTypeFatMassWeight,
		MeasureTypeDiastolicBP,
		MeasureTypeSystolicBP,
		MeasureTypeHeartPulse,
		MeasureTypeTemp,
		MeasureTypeSpO2,
		MeasureTypeBodyTemp,
		MeasureTypeSkinTemp,
		MeasureTypeMuscleMass,
		MeasureTypeHydration,
		MeasureTypeBoneMass,
		MeasureTypePWaveVel,
		MeasureTypeVO2Max,
		MeasureTypeQRSInterval,
		MeasureTypePRInterval,
		MeasureTypeQTInterval,
		MeasureTypeCorrQTInterval,
		MeasureTypeAtrialFib,
	}
}

// MeasureCategory differentiates between real measurements and user objectives.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type MeasureCategory int

// MeasureCategory values
const (
	MeasureCategoryRealMeasure   MeasureCategory = 1 // Real measures
	MeasureCategoryUserObjective MeasureCategory = 2 // User objectives
)

var validMeasureCategoryValues = map[MeasureCategory]struct{}{
	MeasureCategoryRealMeasure:   {},
	MeasureCategoryUserObjective: {},
}

// IsValid checks if v is a valid MeasureCategory.
func (v MeasureCategory) IsValid() bool {
	_, ok := validMeasureCategoryValues[v]

	return ok
}

type getmeasResponse struct {
	Body Measures `json:"body"`
}

// Measures is the response from the Getmeas API call.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type Measures struct {
	UpdateTime    int            `json:"updatetime"` // Note: spec says string, but it's in fact an int
	TimeZone      string         `json:"timezone"`
	MeasureGroups []MeasureGroup `json:"measuregrps"`
}

// Measures are returned in groups.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type MeasureGroup struct {
	GroupID   int64           `json:"grpid"`
	Attrib    int             `json:"attrib"`
	Date      int             `json:"date"`
	CreatedAt int             `json:"created"`
	Category  MeasureCategory `json:"category"`
	DeviceID  string          `json:"deviceid"`
	Measures  []Measure       `json:"measures"`
	Comment   string          `json:"comment"` // Deprecated
}

// Measure is an individual data point.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type Measure struct {
	Value int         `json:"value"`
	Type  MeasureType `json:"type"`
	Unit  int         `json:"unit"`
	Algo  int         `json:"algo"` // Deprecated
	FM    int         `json:"fm"`   // Deprecated
	FW    int         `json:"fw"`   // Deprecated
}

// Getmeas provides measures stored on a specific date.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
func (s *MeasureService) Getmeas(ctx context.Context, measureTypes []MeasureType, category MeasureCategory, opts MeasureGetOptions) (*Measures, *Response, error) {
	// validate category first because it requires less effort
	if !category.IsValid() {
		return nil, nil, errors.New("invalid category")
	}

	measureTypes = filterValidMeasureTypeValues(measureTypes)

	if len(measureTypes) == 0 {
		return nil, nil, errors.New("need at least one measure type")
	}

	const urlPath = "measure"

	form := url.Values{
		"action":   {"getmeas"},
		"category": {fmt.Sprintf("%d", category)},
	}

	if len(measureTypes) == 1 {
		form.Add("meastype", fmt.Sprintf("%d", measureTypes[0]))
	} else {
		form.Add("meastypes", strings.Join(measureTypesToString(measureTypes), ","))
	}

	if !opts.LastUpdate.IsZero() {
		form.Add("lastupdate", fmt.Sprintf("%d", opts.LastUpdate.Unix()))
	} else if !opts.StartDate.IsZero() && !opts.EndDate.IsZero() {
		form.Add("startdate", fmt.Sprintf("%d", opts.StartDate.Unix()))
		form.Add("enddate", fmt.Sprintf("%d", opts.EndDate.Unix()))
	}

	if opts.Offset > 0 {
		form.Add("offset", fmt.Sprintf("%d", opts.Offset))
	}

	measuresResp := new(getmeasResponse)

	resp, err := s.client.PostForm(ctx, urlPath, form, measuresResp)

	return &measuresResp.Body, resp, err
}

func filterValidMeasureTypeValues(values []MeasureType) []MeasureType {
	var validValues []MeasureType

	for _, v := range values {
		if !v.IsValid() {
			continue
		}

		validValues = append(validValues, v)
	}

	return validValues
}

func measureTypesToString(measureTypes []MeasureType) []string {
	s := make([]string, 0, len(measureTypes))

	for _, v := range measureTypes {
		s = append(s, fmt.Sprintf("%d", v))
	}

	return s
}

// ActivityField is a type of metric tracked during an activity.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getactivity
type ActivityField string

// ActivityField values
const (
	ActivityFieldSteps         ActivityField = "steps"         // Number of steps.
	ActivityFieldDistance      ActivityField = "distance"      // Distance travelled (in meters).
	ActivityFieldElevation     ActivityField = "elevation"     // Number of floors climbed.
	ActivityFieldSoft          ActivityField = "soft"          // Duration of soft activities (in seconds).
	ActivityFieldModerate      ActivityField = "moderate"      // Duration of moderate activities (in seconds).
	ActivityFieldIntense       ActivityField = "intense"       // Duration of intense activities (in seconds).
	ActivityFieldActive        ActivityField = "active"        // Sum of intense and moderate activity durations (in seconds).
	ActivityFieldCalories      ActivityField = "calories"      // Active calories burned (in Kcal).
	ActivityFieldTotalCalories ActivityField = "totalcalories" // Total calories burned (in Kcal).
	ActivityFieldHRAverage     ActivityField = "hr_average"    // Average heart rate.
	ActivityFieldHRMin         ActivityField = "hr_min"        // Minimal heart rate.
	ActivityFieldHRMax         ActivityField = "hr_max"        // Maximal heart rate.
	ActivityFieldHRZone0       ActivityField = "hr_zone_0"     // Duration in seconds when heart rate was in a light zone.
	ActivityFieldHRZone1       ActivityField = "hr_zone_1"     // Duration in seconds when heart rate was in a moderate zone.
	ActivityFieldHRZone2       ActivityField = "hr_zone_2"     // Duration in seconds when heart rate was in an intense zone.
	ActivityFieldHRZone3       ActivityField = "hr_zone_3"     // Duration in seconds when heart rate was in maximal zone.
)

var validActivityFieldValues = map[ActivityField]struct{}{
	ActivityFieldSteps:         {},
	ActivityFieldDistance:      {},
	ActivityFieldElevation:     {},
	ActivityFieldSoft:          {},
	ActivityFieldModerate:      {},
	ActivityFieldIntense:       {},
	ActivityFieldActive:        {},
	ActivityFieldCalories:      {},
	ActivityFieldTotalCalories: {},
	ActivityFieldHRAverage:     {},
	ActivityFieldHRMin:         {},
	ActivityFieldHRMax:         {},
	ActivityFieldHRZone0:       {},
	ActivityFieldHRZone1:       {},
	ActivityFieldHRZone2:       {},
	ActivityFieldHRZone3:       {},
}

// IsValid checks if v is a valid ActivityField.
func (v ActivityField) IsValid() bool {
	_, ok := validActivityFieldValues[v]

	return ok
}

// AllActivityFields is the list of all ActivityField values.
func AllActivityFields() []ActivityField {
	return []ActivityField{
		ActivityFieldSteps,
		ActivityFieldDistance,
		ActivityFieldElevation,
		ActivityFieldSoft,
		ActivityFieldModerate,
		ActivityFieldIntense,
		ActivityFieldActive,
		ActivityFieldCalories,
		ActivityFieldTotalCalories,
		ActivityFieldHRAverage,
		ActivityFieldHRMin,
		ActivityFieldHRMax,
		ActivityFieldHRZone0,
		ActivityFieldHRZone1,
		ActivityFieldHRZone2,
		ActivityFieldHRZone3,
	}
}

type getactivityResponse struct {
	Body Activities `json:"body"`
}

// Activities is the response from the Getactivity API call.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getactivity
type Activities struct {
	Activities []Activity `json:"activities"`
}

// Activity aggregates metrics of a single activity.
//
// Fields are populated based on the requested fields.
//
// TODO: consider making fields pointers, so they dont get populated when no data is returned.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getactivity
type Activity struct {
	Date      string `json:"date"`
	Timezone  string `json:"timezone"`
	DeviceID  string `json:"deviceid"`
	Brand     int    `json:"brand"`
	IsTracker bool   `json:"is_tracker"`

	// Fields
	Steps         int     `json:"steps"`
	Distance      float64 `json:"distance"`  // Note: spec says int, but it's in fact a float
	Elevation     float64 `json:"elevation"` // Note: spec says int, but it's in fact a float
	Soft          int     `json:"soft"`
	Moderate      int     `json:"moderate"`
	Intense       int     `json:"intense"`
	Active        int     `json:"active"`
	Calories      float64 `json:"calories"`      // Note: spec says int, but it's in fact a float
	TotalCalories float64 `json:"totalcalories"` // Note: spec says int, but it's in fact a float
	HRAverage     int     `json:"hr_average"`
	HRMin         int     `json:"hr_min"`
	HRMax         int     `json:"hr_max"`
	HRZone0       int     `json:"hr_zone_0"`
	HRZone1       int     `json:"hr_zone_1"`
	HRZone2       int     `json:"hr_zone_2"`
	HRZone3       int     `json:"hr_zone_3"`
}

// Getactivity provides daily aggregated activity data of a user.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getactivity
func (s *MeasureService) Getactivity(ctx context.Context, fields []ActivityField, opts MeasureGetOptions) (*Activities, *Response, error) {
	fields = filterValidActivityFieldValues(fields)

	if len(fields) == 0 {
		return nil, nil, errors.New("need at least one activity field")
	}

	const urlPath = "v2/measure"

	form := url.Values{
		"action":      {"getactivity"},
		"data_fields": {joinActivityFields(fields)},
	}

	if !opts.LastUpdate.IsZero() {
		form.Add("lastupdate", fmt.Sprintf("%d", opts.LastUpdate.Unix()))
	} else if !opts.StartDate.IsZero() && !opts.EndDate.IsZero() {
		form.Add("startdateymd", opts.StartDate.Format("2006-01-02"))
		form.Add("enddateymd", opts.EndDate.Format("2006-01-02"))
	}

	if opts.Offset > 0 {
		form.Add("offset", fmt.Sprintf("%d", opts.Offset))
	}

	activityResp := new(getactivityResponse)

	resp, err := s.client.PostForm(ctx, urlPath, form, activityResp)

	return &activityResp.Body, resp, err
}

func filterValidActivityFieldValues(values []ActivityField) []ActivityField {
	var validValues []ActivityField

	for _, v := range values {
		if !v.IsValid() {
			continue
		}

		validValues = append(validValues, v)
	}

	return validValues
}

func joinActivityFields(fields []ActivityField) string {
	s := make([]string, 0, len(fields))

	for _, f := range fields {
		s = append(s, string(f))
	}

	return strings.Join(s, ",")
}

// IntradayActivityField is a type of metric tracked during an activity.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getintradayactivity
type IntradayActivityField string

// IntradayActivityField values
const (
	IntradayActivityFieldSteps     IntradayActivityField = "steps"      // Number of steps.
	IntradayActivityFieldElevation IntradayActivityField = "elevation"  // Number of floors climbed.
	IntradayActivityFieldCalories  IntradayActivityField = "calories"   // Active calories burned (in Kcal).
	IntradayActivityFieldDistance  IntradayActivityField = "distance"   // Distance travelled (in meters).
	IntradayActivityFieldStroke    IntradayActivityField = "stroke"     // Number of strokes performed.
	IntradayActivityFieldPoolLap   IntradayActivityField = "pool_lap"   // Number of pool_lap performed.
	IntradayActivityFieldDuration  IntradayActivityField = "duration"   // Duration of the activity (in seconds).
	IntradayActivityFieldHeartRate IntradayActivityField = "heart_rate" // Measured heart rate.
	IntradayActivityFieldSpO2Auto  IntradayActivityField = "spo2_auto"  // SpO2 measurement automatically tracked by a device tracker.
)

var validIntradayActivityFieldValues = map[IntradayActivityField]struct{}{
	IntradayActivityFieldSteps:     {},
	IntradayActivityFieldElevation: {},
	IntradayActivityFieldCalories:  {},
	IntradayActivityFieldDistance:  {},
	IntradayActivityFieldStroke:    {},
	IntradayActivityFieldPoolLap:   {},
	IntradayActivityFieldDuration:  {},
	IntradayActivityFieldHeartRate: {},
	IntradayActivityFieldSpO2Auto:  {},
}

// IsValid checks if v is a valid IntradayActivityField.
func (v IntradayActivityField) IsValid() bool {
	_, ok := validIntradayActivityFieldValues[v]

	return ok
}

// AllIntradayActivityFields is the list of all IntradayActivityField values.
func AllIntradayActivityFields() []IntradayActivityField {
	return []IntradayActivityField{
		IntradayActivityFieldSteps,
		IntradayActivityFieldElevation,
		IntradayActivityFieldCalories,
		IntradayActivityFieldDistance,
		IntradayActivityFieldStroke,
		IntradayActivityFieldPoolLap,
		IntradayActivityFieldDuration,
		IntradayActivityFieldHeartRate,
		IntradayActivityFieldSpO2Auto,
	}
}

type getintradayactivityResponse struct {
	Body IntradayActivities `json:"body"`
}

// IntradayActivities is the response from the Getintradayactivity API call.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getintradayactivity
type IntradayActivities struct {
	Series map[string]IntradayActivity `json:"series"`
}

// IntradayActivity aggregates metrics of a single activity.
//
// Fields are populated based on the requested fields.
//
// TODO: consider making fields pointers, so they dont get populated when no data is returned.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getintradayactivity
type IntradayActivity struct {
	DeviceID string `json:"deviceid"`
	Model    string `json:"model"`
	ModelID  int    `json:"model_id"`

	// Fields
	Steps     int     `json:"steps"`
	Elevation float64 `json:"elevation"` // Note: spec says int, but it's in fact a float
	Calories  float64 `json:"calories"`  // Note: spec says int, but it's in fact a float
	Distance  float64 `json:"distance"`  // Note: spec says int, but it's in fact a float
	Stroke    int     `json:"stroke"`
	PoolLap   int     `json:"pool_lap"`
	Duration  int     `json:"duration"`
	HeartRate int     `json:"heart_rate"`
	SpO2      int     `json:"spo2_auto"`
}

// Getintradayactivity provides activity data for the user with a fine granularity.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getintradayactivity
func (s *MeasureService) Getintradayactivity(ctx context.Context, fields []IntradayActivityField, opts MeasureGetOptions) (*IntradayActivities, *Response, error) {
	fields = filterValidIntradayActivityFieldValues(fields)

	if len(fields) == 0 {
		return nil, nil, errors.New("need at least one intraday activity data field")
	}

	const urlPath = "v2/measure"

	form := url.Values{
		"action":      {"getintradayactivity"},
		"data_fields": {joinIntradayActivityFields(fields)},
	}

	if !opts.StartDate.IsZero() && !opts.EndDate.IsZero() {
		form.Add("startdate", fmt.Sprintf("%d", opts.StartDate.Unix()))
		form.Add("enddate", fmt.Sprintf("%d", opts.EndDate.Unix()))
	}

	intradayactivityResp := new(getintradayactivityResponse)

	resp, err := s.client.PostForm(ctx, urlPath, form, intradayactivityResp)

	return &intradayactivityResp.Body, resp, err
}

func filterValidIntradayActivityFieldValues(values []IntradayActivityField) []IntradayActivityField {
	var validValues []IntradayActivityField

	for _, v := range values {
		if !v.IsValid() {
			continue
		}

		validValues = append(validValues, v)
	}

	return validValues
}

func joinIntradayActivityFields(fields []IntradayActivityField) string {
	s := make([]string, 0, len(fields))

	for _, f := range fields {
		s = append(s, string(f))
	}

	return strings.Join(s, ",")
}

// WorkoutField is a type of metric tracked during workout sessions.
//
// Withings API docs: https://developer.withings.com/api-reference/#operation/measurev2-getworkouts
type WorkoutField string

// WorkoutField values
const (
	WorkoutFieldCalories          WorkoutField = "calories"            // Active calories burned (in Kcal).
	WorkoutFieldIntensity         WorkoutField = "intensity"           // Intensity.
	WorkoutFieldManualDistance    WorkoutField = "manual_distance"     // Distance travelled manually entered by user (in meters).
	WorkoutFieldManualCalories    WorkoutField = "manual_calories"     // Active calories burned manually entered by user (in Kcal).
	WorkoutFieldHRAverage         WorkoutField = "hr_average"          // Average heart rate.
	WorkoutFieldHRMin             WorkoutField = "hr_min"              // Minimal heart rate.
	WorkoutFieldHRMax             WorkoutField = "hr_max"              // Maximal heart rate.
	WorkoutFieldHRZone0           WorkoutField = "hr_zone_0"           // Duration in seconds when heart rate was in a light zone.
	WorkoutFieldHRZone1           WorkoutField = "hr_zone_1"           // Duration in seconds when heart rate was in a moderate zone.
	WorkoutFieldHRZone2           WorkoutField = "hr_zone_2"           // Duration in seconds when heart rate was in an intense zone.
	WorkoutFieldHRZone3           WorkoutField = "hr_zone_3"           // Duration in seconds when heart rate was in maximal zone.
	WorkoutFieldPauseDuration     WorkoutField = "pause_duration"      // Total pause time in second filled by user.
	WorkoutFieldAlgoPauseDuration WorkoutField = "algo_pause_duration" // Total pause time in seconds detected by Withings device (swim only).
	WorkoutFieldSpO2Average       WorkoutField = "spo2_average"        // Average percent of SpO2 percent value during a workout.
	WorkoutFieldSteps             WorkoutField = "steps"               // Number of steps.
	WorkoutFieldDistance          WorkoutField = "distance"            // Distance travelled (in meters).
	WorkoutFieldElevation         WorkoutField = "elevation"           // Number of floors climbed.
	WorkoutFieldPoolLaps          WorkoutField = "pool_laps"           // Number of pool laps.
	WorkoutFieldStrokes           WorkoutField = "strokes"             // Number of strokes.
	WorkoutFieldPoolLength        WorkoutField = "pool_length"         // Length of the pool.
)

var validWorkoutFieldValues = map[WorkoutField]struct{}{
	WorkoutFieldCalories:          {},
	WorkoutFieldIntensity:         {},
	WorkoutFieldManualDistance:    {},
	WorkoutFieldManualCalories:    {},
	WorkoutFieldHRAverage:         {},
	WorkoutFieldHRMin:             {},
	WorkoutFieldHRMax:             {},
	WorkoutFieldHRZone0:           {},
	WorkoutFieldHRZone1:           {},
	WorkoutFieldHRZone2:           {},
	WorkoutFieldHRZone3:           {},
	WorkoutFieldPauseDuration:     {},
	WorkoutFieldAlgoPauseDuration: {},
	WorkoutFieldSpO2Average:       {},
	WorkoutFieldSteps:             {},
	WorkoutFieldDistance:          {},
	WorkoutFieldElevation:         {},
	WorkoutFieldPoolLaps:          {},
	WorkoutFieldStrokes:           {},
	WorkoutFieldPoolLength:        {},
}

// IsValid checks if v is a valid WorkoutField.
func (v WorkoutField) IsValid() bool {
	_, ok := validWorkoutFieldValues[v]

	return ok
}

// AllWorkoutFields is the list of all supported workout fields.
func AllWorkoutFields() []WorkoutField {
	return []WorkoutField{
		WorkoutFieldCalories,
		WorkoutFieldIntensity,
		WorkoutFieldManualDistance,
		WorkoutFieldManualCalories,
		WorkoutFieldHRAverage,
		WorkoutFieldHRMin,
		WorkoutFieldHRMax,
		WorkoutFieldHRZone0,
		WorkoutFieldHRZone1,
		WorkoutFieldHRZone2,
		WorkoutFieldHRZone3,
		WorkoutFieldPauseDuration,
		WorkoutFieldAlgoPauseDuration,
		WorkoutFieldSpO2Average,
		WorkoutFieldSteps,
		WorkoutFieldDistance,
		WorkoutFieldElevation,
		WorkoutFieldPoolLaps,
		WorkoutFieldStrokes,
		WorkoutFieldPoolLength,
	}
}

type getworkoutsResponse struct {
	Body Workouts `json:"body"`
}

// Workouts is the response from the Getworkouts API call.
//
// Withings API docs: https://developer.withings.com/api-reference/#operation/measurev2-getworkouts
type Workouts struct {
	Series []Workout `json:"series"`
}

// Workout aggregates data related to workout sessions from different trackers.
//
// Fields are populated based on the requested fields.
//
// TODO: consider making fields pointers, so they dont get populated when no data is returned.
//
// Withings API docs: https://developer.withings.com/api-reference/#operation/measurev2-getworkouts
type Workout struct {
	Category  int    `json:"category"`
	Timezone  string `json:"timezone"`
	Model     int    `json:"model"`
	Attrib    int    `json:"attrib"`
	Startdate int64  `json:"startdate"`
	Enddate   int64  `json:"enddate"`
	Date      string `json:"date"`
	Modified  int64  `json:"modified"`
	DeviceID  string `json:"deviceid"`

	Data WorkoutData `json:"data"`
}

type WorkoutData struct {
	Calories          float64 `json:"calories"` // Note: spec says int, but it's in fact a float
	Intensity         int     `json:"intensity"`
	ManualDistance    int     `json:"manual_distance"`
	ManualCalories    int     `json:"manual_calories"`
	HrAverage         int     `json:"hr_average"`
	HrMin             int     `json:"hr_min"`
	HrMax             int     `json:"hr_max"`
	HrZone0           int     `json:"hr_zone_0"`
	HrZone1           int     `json:"hr_zone_1"`
	HrZone2           int     `json:"hr_zone_2"`
	HrZone3           int     `json:"hr_zone_3"`
	PauseDuration     int     `json:"pause_duration"`
	AlgoPauseDuration int     `json:"algo_pause_duration"`
	SpO2Average       int     `json:"spo2_average"`
	Steps             int     `json:"steps"`
	Distance          float64 `json:"distance"`  // Note: spec says int, but it's in fact a float
	Elevation         float64 `json:"elevation"` // Note: spec says int, but it's in fact a float
	PoolLaps          int     `json:"pool_laps"`
	Strokes           int     `json:"strokes"`
	PoolLength        int     `json:"pool_length"`
}

// Getworkouts provides data relevant to workout sessions from the different trackers.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getworkouts
func (s *MeasureService) Getworkouts(ctx context.Context, fields []WorkoutField, opts MeasureGetOptions) (*Workouts, *Response, error) {
	fields = filterValidWorkoutFieldValues(fields)

	if len(fields) == 0 {
		return nil, nil, errors.New("need at least one workout data field")
	}

	const urlPath = "v2/measure"

	form := url.Values{
		"action":      {"getworkouts"},
		"data_fields": {joinWorkoutFields(fields)},
	}

	if !opts.LastUpdate.IsZero() {
		form.Add("lastupdate", fmt.Sprintf("%d", opts.LastUpdate.Unix()))
	} else if !opts.StartDate.IsZero() && !opts.EndDate.IsZero() {
		form.Add("startdateymd", opts.StartDate.Format("2006-01-02"))
		form.Add("enddateymd", opts.EndDate.Format("2006-01-02"))
	}

	if opts.Offset > 0 {
		form.Add("offset", fmt.Sprintf("%d", opts.Offset))
	}

	getworkoutsResp := new(getworkoutsResponse)

	resp, err := s.client.PostForm(ctx, urlPath, form, getworkoutsResp)

	return &getworkoutsResp.Body, resp, err
}

func filterValidWorkoutFieldValues(values []WorkoutField) []WorkoutField {
	var validValues []WorkoutField

	for _, v := range values {
		if !v.IsValid() {
			continue
		}

		validValues = append(validValues, v)
	}

	return validValues
}

func joinWorkoutFields(fields []WorkoutField) string {
	s := make([]string, 0, len(fields))

	for _, f := range fields {
		s = append(s, string(f))
	}

	return strings.Join(s, ",")
}
