// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package features

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/yamlfmt/pkg/yaml"
)

type BlockScalarStyle string

const (
	LiteralBlockStyle BlockScalarStyle = "literal"
	FoldedBlockStyle  BlockScalarStyle = "folded"
)

var ErrUnrecognizedBlockScalarStyle = errors.New("unrecognized block scalar style")

// FeatureForceBlockScalar returns a YAMLFeatureFunc that forces multiline strings
// to use block scalar style (literal "|" or folded ">").
func FeatureForceBlockScalar(style BlockScalarStyle) (YAMLFeatureFunc, error) {
	var toStyle yaml.Style
	switch style {
	case LiteralBlockStyle:
		toStyle = yaml.LiteralStyle
	case FoldedBlockStyle:
		toStyle = yaml.FoldedStyle
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnrecognizedBlockScalarStyle, style)
	}

	var forceStyle YAMLFeatureFunc
	forceStyle = func(n yaml.Node) error {
		for _, c := range n.Content {
			// Only process scalar nodes that are strings with newlines
			if c.Kind == yaml.ScalarNode && c.Tag == "!!str" && strings.Contains(c.Value, "\n") {
				c.Style = toStyle
			}
			if err := forceStyle(*c); err != nil {
				return err
			}
		}
		return nil
	}
	return forceStyle, nil
}
