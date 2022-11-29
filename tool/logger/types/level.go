// Copyright 2022 Meta Platforms, Inc. and affiliates.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Copyright (c) Facebook, Inc. and its affiliates.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Package types of logger unifies different types of loggers into
// interfaces Logger. For example it allows to upgrade simple fmt.Printf
// to be a fully functional Logger. Therefore multiple wrappers are implemented
// here to provide different functions which could be missing in some loggers.
package types

import (
	"fmt"
	"strings"
)

// Level is a severity of a message/Entry.
//
// There are two ways to use Level:
// 1. To define a severity of a specific message/Entry (when logging).
// 2. To define a severity of messages/Entries be actually logged (when configuring a Logger).
type Level int

const (
	// LevelUndefined is an erroneous value of log-level which just corresponds
	// to the zero-value.
	LevelUndefined = Level(iota)

	// LevelNone means to do not log anything.
	//
	// But even if a Logger is setup with this level, messages with levels Panic
	// and Fatal will still cause panics and os.Exit, just we except this to avoid
	// sending log messages (some Logger implementations may ignore this rule).
	LevelNone

	// LevelFatal means non-recoverable error case.
	//
	// If a message sent with this level then os.Exit is invoked in the end of processing the message.
	//
	// If a Logger is setup with this level, it will just ignore messages with level higher than Panic,
	// will silently panic on a Panic message and will loudly exit on a Fatal message.
	//
	// Some Logger implementations may ignore the rule about panicking "silently".
	LevelFatal

	// LevelPanic means a panic case (basically a recoverable, but critical problem).
	//
	// If a message sent with this level then panic() is invoked in the end of processing the message.
	//
	// If a Logger is setup with this level, it will ignore messages with level higher than Panic.
	LevelPanic

	// LevelError means an error case (not a critical problem, but yet an error).
	//
	// A message with this level is just logged if the Logger is setup with level no less that this.
	LevelError

	// LevelWarning means an odd/unexpected/wrong/undesired/etc case. Basically this is something
	// to keep an eye on. For example which could explain odd/erroneous behavior of the application
	// in future.
	//
	// A message with this level is just logged if the Logger is setup with level no less that this.
	LevelWarning

	// LevelInfo means an information message, essential enough to notify the end user (who is
	// not a developer of the application), but about something benign and that does not
	// says to any wrong behavior of the application.
	//
	// A message with this level is just logged if the Logger is setup with level no less that this.
	//
	// Recommended as the default level.
	LevelInfo

	// LevelDebug means a message non-essential for the end user, but useful for debugging
	// the application by its developer.
	//
	// A message with this level is just logged if the Logger is setup with level no less that this.
	LevelDebug

	// LevelTrace is for all other messages.
	//
	// For example, sometimes in complex processes/application it sometimes useful
	// to have an extra layer to put an insane amount put all useless messages, which
	// may suddenly become very useful for some very difficult debug process.
	//
	// A message with this level is just logged if the Logger is setup with level no less that this.
	LevelTrace

	// EndOfLevel just defines a maximum level + 1.
	EndOfLevel
)

// Byte is similar to String, but returns a single byte, which
// describes in a human-readable way the logging level.
func (logLevel Level) Byte() byte {
	switch logLevel {
	case LevelUndefined:
		return '?'
	case LevelNone:
		return '-'
	case LevelDebug:
		return 'D'
	case LevelInfo:
		return 'I'
	case LevelWarning:
		return 'W'
	case LevelError:
		return 'E'
	case LevelPanic:
		return 'P'
	case LevelFatal:
		return 'F'
	}
	return 'U'
}

// String just implements fmt.Stringer, flag.Value and pflag.Value.
func (logLevel Level) String() string {
	switch logLevel {
	case LevelUndefined:
		return "undefined"
	case LevelNone:
		return "none"
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelPanic:
		return "panic"
	case LevelFatal:
		return "fatal"
	}
	return fmt.Sprintf("unknown_%d", logLevel)
}

// Set updates the logging level values based on the passed string value.
// This method just implements flag.Value and pflag.Value.
func (logLevel *Level) Set(value string) error {
	newLogLevel, err := ParseLogLevel(value)
	if err != nil {
		return err
	}
	*logLevel = newLogLevel
	return nil
}

// Type just implements pflag.Value.
func (logLevel *Level) Type() string {
	return "Level"
}

// ParseLogLevel parses incoming string into a Level and returns
// LevelUndefined with an error if an unknown logging level was passed.
func ParseLogLevel(in string) (Level, error) {
	switch strings.ToLower(in) {
	case "t", "trace":
		return LevelTrace, nil
	case "d", "debug":
		return LevelDebug, nil
	case "i", "info":
		return LevelInfo, nil
	case "w", "warn", "warning":
		return LevelWarning, nil
	case "e", "err", "error":
		return LevelError, nil
	case "p", "panic":
		return LevelPanic, nil
	case "f", "fatal":
		return LevelFatal, nil
	case "n", "none":
		return LevelNone, nil
	}
	var allowedValues []string
	for logLevel := LevelFatal; logLevel <= LevelDebug; logLevel++ {
		allowedValues = append(allowedValues, logLevel.String())
	}
	return LevelUndefined, fmt.Errorf("unknown logging level '%s', known values are: %s",
		in, strings.Join(allowedValues, ", "))
}
