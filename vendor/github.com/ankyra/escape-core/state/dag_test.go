/*
Copyright 2017 Ankyra

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package state

import (
	. "gopkg.in/check.v1"
)

func (s *suite) Test_GetDeploymentStateDAG_empty_env(c *C) {
	prj, _ := NewProjectState("my-project")
	env := prj.GetEnvironmentStateOrMakeNew("my-env")
	dag, err := env.GetDeploymentStateDAG("build")
	c.Assert(err, IsNil)
	c.Assert(dag, HasLen, 0)
}

func (s *suite) Test_GetDeploymentStateDAG_one_deployment(c *C) {
	stage := "deploy"
	prj, _ := NewProjectState("my-project")
	env := prj.GetEnvironmentStateOrMakeNew("my-env")
	depl1 := env.GetOrCreateDeploymentState("depl1")
	depl1.GetStageOrCreateNew(stage)
	dag, err := env.GetDeploymentStateDAG(stage)
	c.Assert(err, IsNil)
	c.Assert(dag, HasLen, 1)
	c.Assert(dag[0].Node, DeepEquals, depl1)
	c.Assert(dag[0].AndThen, HasLen, 0)
}

func (s *suite) Test_GetDeploymentStateDAG_two_deployments_one_provider(c *C) {
	stage := "deploy"
	prj, _ := NewProjectState("my-project")
	env := prj.GetEnvironmentStateOrMakeNew("my-env")
	depl1 := env.GetOrCreateDeploymentState("depl1")
	depl2 := env.GetOrCreateDeploymentState("depl2")
	st := depl1.GetStageOrCreateNew(stage)
	st.Providers["whatever"] = "depl2"
	depl2.GetStageOrCreateNew(stage)

	dag, err := env.GetDeploymentStateDAG(stage)
	c.Assert(err, IsNil)
	c.Assert(dag, HasLen, 1)
	c.Assert(dag[0].Node, DeepEquals, depl2)
	c.Assert(dag[0].AndThen, HasLen, 1)
	c.Assert(dag[0].AndThen[0].Node, DeepEquals, depl1)
	c.Assert(dag[0].AndThen[0].AndThen, HasLen, 0)

	tsort, err := env.GetDeploymentStateTopologicalSort(stage)
	c.Assert(err, IsNil)
	c.Assert(tsort, HasLen, 2)
	c.Assert(tsort[0], DeepEquals, depl2)
	c.Assert(tsort[1], DeepEquals, depl1)
}

func (s *suite) Test_GetDeploymentStateDAG(c *C) {
	// For deployment graph:
	//
	// A -> B, E
	// B -> C, D
	// C -> D
	// D
	// E

	stage := "deploy"
	prj, _ := NewProjectState("my-project")
	env := prj.GetEnvironmentStateOrMakeNew("my-env")
	deplA := env.GetOrCreateDeploymentState("deplA")
	deplB := env.GetOrCreateDeploymentState("deplB")
	deplC := env.GetOrCreateDeploymentState("deplC")
	deplD := env.GetOrCreateDeploymentState("deplD")
	deplE := env.GetOrCreateDeploymentState("deplE")

	stA := deplA.GetStageOrCreateNew(stage)
	stA.Providers["b"] = "deplB"
	stA.Providers["e"] = "deplE"

	stB := deplB.GetStageOrCreateNew(stage)
	stB.Providers["c"] = "deplC"
	stB.Providers["d"] = "deplD"

	stC := deplC.GetStageOrCreateNew(stage)
	stC.Providers["d"] = "deplD"

	deplD.GetStageOrCreateNew(stage)
	deplE.GetStageOrCreateNew(stage)

	dag, err := env.GetDeploymentStateDAG(stage)
	c.Assert(err, IsNil)
	c.Assert(dag, HasLen, 2)

	var bDag, cDag, dDag, eDag *DAGNode
	if dag[0].Node.Name == "deplD" {
		dDag = dag[0]
		eDag = dag[1]
	} else {
		dDag = dag[1]
		eDag = dag[0]
	}
	c.Assert(dDag.Node, DeepEquals, deplD)
	c.Assert(dDag.AndThen, HasLen, 2)

	if dDag.AndThen[0].Node.Name == "deplB" {
		bDag = dDag.AndThen[0]
		cDag = dDag.AndThen[1]
	} else {
		bDag = dDag.AndThen[1]
		cDag = dDag.AndThen[0]
	}
	c.Assert(bDag.Node, DeepEquals, deplB)
	c.Assert(bDag.AndThen, HasLen, 1)
	c.Assert(bDag.AndThen[0].Node, DeepEquals, deplA)

	c.Assert(cDag.Node, DeepEquals, deplC)
	c.Assert(cDag.AndThen, HasLen, 1)
	c.Assert(cDag.AndThen[0].Node, DeepEquals, deplB)

	c.Assert(eDag.Node, DeepEquals, deplE)
	c.Assert(eDag.AndThen, HasLen, 1)
	c.Assert(eDag.AndThen[0].Node, DeepEquals, deplA)

	i := 0
	for i < 1000 {
		tsort, err := env.GetDeploymentStateTopologicalSort(stage)
		c.Assert(err, IsNil)
		for ix, depl := range tsort {
			st := depl.GetStageOrCreateNew(stage)
			for _, deplName := range st.Providers {
				found := false
				for depIx, depDepl := range tsort {
					if depDepl.Name == deplName {
						found = true
						c.Assert(depIx < ix, Equals, true, Commentf("Deployment '%s' should happen before '%s'", deplName, depl.Name))
					}
				}
				c.Assert(found, Equals, true, Commentf("Missing deployment '%s' in topological sort", depl.Name))
			}
		}
		i += 1
	}
}

type hasItemChecker struct{}

var HasItem = &hasItemChecker{}

func (*hasItemChecker) Info() *CheckerInfo {
	return &CheckerInfo{Name: "HasItem", Params: []string{"obtained", "expected to have item"}}
}
func (*hasItemChecker) Check(params []interface{}, names []string) (bool, string) {
	obtained := params[0]
	expectedItem := params[1]
	switch obtained.(type) {
	case []interface{}:
		for _, v := range obtained.([]interface{}) {
			if v == expectedItem {
				return true, ""
			}
		}
	case []string:
		for _, v := range obtained.([]string) {
			if v == expectedItem {
				return true, ""
			}
		}
	default:
		return false, "Unexpected type."
	}
	return false, "Item not found"
}
