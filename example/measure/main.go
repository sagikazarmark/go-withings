/*
Package measure demonstrates the use of the Measure service.

Look at the "auth" example to get an access token.

Set the following environment variable:
	export WITHINGS_ACCESS_TOKEN="<YOUR ACCESS TOKEN>"

Then run the application:
	go run main.go
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"

	"github.com/sagikazarmark/go-withings/withings"
)

func main() {
	client := withings.NewClient(
		oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: os.Getenv("WITHINGS_ACCESS_TOKEN"),
		})),
	)

	now := time.Now()

	opts := withings.MeasureGetOptions{
		LastUpdate: now.Add(-24 * time.Hour),
		StartDate:  now.Add(-24 * time.Hour),
		EndDate:    now,
	}

	measures, _, err := client.Measure.Getmeas(
		context.Background(),
		withings.AllMeasureTypes(),
		withings.MeasureCategoryRealMeasure,
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	activities, _, err := client.Measure.Getactivity(
		context.Background(),
		withings.AllActivityFields(),
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	intradayactivities, _, err := client.Measure.Getintradayactivity(
		context.Background(),
		withings.AllIntradayActivityFields(),
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	workouts, _, err := client.Measure.Getworkouts(
		context.Background(),
		withings.AllWorkoutFields(),
		opts,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Measures: %#v\n\n", measures)
	fmt.Printf("Activities: %#v\n\n", activities)
	fmt.Printf("Intraday activities: %#v\n\n", intradayactivities)
	fmt.Printf("Workouts: %#v\n\n", workouts)
}
