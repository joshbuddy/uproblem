// Package kataprob finds problems from games that have katago-analysis data
// attached. It is assumed that all games have the relevant katago analysis data
// attached before this point.
package kataprob

import (
	"fmt"
	"math"

	"github.com/golang/glog"
	"github.com/otrego/clamshell/go/movetree"
	"github.com/otrego/clamshell/katago"
)

// FindBlunders finds positions (paths) that result from big swings in points.
func FindBlunders(g *movetree.MoveTree) ([]movetree.Path, error) {
	blunderAmt := 3.0

	var cur movetree.Path
	var found []movetree.Path
	if g.Root == nil {
		return found, nil
	}

	var prevLead float64

	for n := g.Root; n != nil; n = n.Next(0) {
		glog.V(3).Infof("VarNum %v\n", n.VarNum())
		glog.V(3).Infof("MoveNum %v\n", n.MoveNum())

		// We assume alternating moves. Lead is always presented as
		pl := prevLead
		cur = append(cur, n.VarNum())
		glog.V(3).Infof("PrevLead %v\n", prevLead)

		d := n.AnalysisData()
		if d == nil {
			glog.Infof("nil analysis data")
			continue
		}

		katad, ok := d.(*katago.AnalysisResult)
		if !ok {
			glog.V(2).Infof("not analysisResult")
			continue
		}
		if katad.RootInfo == nil {
			// This
			glog.Errorf("no RootInfo for at move %v", n.MoveNum())
			continue
		}

		lead := katad.RootInfo.ScoreLead
		var nextLead float64
		//if
		nextLead = -1 * lead
		//glog.V(3).Infof("Next ScoreLead: %f:\n", nextLead)
		delta := nextLead - pl
		fmt.Printf("%d nextLead: %f Delta: %f katad.RootInfo.ScoreLead: %f\n", n.MoveNum(), nextLead, delta, katad.RootInfo.ScoreLead)

		if delta >= math.Abs(blunderAmt) {
			found = append(found, cur.Clone())
		}

		// prevLead is always with respect to current player
		prevLead = nextLead
	}

	return found, nil
}
