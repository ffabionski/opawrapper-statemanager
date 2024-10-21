package examplerego

import rego.v1

state := {"counter": data.counter - 1}

default allow := false

allow if {
	data.counter > 0
}
