package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/go/sgf"
	"github.com/otrego/clamshell/katago"
	"github.com/otrego/clamshell/katago/kataprob"
)

func tryThing() error {
	fi := "40391047-169-quacimodo-marcmaniez.sgf"

	content, err := ioutil.ReadFile(fi)
	if err != nil {
		return err
	}

	g, err := sgf.FromString(string(content)).Parse()
	if err != nil {
		return err
	}

	analysisBytes, err := ioutil.ReadFile("analysis.json")
	if err != nil {
		return err
	}
	result := &katago.AnalysisList{}
	if err := json.Unmarshal(analysisBytes, result); err != nil {
		return err
	}

	if err := result.AddToGame(g); err != nil {
		return err
	}
	fmt.Print("children \n")
	fmt.Printf("First move info %#v\n", g.Root.Children[0].AnalysisData().(*katago.AnalysisResult).MoveInfos[0])
	paths, err := kataprob.FindBlunders(g)
	if err != nil {
		return err
	}

	spew.Dump(paths)

	blunderNodes := Map(paths, func(p movetree.Path) movetree.Node { return *p.Apply(g.Root)})

	sort.Slice(blunderNodes, func(i, j int) bool {
		iDelta := scoreDelta(blunderNodes[i])
		jDelta := scoreDelta(blunderNodes[j])
		return iDelta > jDelta
	})
	fmt.Printf("Blunder nodes %#v\n", blunderNodes)

	return nil
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
			result[i] = fn(t)
	}
	return result
}

func scoreDelta (n movetree.Node) float64 {
	nScoreLead := n.AnalysisData().(*katago.AnalysisResult).MoveInfos[0].ScoreLead
	pScoreLead := n.Parent.AnalysisData().(*katago.AnalysisResult).MoveInfos[0].ScoreLead
	return nScoreLead - pScoreLead
}