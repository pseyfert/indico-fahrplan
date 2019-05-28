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

package fahrplanoutput

import (
	"fmt"
	"time"

	"github.com/pseyfert/indico-fahrplan/indicoinput"
)

// https://stackoverflow.com/a/47342272
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

func NewEvent(c indicoinput.Contribution) (retval Event) {
	retval.Id = c.Id
	retval.Date = c.StartTime
	retval.Duration = fmtDuration(time.Minute * time.Duration(c.Duration))
	retval.Title = c.Title
	retval.Abstract = c.Description

	return
}

func NewRoom(location string, events []Event) (retval Room) {
	retval.Name = location
	retval.Events = make([]Event, 0, len(events))
	for _, e := range events {
		retval.Events = append(retval.Events, e)
	}

	return
}

func contributionsToRooms(contribs []indicoinput.Contribution) []Room {
	var rooms map[string][]Event
	for _, c := range contribs {
		if _, found := rooms[c.Location]; !found {
			rooms[c.Location] = make([]Event, 0)
		}
		rooms[c.Location] = append(rooms[c.Location], NewEvent(c))
	}

	retval := make([]Room, 0, len(rooms))
	for k, v := range rooms {
		retval = append(retval, NewRoom(k, v))
	}

	return retval
}

func NewDay(contribs []indicoinput.Contribution) (retval Day) {
	retval.Rooms = contributionsToRooms(contribs)

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
