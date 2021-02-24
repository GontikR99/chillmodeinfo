// +build wasm,web

package vuguutil

import (
	"fmt"
	"github.com/GontikR99/chillmodeinfo/pkg/dom"
	"github.com/GontikR99/chillmodeinfo/pkg/dom/document"
)

// Instead of marking elements by id, we create a faux attribute 'node-id' to take it's place.
func GetElementByNodeId(nodeId string) dom.Element {
	return document.QuerySelector(fmt.Sprintf("[node-id='%s']", nodeId))
}