package rpeat

import (
	"os"
)

type Artifacts struct {
	Artifact    []Artifact `json:"Artifact,omitempty" xml:"Artifact,omitempty"`
	Permissions Permission `json:"-" xml:"-"`
	Retain      string     `json:"Retain,omitempty" xml:"Retain,omitempty"`
	Type        string     `json:"Type,omitempty" xml:"Type,omitempty"`
}

type Artifact struct {
	// An Artifact represents a single file output of a process
	// which is triggered upon exit of process with action including
	// asynchronous upload for secure storage and access

	// Src is local file generated from rpeat® process
	Src string `json:"Src" xml:"Src"`
	// Dst is path to external object including resource URI
	//  e.g.
	//
	//  s3://
	//  file://
	//  rpeatio[:file]//
	//  rpeatio:slack//
	//
	Dst string `json:"Dst" xml:"Dst"`

	// Artifact overrided
	Permissions Permission `json:"-" xml:"-"`
	Retain      string     `json:"Retain" xml:"Retain"`
	Type        string     `json:"Type" xml:"Type"`

	// Unique ID of artifact
	UUID string      `json:"UUID" xml:"UUID"`
	Stat os.FileInfo `json:"Stat" xml:"Stat"`

	JobUUID string `json:"JobUUID" xml:"JobUUID"`
	RunUUID string `json:"RunUUID" xml:"RunUUID"`
	SysPerm string `json:"SysPerm" xml:"SysPerm"`
}

func (job *Job) HasArtifacts() bool {
	artifacts := job.Artifacts
	if len(artifacts.Artifact) > 0 {
		return true
	}
	return false
}

// read + upload artifacts
func (job *Job) processArtifacts() {
}
