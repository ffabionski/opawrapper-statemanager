package examplerego

# B comunica con C, fino a quando A comunica con B

import rego.v1

default allow := false

allow if {
	input.source == "a"
	input.dest == "b"
}

allow if {
	input.source == "b"
	input.dest == "c"
	data.ab == false
}

state["ab"] if {
	input.source == "a"
	input.dest == "b"
}
