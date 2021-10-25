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

// Measure types
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

// AllMeasureTypes is the list of all supported measure types.
//
// TODO: make this a method instead?
var AllMeasureTypes = []MeasureType{
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

// Category differentiates between real measurements and user objectives.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type Category int

// Categories
const (
	CategoryReal          Category = 1 // Real measures
	CategoryUserObjective Category = 2 // User objectives
)

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
	GroupID   int64     `json:"grpid"`
	Attrib    int       `json:"attrib"`
	Date      int       `json:"date"`
	CreatedAt int       `json:"created"`
	Category  Category  `json:"category"`
	DeviceID  string    `json:"deviceid"`
	Measures  []Measure `json:"measures"`
	Comment   string    `json:"comment"` // Deprecated
}

// Measure is an individual data point.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
type Measure struct {
	Value int `json:"value"`
	Type  int `json:"type"`
	Unit  int `json:"unit"`
	Algo  int `json:"algo"` // Deprecated
	FM    int `json:"fm"`   // Deprecated
	FW    int `json:"fw"`   // Deprecated
}

// Getmeas provides measures stored on a specific date.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measure-getmeas
func (s *MeasureService) Getmeas(ctx context.Context, measureTypes []MeasureType, category Category, opts MeasureGetOptions) (*Measures, *Response, error) {
	if len(measureTypes) == 0 {
		return nil, nil, errors.New("need at least one measure type")
	}

	const u = "measure"

	form := url.Values{
		"action":   {"getmeas"},
		"category": {fmt.Sprintf("%d", category)},
	}

	if len(measureTypes) == 1 {
		form.Add("meastype", fmt.Sprintf("%d", measureTypes[0]))
	} else {
		form.Add("meastypes", joinMeasureTypes(measureTypes))
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

	resp, err := s.client.PostForm(ctx, u, form, measuresResp)

	return &measuresResp.Body, resp, err
}

func joinMeasureTypes(ints []MeasureType) string {
	s := make([]string, 0, len(ints))

	for _, i := range ints {
		s = append(s, fmt.Sprintf("%d", i))
	}

	return strings.Join(s, ",")
}

// ActivityField is a type of metric tracked during an activity.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getactivity
type ActivityField string

// Activity fields
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

// AllActivityFields is the list of all supported activity fields.
//
// TODO: make this a method instead?
var AllActivityFields = []ActivityField{
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
	if len(fields) == 0 {
		return nil, nil, errors.New("need at least one activity data field")
	}

	const u = "v2/measure"

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

	resp, err := s.client.PostForm(ctx, u, form, activityResp)

	return &activityResp.Body, resp, err
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

// Activity fields
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

// AllIntradayActivityFields is the list of all supported intraday activity fields.
//
// TODO: make this a method instead?
var AllIntradayActivityFields = []IntradayActivityField{
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
	if len(fields) == 0 {
		return nil, nil, errors.New("need at least one intraday activity data field")
	}

	const u = "v2/measure"

	form := url.Values{
		"action":      {"getintradayactivity"},
		"data_fields": {joinIntradayActivityFields(fields)},
	}

	if !opts.StartDate.IsZero() && !opts.EndDate.IsZero() {
		form.Add("startdate", fmt.Sprintf("%d", opts.StartDate.Unix()))
		form.Add("enddate", fmt.Sprintf("%d", opts.EndDate.Unix()))
	}

	intradayactivityResp := new(getintradayactivityResponse)

	resp, err := s.client.PostForm(ctx, u, form, intradayactivityResp)

	return &intradayactivityResp.Body, resp, err
}

func joinIntradayActivityFields(fields []IntradayActivityField) string {
	s := make([]string, 0, len(fields))

	for _, f := range fields {
		s = append(s, string(f))
	}

	return strings.Join(s, ",")
}

// Getworkouts provides data relevant to workout sessions from the different trackers.
//
// Withings API docs: https://developer.withings.com/api-reference#operation/measurev2-getworkouts
func (s *MeasureService) Getworkouts(ctx context.Context) (*Activities, *Response, error) {
	panic("not implemented")
}
