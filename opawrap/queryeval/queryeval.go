package queryeval

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

// TODO: abstract from the policy engine used to evaluate policies (REST calls???)

/*
 * It uses OPA's Go library to evaluate queries...
 * Extract the new state and the real output of the evaluated query
 *
 * data - the data that OPA need to evalute the query
 * input - input provided by the user needed to eveluate the query
 * w - to handle http errors
 * ctx - context of the http request
 */
func OPA(data map[string]any, input any,
	w http.ResponseWriter, ctx context.Context) (map[string]any, map[string]any) {

	packageRego := "data.examplerego"

	store := inmem.NewFromObject(data)

	tx := storage.NewTransactionOrDie(ctx, store)

	re := rego.New(
		rego.Query(packageRego),
		rego.Load([]string{os.Args[1]}, nil),
		rego.Store(store),
		rego.Transaction(tx),
	)

	// Create a prepared query that can be evaluated
	query, err := re.PrepareForEval(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w,
			"OPA couldn't prepare to evalute the query",
			http.StatusInternalServerError)
	}

	// Execute the prepared query
	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Println(err)
		http.Error(w,
			"OPA couldn't evaluate the query",
			http.StatusInternalServerError)
	}

	// Manipulate the result to divide the state from the output
	resultRaw := rs[0].Expressions[0].Value

	result, ok := resultRaw.(map[string]any)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}

	stateRaw, ok := result["state"]
	if ok {
		delete(result, "state")
	}

	state := stateRaw.(map[string]any)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return state, result
}
