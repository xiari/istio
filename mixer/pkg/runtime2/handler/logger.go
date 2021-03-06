// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"istio.io/istio/mixer/pkg/adapter"
	"istio.io/istio/pkg/log"
)

type logger struct {
	l *zap.Logger
	s *zap.SugaredLogger
}

var _ adapter.Logger = &logger{}

// stackDepth is used
// to determine how many levels of stack to skip
// it should be 1 if logging through a struct and
// it should be 2 if logging through an interface
// If depth=1 is used with an interface, filename appears as <autogenerated>
// logger is used through the adapter.Logger interface.
const stackDepth = 2

func newLogger(name string) logger {
	l := log.With(zap.String("adapter", name)).WithOptions(zap.AddCallerSkip(stackDepth))
	s := l.Sugar()

	return logger{
		l: l,
		s: s,
	}
}

// VerbosityLevel from adapter.Logger.
func (l logger) VerbosityLevel(level adapter.VerbosityLevel) bool {
	switch level {
	case 0:
		return l.l.Core().Enabled(zap.ErrorLevel)
	case 1:
		return l.l.Core().Enabled(zap.WarnLevel)
	case 2:
		return l.l.Core().Enabled(zap.InfoLevel)
	}

	if level >= 3 {
		return l.l.Core().Enabled(zap.DebugLevel)
	}

	// < 0
	return false
}

// Infof from adapter.Logger.
func (l logger) Infof(format string, args ...interface{}) {
	l.s.Infof(format, args...)
}

// Warningf from adapter.Logger.
func (l logger) Warningf(format string, args ...interface{}) {
	l.s.Warnf(format, args...)
}

// Errorf from adapter.Logger.
func (l logger) Errorf(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args...)
	l.s.Error(s)
	return errors.New(s)
}
