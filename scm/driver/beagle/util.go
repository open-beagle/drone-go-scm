// // Copyright 2017 Drone.IO Inc. All rights reserved.
// // Use of this source code is governed by a BSD-style
// // license that can be found in the LICENSE file.

package beagle

import (
	"net/url"
	"strconv"

	"github.com/drone/go-scm/scm"
)

func encodeMemberListOptions(opts scm.ListOptions) string {
	params := url.Values{}
	// params.Set("membership", "true")
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.Size != 0 {
		params.Set("per_page", strconv.Itoa(opts.Size))
	}
	if len(opts.URL) > 0 {
		params.Set("scm", opts.URL)
	}
	return params.Encode()
}
