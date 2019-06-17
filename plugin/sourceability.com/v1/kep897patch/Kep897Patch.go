package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"sigs.k8s.io/kustomize/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	Patches []Patch `json:"patches"`

	ldr ifc.Loader
	rf  *resmap.Factory
}

type Patch struct {
	// Path is a relative file path to the patch file.
	Path string `json:"path,omitempty"`

	// Target points to the resources that the patch is applied to
	Target PatchTarget `json:"target,omitempty"`

	// Type is one of `StrategicMergePatch` or `JsonPatch`
	Type string `json:"type,omitempty"`
}

// PatchTarget specifies a set of resources
type PatchTarget struct {
	// Group of the target
	Group string `json:"group,omitempty"`

	// Version of the target
	Version string `json:"version,omitemtpy"`

	// Kind of the target
	Kind string `json:"kind,omitempty"`

	// Name of the target
	// The name could be with wildcard to match a list of Resources
	Name string `json:"name,omitempty"`

	// MatchAnnotations is a map of key-value pairs.
	// A Resource matches it will be appied the patch
	MatchAnnotations map[string]string `json:"matchAnnotations,omitempty"`

	// LabelSelector is a map of key-value pairs.
	// A Resource matches it will be applied the patch.
	LabelSelector map[string]string `json:"labelSelector,omitempty"`
}

var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) (err error) {
	p.Patches = nil
	p.ldr = ldr
	p.rf = rf
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Transform(m resmap.ResMap) error {
	for _, patch := range p.Patches {
		if patch.Type != "StrategicMergePatch" {
			return errors.New(fmt.Sprintf("StrategicMergePatch is the only supported patch type, \"%s\" given", patch.Type))
		}

		matchingMap := getMatchingResources(m, patch.Target)
		if matchingMap.Size() < 1 {
			return fmt.Errorf("0 resources matching %v", patch.Target)
		}

		patchTemplate, err := loadPatchTemplate(p, patch)

		if err != nil {
			return err
		}

		patches := createPatches(matchingMap, patchTemplate)

		err = applyPatches(m, p, patches)

		if err != nil {
			return err
		}
	}

	return nil
}

func getMatchingResources(m resmap.ResMap, target PatchTarget) resmap.ResMap {
	matchingMap := resmap.New()

	for _, resource := range m.Resources() {
		if matchesResource(resource, target) {
			matchingMap.Append(resource)
		}
	}

	return matchingMap
}

func matchesResource(resource *resource.Resource, target PatchTarget) bool {
	kunstructured := resource.Kunstructured

	if target.Group != "" && target.Group != kunstructured.GetGvk().Group {
		return false
	}

	if target.Version != "" && target.Version != kunstructured.GetGvk().Version {
		return false
	}

	if target.Kind != "" && target.Kind != kunstructured.GetGvk().Kind {
		return false
	}

	if target.Name != "" {
		match, _ := filepath.Match(target.Name, kunstructured.GetName())
		if !match {
			return false
		}
	}

	if !mapMatches(kunstructured.GetLabels(), target.LabelSelector) {
		return false
	}

	if !mapMatches(kunstructured.GetAnnotations(), target.MatchAnnotations) {
		return false
	}

	return true
}

func mapMatches(check map[string]string, expected map[string]string) bool {
	for k, v := range expected {
		if check[k] != v {
			return false
		}
	}

	return true
}

func createPatches(m resmap.ResMap, patchTemplate *resource.Resource) []*resource.Resource {
	var patches []*resource.Resource

	for _, resource := range m.Resources() {
		newPatch := patchTemplate.DeepCopy()
		newPatch.Kunstructured.SetName(resource.GetName())

		patches = append(patches, newPatch)
	}

	return patches
}

func loadPatchTemplate(p *plugin, patch Patch) (*resource.Resource, error) {
	patchResources, err := p.rf.RF().SliceFromPatches(p.ldr, []types.PatchStrategicMerge{types.PatchStrategicMerge(patch.Path)})
	if err != nil {
		return nil, err
	}

	if len(patchResources) != 1 {
		return nil, errors.New(fmt.Sprintf("Expected to find 1 resource in file \"%s\", found %d", patch.Path, len(patchResources)))
	}

	return patchResources[0], nil
}

func applyPatches(m resmap.ResMap, p *plugin, patches []*resource.Resource) error {
	patchTransformer, err := transformer.NewFactoryImpl().MakePatchTransformer(patches, p.rf.RF())
	if err != nil {
		return err
	}

	err = patchTransformer.Transform(m)
	if err != nil {
		return err
	}

	return nil
}
