package beagle

// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/drone/go-scm/scm"
)

type webhookService struct {
	client *wrapper
}

func (s *webhookService) Parse(req *http.Request, fn scm.SecretFunc) (scm.Webhook, error) {
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}
	var hook scm.Webhook
	err = json.Unmarshal(data, &hook)
	if err != nil {
		return hook, err
	}
	token, err := fn(hook)
	if err != nil {
		return hook, err
	} else if token == "" {
		return hook, nil
	}
	return hook, nil
}
