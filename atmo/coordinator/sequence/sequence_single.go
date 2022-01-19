package sequence

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/suborbital/atmo/directive/executable"
	"github.com/suborbital/reactr/request"
	"github.com/suborbital/reactr/rt"
)

var ErrMissingFQFN = errors.New("callableFn missing FQFN")

func (seq Sequence) ExecSingleFn(fn executable.CallableFn, req *request.CoordinatedRequest) (*FnResult, error) {
	start := time.Now()
	defer func() {
		seq.log.Debug("fn", fn.Fn, "executed in", time.Since(start).Milliseconds(), "ms")
	}()

	if fn.FQFN == "" {
		return nil, ErrMissingFQFN
	}

	var jobResult []byte
	var runErr rt.RunErr

	// Do will execute the job locally if possible or find a remote peer to execute it
	res, err := seq.exec.Do(fn.FQFN, req, seq.ctx)
	if err != nil {
		// check if the error type is rt.RunErr, because those are handled differently
		if returnedErr, isRunErr := err.(rt.RunErr); isRunErr {
			runErr = returnedErr
		} else {
			return nil, errors.Wrap(err, "failed to exec.Do")
		}
	} else if res != nil {
		jobResult = res.([]byte)
	} else {
		return nil, nil
	}

	// runErr would be an actual error returned from a function
	// should find a better way to determine if a RunErr is "non-nil"
	if runErr.Code != 0 || runErr.Message != "" {
		seq.log.Debug("fn", fn.Fn, "returned an error")
	} else if jobResult == nil {
		seq.log.Debug("fn", fn.Fn, "returned a nil result")
	}

	cResponse := &request.CoordinatedResponse{}

	if jobResult != nil {
		if err := json.Unmarshal(jobResult, cResponse); err != nil {
			// handle backwards-compat
			cResponse.Output = jobResult
		}
	}

	result := &FnResult{
		FQFN:     fn.FQFN,
		Key:      fn.Key(),
		Response: cResponse,
		RunErr:   runErr,
	}

	return result, nil
}
