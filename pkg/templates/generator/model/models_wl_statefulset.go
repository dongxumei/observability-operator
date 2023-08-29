package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type CompStatefulSet struct {
	GenericModel         `yaml:",inline"`
	Selector             map[string]string `yaml:"selector,omitempty"`
	Replicas             int               `yaml:"replicas,omitempty"`
	UpdateStrategy       string            `yaml:"updateStrategy,omitempty"`
	PodManagementPolicy  string            `yaml:"podManagementPolicy,omitempty"`
	MinReadySeconds      int32             `yaml:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int32            `yaml:"revisionHistoryLimit,omitempty"`
}

func (a *CompStatefulSet) AvailableOperations() map[generator.Operation][]string {
	return map[generator.Operation][]string{}
}

func (g *CompStatefulSet) Type() string {
	return ""
}

func (a *CompStatefulSet) Args() []string {
	return []string{A_Labels, A_Selector, A_Replica, A_UpdateStrategy, A_PodManagementPolicy, A_MinReadySeconds, A_RevisionHistoryLimit}
}

func (a *CompStatefulSet) ArgsExample() string {
	return "label1:values;label2:value2 label1:values replicas updateStrategy podManagementPolicy minReadySeconds RevisionHistoryLimit"
}

func (a *CompStatefulSet) ParseArgs(input string) (err error) {
	sections := strings.Split(input, " ")
	if len(sections) != len(a.Args()) {
		return fmt.Errorf("invalid args, required pattern: %s", a.ArgsExample())
	}

	if sections[0] != "_" {
		err = convertAndMergeLabels(a.Labels, sections[0])
		if err != nil {
			return err
		}
	}

	if sections[1] != "_" {
		err = convertAndMergeLabels(a.Selector, sections[1])
		if err != nil {
			return err
		}
	}

	if sections[2] != "_" {
		a.Replicas, err = strconv.Atoi(sections[2])
		if err != nil {
			return err
		}
	}

	if sections[3] != "_" {
		val, err := convertUpdateStrategy(sections[3])
		if err != nil {
			return err
		}
		a.UpdateStrategy = val
	}

	if sections[4] != "_" {
		val, err := convertPodManagementPolicy(sections[4])
		if err != nil {
			return err
		}
		a.PodManagementPolicy = val
	}

	if sections[5] != "_" {
		val, err := strconv.Atoi(sections[5])
		if err != nil {
			return err
		}
		a.MinReadySeconds = int32(val)
	}

	if sections[6] != "_" {
		val, err := strconv.Atoi(sections[6])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.RevisionHistoryLimit = &int32Val
	}

	return nil
}

func convertPodManagementPolicy(value string) (string, error) {
	if value == "OrderedReady" || value == "Parallel" {
		return value, nil
	}
	return "", fmt.Errorf("invalid values")
}
