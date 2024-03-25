package usecase

import (
	"context"
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/skyrocketOoO/zanazibar-dag/domain"
)

type VisualUsecase struct {
}

func NewVisualUsecase() *VisualUsecase {
	return &VisualUsecase{}
}

func (u *VisualUsecase) SeeTree(c context.Context, node domain.Node, maxDepth int) (string, error) {
	graph := genGraph(genTree())
	address := "internal/usecase/html/tree.html"
	f, err := os.Create(address)
	if err != nil {
		panic(err)

	}
	graph.Render(io.MultiWriter(f))
	return address, nil
}

func genTree() []opts.TreeData {
	var TreeNodes = []*opts.TreeData{
		{
			Name: "Node33332",
			Children: []*opts.TreeData{
				{
					Name: "Chield1",
				},
			},
		},
		{
			Name: "Node2",
			Children: []*opts.TreeData{
				{
					Name: "Chield1",
				},
				{
					Name: "Chield2",
				},
				{
					Name: "Chield3",
				},
			},
		},
		{
			Name: "Node3",
			// Collapsed: opts.Bool(true),
			Children: []*opts.TreeData{
				{
					Name: "Chield1",
				},
				{
					Name: "Chield2",
				},
				{
					Name: "Chield3",
				},
			},
		},
	}

	var Tree = []opts.TreeData{
		{
			Name:     "Root",
			Children: TreeNodes,
		},
	}
	return Tree
}

func genGraph(treeData []opts.TreeData) *charts.Tree {
	graph := charts.NewTree()
	graph.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "95vh"}),
		charts.WithTitleOpts(opts.Title{Title: "basic tree example"}),
		//charts.WithTooltipOpts(opts.Tooltip{Show: false}),
	)
	graph.AddSeries("tree", treeData).
		SetSeriesOptions(
			charts.WithTreeOpts(
				opts.TreeChart{
					Layout:           "orthogonal",
					Orient:           "LR",
					InitialTreeDepth: -1,
					Leaves: &opts.TreeLeaves{
						Label: &opts.Label{Show: true, Position: "right", Color: "Black"},
					},
				},
			),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "top", Color: "Black"}),
		)
	return graph
}
