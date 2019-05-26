# Copyright (C) 2019 Paul Seyfert
# Author: Paul Seyfert <pseyfert.mathphys@gmail.com>
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

wget 'https://indico.cern.ch/export/timetable/708041.xml' --output-document=indico.timetable.xml
wget 'https://indico.cern.ch/export/event/708041.xml?detail=contributions' --output-document=indico.event.detail.contributions.xml
wget https://fahrplan.events.ccc.de/congress/2018/Fahrplan/schedule.xml --output-document=fahrplan.xml
