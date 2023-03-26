// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved

package events

import "errors"

type responseIsNotJSON struct{}

func (n responseIsNotJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("not json")
}
