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

	data.FillTimes()
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

func indicoDayToFahrplanDays(contribs [][]indicoinput.Contribution) (retval []fahrplanoutput.Day) {
	retval = make([]fahrplanoutput.Day, 0, len(contribs))
	for _, cs := range contribs {
		retval = append(retval, ConvertSingleDay(cs))
	}
	return
}

func FahrplanFromApi(data indicoinput.Apiresult) (retval fahrplanoutput.Schedule, err error) {
	data.FillTimes()

	idays, err := IndicoToDays(data)
	if err != nil {
		return
	}

	fdays := indicoDayToFahrplanDays(idays)

	retval.Days = fdays
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

	return
}

func contributionsToRooms_singleDay(contribs []indicoinput.Contribution) []fahrplanoutput.Room {
	rooms := make(map[string][]fahrplanoutput.Event)
	for _, c := range contribs {
		if _, found := rooms[c.Location]; !found {
			rooms[c.Location] = make([]fahrplanoutput.Event, 0)
		}
		rooms[c.Location] = append(rooms[c.Location], ConvertSingleContribution(c))
	}

	retval := make([]fahrplanoutput.Room, 0, len(rooms))
	for k, v := range rooms {
		retval = append(retval, fahrplanoutput.NewRoom(k, v))
	}

	return retval
}

func ConvertSingleDay(contribs []indicoinput.Contribution) (retval fahrplanoutput.Day) {
	retval.Rooms = contributionsToRooms_singleDay(contribs)

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
