/*
   Copyright (c) 2015 Andrey Sibiryov <me@kobology.ru>
   Copyright (c) 2015 Other contributors as noted in the AUTHORS file.

   This file is part of GORB - Go Routing and Balancing.

   GORB is free software; you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation; either version 3 of the License, or
   (at your option) any later version.

   GORB is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package pulse

// Pulse is an interface for an active health check for a backend.
type Pulse interface {
	Loop(ID, chan Status)
	Stop()
	Info() Metrics
}

// New creates a new Pulse from the provided endpoint and options.
func New(address string, port uint16, opts *Options) Pulse {
	switch opts.Type {
	case "tcp":
		return NewTCPPulse(address, port, opts)
	case "http":
		return NewHTTPPulse(address, port, opts)
	default:
		return nil
	}
}
