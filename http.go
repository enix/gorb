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

package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kobolog/gorb/core"
	"github.com/kobolog/gorb/network"
	"github.com/kobolog/gorb/util"

	"github.com/gorilla/mux"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.Write(util.MustMarshal(obj, util.JSONOptions{Indent: true}))
}

func writeError(w http.ResponseWriter, err error) {
	var code int

	switch err {
	case core.ErrIpvsSyscallFailed:
		code = http.StatusInternalServerError
	case core.ErrObjectExists:
		code = http.StatusConflict
	case core.ErrObjectNotFound:
		code = http.StatusNotFound
	default:
		code = http.StatusBadRequest
	}

	writeErrorWithCode(w, err, code)
}

func writeErrorWithCode(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(util.MustMarshal(&errorResponse{err.Error()}, util.JSONOptions{Indent: true}))
}

type serviceCreateHandler struct {
	ctx *core.Context
}

func (h serviceCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		opts core.ServiceOptions
		vars = mux.Vars(r)
	)

	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		writeError(w, err)
	} else if err := h.ctx.CreateService(vars["vsID"], &opts); err != nil {
		writeError(w, err)
	}
}

type backendCreateHandler struct {
	ctx *core.Context
}

func (h backendCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		opts core.BackendOptions
		vars = mux.Vars(r)
	)

	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		writeError(w, err)
	} else if err := h.ctx.CreateBackend(vars["vsID"], vars["rsID"], &opts); err != nil {
		writeError(w, err)
	}
}

type backendUpdateHandler struct {
	ctx *core.Context
}

func (h backendUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		opts core.BackendOptions
		vars = mux.Vars(r)
	)

	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		writeError(w, err)
	} else if _, err := h.ctx.UpdateBackend(vars["vsID"], vars["rsID"], opts.Weight); err != nil {
		writeError(w, err)
	}
}

type serviceRemoveHandler struct {
	ctx *core.Context
}

func (h serviceRemoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if _, err := h.ctx.RemoveService(vars["vsID"]); err != nil {
		writeError(w, err)
	}
}

type backendRemoveHandler struct {
	ctx *core.Context
}

func (h backendRemoveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if _, err := h.ctx.RemoveBackend(vars["vsID"], vars["rsID"]); err != nil {
		writeError(w, err)
	}
}

type serviceListHandler struct {
	ctx *core.Context
}

func (h serviceListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if list, err := h.ctx.ListServices(); err != nil {
		writeError(w, err)
	} else {
		writeJSON(w, list)
	}
}

type serviceStatusHandler struct {
	ctx *core.Context
}

func (h serviceStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if opts, err := h.ctx.GetService(vars["vsID"]); err != nil {
		writeError(w, err)
	} else {
		writeJSON(w, opts)
	}
}

type backendStatusHandler struct {
	ctx *core.Context
}

func (h backendStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if opts, err := h.ctx.GetBackend(vars["vsID"], vars["rsID"]); err != nil {
		writeError(w, err)
	} else {
		writeJSON(w, opts)
	}
}

func interfacesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	interfaceName := vars["interface"]
	var err error
	code := http.StatusBadRequest
	var result interface{}

	switch r.Method {
	case "GET":
		result, err = network.GetInterface(interfaceName)
		if err != nil && strings.Contains(err.Error(), "does not exist") {
			code = http.StatusNotFound
		}
	case "PUT":
		opts := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			writeError(w, err)
			return
		}

		result, err = network.AddInterface(interfaceName, opts["type"])
		if err != nil && strings.Contains(err.Error(), "File exists") {
			code = http.StatusConflict
		}
	case "DELETE":
		err = network.DeleteInterface(interfaceName)
		if err != nil && strings.Contains(err.Error(), "Cannot find") {
			code = http.StatusNotFound
		}
	}

	if err != nil {
		writeErrorWithCode(w, err, code)
	} else {
		writeJSON(w, result)
	}
}
