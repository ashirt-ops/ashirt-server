// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

/*
	remux is a small module of re-wrapped functions. These exists to provide a more ergonomic
	interface into using chi / modules that utilize google mux implicitly (e.g. dissectors)

	Larger wrappers that act more as decorators are noted as such.
*/

package remux
