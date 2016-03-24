/*
Strava Commute Uploader

(c) Douglas Watson, 2016

This application connects to the Strava API and uploads a commute for every day
in a specified time range

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/strava/go.strava"
)

func main() {

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("USAGE:\n\tcommutes config.yaml")
		os.Exit(1)
	}

	config, err := ConfigFromYaml(flag.Arg(0))
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to Strava API
	client := strava.NewClient(config.AccessToken)
	service := strava.NewActivitiesService(client)

	// Post commute every day except weekends or dates in exclusion periods
	i := 0 // index of current/next exclusion period
	for date := config.StartDate; date.Before(config.EndDate); date = date.AddDate(0, 0, 1) {

		// Is date in exclusion period?
		// This method isn't very general but works for my purpose.
		if date.After(config.ExcludeDates[i][0]) && date.Before(config.ExcludeDates[i][1]) {
			continue
		}
		if date.After(config.ExcludeDates[i][1]) && i < len(config.ExcludeDates)-1 {
			i++
		}

		// Optionally skip weekends
		if config.OnlyWeekdays &&
			(date.Weekday() == time.Saturday || date.Weekday() == time.Sunday) {
			continue
		}
		log.Println("Post commute:", date, ":", date.Weekday().String())

		err := Upload(service, date, config.Duration*60, config.Distance*1000.0, config.GearID)
		if err != nil {
			log.Fatal("Error uploading", err)
		}
	}
	log.Println("Done.")
}

// Upload will create a new commute.
func Upload(service *strava.ActivitiesService, startTime time.Time, duration int, distance float64, gearID string) error {

	activity, err := service.Create("Commute", strava.ActivityTypes.Ride, startTime, duration).
		Distance(distance).
		Do()

	if err != nil {
		if e, ok := err.(strava.Error); ok && e.Message == "Authorization Error" {
			log.Printf("Make sure your token has 'write' permissions; check README.md for instructions.")
		}

		return err
	}

	// Mark as commute and make private
	activity, err = service.Update(activity.Id).
		Commute(true).
		Private(true).
		Gear(gearID).
		Do()

	return err
}
