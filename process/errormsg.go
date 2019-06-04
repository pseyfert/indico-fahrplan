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

	dayvider "github.com/pseyfert/go-dayvider"
	"github.com/pseyfert/indico-fahrplan/indicoinput"
)

func blockinfo(blocks []dayvider.Block, contrib indicoinput.Contribution) (retval string) {
	e := len(blocks) - 1

	retval = retval + fmt.Sprintf("failing contributions: %s\n", contrib.Title)
	retval = retval + fmt.Sprintf("... starts at: %s\n", contrib.StartTime.String())
	retval = retval + fmt.Sprintf("\n")
	retval = retval + fmt.Sprintf("first block:\n")
	retval = retval + fmt.Sprintf("   start: %s\n", blocks[0].Start.String())
	retval = retval + fmt.Sprintf("   end:   %s\n", blocks[0].End.String())
	retval = retval + fmt.Sprintf(" first/last contrib:\n")
	retval = retval + fmt.Sprintf("   start:     %s\n", (*blocks[0].Event).Bookings[blocks[0].Seed].Start.String())
	retval = retval + fmt.Sprintf("   maybe-end: %s\n", (*blocks[0].Event).Bookings[blocks[0].Last-1].End.String())

	retval = retval + fmt.Sprintf("last block:\n")
	retval = retval + fmt.Sprintf("   start: %s\n", blocks[e].Start.String())
	retval = retval + fmt.Sprintf("   end:   %s\n", blocks[e].End.String())
	retval = retval + fmt.Sprintf(" first/last contrib:\n")
	retval = retval + fmt.Sprintf("   start:     %s\n", (*blocks[e].Event).Bookings[blocks[e].Seed].Start.String())
	retval = retval + fmt.Sprintf("   maybe-end: %s\n", (*blocks[e].Event).Bookings[blocks[e].Last-1].End.String())
	return
}
