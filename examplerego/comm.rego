package examplerego

# b comunica con c, fintanto che a non ha comunicato con b

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
