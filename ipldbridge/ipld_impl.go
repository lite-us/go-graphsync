package ipldbridge

import (
	"bytes"
	"context"

	"github.com/ipld/go-ipld-prime/fluent"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/encoding/dagcbor"
	free "github.com/ipld/go-ipld-prime/impl/free"
	ipldtraversal "github.com/ipld/go-ipld-prime/traversal"
	ipldselector "github.com/ipld/go-ipld-prime/traversal/selector"
	selectorbuilder "github.com/ipld/go-ipld-prime/traversal/selector/builder"
)

// TraversalConfig is an alias from ipld, in case it's renamed/moved.
type TraversalConfig = ipldtraversal.TraversalConfig

type ipldBridge struct {
}

// NewIPLDBridge returns an IPLD Bridge.
func NewIPLDBridge() IPLDBridge {
	return &ipldBridge{}
}

func (rb *ipldBridge) ExtractData(node ipld.Node, buildFn func(SimpleNode) interface{}) (interface{}, error) {
	var value interface{}
	err := fluent.Recover(func() {
		simpleNode := fluent.WrapNode(node)
		value = buildFn(simpleNode)
	})
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (rb *ipldBridge) BuildNode(buildFn func(NodeBuilder) ipld.Node) (ipld.Node, error) {
	var node ipld.Node
	err := fluent.Recover(func() {
		nb := fluent.WrapNodeBuilder(free.NodeBuilder())
		node = buildFn(nb)
	})
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (rb *ipldBridge) BuildSelector(buildFn func(SelectorSpecBuilder) SelectorSpec) (ipld.Node, error) {
	var node ipld.Node
	err := fluent.Recover(func() {
		ssb := selectorbuilder.NewSelectorSpecBuilder(free.NodeBuilder())
		node = buildFn(ssb).Node()
	})
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (rb *ipldBridge) Traverse(ctx context.Context, loader Loader, root ipld.Link, s Selector, fn AdvVisitFn) error {
	node, err := root.Load(ctx, LinkContext{}, free.NodeBuilder(), loader)
	if err != nil {
		return err
	}
	return TraversalProgress{
		Cfg: &TraversalConfig{
			Ctx:        ctx,
			LinkLoader: loader,
		},
	}.TraverseInformatively(node, s, fn)
}

func (rb *ipldBridge) EncodeNode(node ipld.Node) ([]byte, error) {
	var buffer bytes.Buffer
	err := dagcbor.Encoder(node, &buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (rb *ipldBridge) DecodeNode(encoded []byte) (ipld.Node, error) {
	reader := bytes.NewReader(encoded)
	return dagcbor.Decoder(free.NodeBuilder(), reader)
}

func (rb *ipldBridge) ParseSelector(selector ipld.Node) (Selector, error) {
	return ipldselector.ParseSelector(selector)
}
