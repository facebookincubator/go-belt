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

package types

// Option is an abstract option which changes the behavior of a Logger.
type Option interface {
	apply(*Config)
}

// Options is a collection of Option-s.
type Options []Option

// Config returns the resulting configuration of a collection of Options.
func (s Options) Config() Config {
	cfg := Config{
		GetCallerFunc: DefaultGetCallerFunc,
	}
	for _, opt := range s {
		opt.apply(&cfg)
	}
	return cfg
}

// Config is a resulting configuration of a collection of Options.
type Config struct {
	// GetCallerPC defines the line of code which invoked the logging entry,
	// see the description of "GetCallerPC" .
	GetCallerFunc GetCallerPC

	// ImplementationSpecificOptions is a set of Logger-implementation-specific
	// options.
	ImplementationSpecificOptions []any
}

// OptionGetCallerFunc overrides GetCallerPC (see description of "GetCallerPC").
type OptionGetCallerFunc GetCallerPC

func (opt OptionGetCallerFunc) apply(cfg *Config) {
	cfg.GetCallerFunc = GetCallerPC(opt)
}
