// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

// Exercise 4.10: Modify issues to report the results in age categories, say less than a month old, less than a year old, and more than a year old.

const (
	Day   = 24 * time.Hour
	Month = 30 * Day  // duration in nanosecods
	Year  = 365 * Day // duration in nanosecods
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d issues:\n", result.TotalCount)

	groups := make(map[string][]*github.Issue)

	for _, item := range result.Items {
		age := time.Now().Sub(item.CreatedAt)
		if age < Month {
			groups["month"] = append(groups["month"], item)
		} else if age <= Year {
			groups["year"] = append(groups["year"], item)
		} else {
			groups["year+"] = append(groups["year+"], item)
		}
	}

	for k, v := range groups {
		fmt.Printf("%q old issues:\n", k)
		for _, item := range v {
			fmt.Printf("#%-5d %9.9s %.55s\n",
				item.Number, item.User.Login, item.Title)
		}
	}

}
