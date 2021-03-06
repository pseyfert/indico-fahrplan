/*
 * Copyright (C) 2019 Paul Seyfert
 * Author: Paul Seyfert <pseyfert.mathphys@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package process

import (
	"fmt"
	"time"

	dayvider "github.com/pseyfert/go-dayvider"
	"github.com/pseyfert/indico-fahrplan/fahrplanoutput"
	"github.com/pseyfert/indico-fahrplan/indicoinput"
)

func IndicoToDayvider(data indicoinput.Apiresult) dayvider.Event {
	bookings := make([]dayvider.Booking, 0, len(data.Results.Conference.Contributions.Contributions))

	data.Parse()
	for _, contrib := range data.Results.Conference.Contributions.Contributions {
		if !contrib.TimeLess {
			start := contrib.StartTime
			end := contrib.EndTime
			bookings = append(bookings, dayvider.Booking{Start: start, End: end})
		}
	}
	event := dayvider.NewEvent(bookings)
	return event
}

func indicoDayToFahrplanDays(contribs [][]indicoinput.Contribution, refdate time.Time) (retval []fahrplanoutput.Day) {
	retval = make([]fahrplanoutput.Day, 0, len(contribs))
	for i, cs := range contribs {
		daystarts := refdate.Add(24 * time.Duration(i) * time.Hour)
		retval = append(retval, ConvertSingleDay(cs, daystarts))
	}
	return
}

func FahrplanFromApi(data indicoinput.Apiresult) (retval fahrplanoutput.Schedule, err error) {
	data.Parse()

	idays, refdate, err := IndicoToDays(data)
	if err != nil {
		return
	}

	fdays := indicoDayToFahrplanDays(idays, refdate)

	retval.Days = fdays

	for i := range retval.Days {
		retval.Days[i].Index = i + 1
		retval.Days[i].Start = retval.Days[i].Start.In(data.Results.Conference.TimezoneLocation)
		retval.Days[i].End = retval.Days[i].End.In(data.Results.Conference.TimezoneLocation)
		retval.Days[i].Date = retval.Days[i].Start.Format("2006-01-02")

		for j := range retval.Days[i].Rooms {
			for k := range retval.Days[i].Rooms[j].Events {
				retval.Days[i].Rooms[j].Events[k].Date = retval.Days[i].Rooms[j].Events[k].Date.In(data.Results.Conference.TimezoneLocation)
				retval.Days[i].Rooms[j].Events[k].Start = retval.Days[i].Rooms[j].Events[k].Date.Format("15:04")
			}
		}
	}

	retval.Conference.Days = len(fdays)
	// retval.Conference.Start = idays[0][0].StartTime.Format("2006-01-02")
	// retval.Conference.End = idays[len(idays)-1][len(idays[len(idays)-1])-1].EndTime.Format("2006-01-02")
	retval.Conference.Start = data.Results.Conference.StartTime.In(data.Results.Conference.TimezoneLocation).Format("2006-01-02")
	retval.Conference.End = data.Results.Conference.EndTime.In(data.Results.Conference.TimezoneLocation).Format("2006-01-02")
	retval.Conference.SlotDuration = "00:10"
	retval.Conference.Acronym = "CERN"
	retval.Conference.Title = "IndicoImport"
	retval.Conference.Url = fmt.Sprintf("https://indico.cern.ch/event/%d", data.Results.Conference.Id)
	retval.Version = "DebugDuty"
	return
}

// https://stackoverflow.com/a/47342272
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func ConvertSingleContribution(c indicoinput.Contribution) (retval fahrplanoutput.Event) {
	retval.Id = c.Id
	retval.Date = c.StartTime
	retval.Duration = fmtDuration(time.Minute * time.Duration(c.Duration))
	retval.Title = c.Title
	retval.Abstract = c.Description
	if c.Track != "" && c.Session != c.Track {
		retval.Track = c.Track + "(" + c.Session + ")"
	} else if c.Track != "" {
		retval.Track = c.Track
	} else if c.Session != "" {
		retval.Track = c.Session
	}

	if len(c.Folders.Folders) != 0 {
		retval.Attachments.Attachments = make([]fahrplanoutput.Attachment, 0, 1+len(c.Folders.Folders))
		retval.Attachments.Attachments = append(retval.Attachments.Attachments, fahrplanoutput.Attachment{Href: c.Url, Name: "Contribution"})
		for _, f := range c.Folders.Folders {
			for _, a := range f.Attachments.Attachments {
				var newattachment fahrplanoutput.Attachment
				newattachment.Href = a.Url
				if a.Title != "" && a.Title != a.Filename && a.Filename != "" {
					newattachment.Name = a.Title + " (" + a.Filename + ")"
				} else if a.Title != "" {
					newattachment.Name = a.Title
				} else if a.Filename != "" {
					newattachment.Name = a.Filename
				} else {
					newattachment.Name = "unnamed attachment"
				}
				retval.Attachments.Attachments = append(retval.Attachments.Attachments, newattachment)
			}
		}
	}

	if len(c.Speakers.Contributionparticipation) > 0 {
		retval.Persons.Persons = make([]fahrplanoutput.Person, len(c.Speakers.Contributionparticipation))
		for i, p_in := range c.Speakers.Contributionparticipation {
			retval.Persons.Persons[i].Id = p_in.Id
			retval.Persons.Persons[i].Name = p_in.FullName
		}
	}

	return
}

func contributionsToRooms_singleDay(contribs []indicoinput.Contribution) []fahrplanoutput.Room {
	rooms := make(map[string][]fahrplanoutput.Event)
	for _, c := range contribs {
		if _, found := rooms[c.CombinedLocation]; !found {
			rooms[c.CombinedLocation] = make([]fahrplanoutput.Event, 0)
		}
		rooms[c.CombinedLocation] = append(rooms[c.CombinedLocation], ConvertSingleContribution(c))
	}

	retval := make([]fahrplanoutput.Room, 0, len(rooms))
	for k, v := range rooms {
		retval = append(retval, fahrplanoutput.NewRoom(k, v))
	}

	return retval
}

func ConvertSingleDay(contribs []indicoinput.Contribution, refdate time.Time) (retval fahrplanoutput.Day) {
	retval.Rooms = contributionsToRooms_singleDay(contribs)

	// use dayvider's days if no contribution available
	if len(contribs) < 1 {
		retval.Start = refdate
		retval.End = refdate.Add(24 * time.Hour)
		return
	}
	// sane default, fixable later
	retval.Start = contribs[0].StartTime
	retval.End = contribs[0].EndTime
	for _, c := range contribs {
		if c.StartTime.Before(retval.Start) {
			retval.Start = c.StartTime
		}
		if c.EndTime.After(retval.End) {
			retval.End = c.EndTime
		}
	}

	return
}
