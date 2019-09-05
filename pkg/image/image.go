package image

import (
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/opencontainers/go-digest"
	"os"
	"time"
)

type diffID digest.Digest

// Image is the image's config object
type Image struct {
	// ID is a unique 64 character identifier of the image
	ID string `json:"id,omitempty"`
	// Parent is the ID of the parent image
	Parent string `json:"parent,omitempty"`
	// Comment is the commit message that was set when committing the image
	Comment string `json:"comment,omitempty"`
	// Created is the timestamp at which the image was created
	Created time.Time `json:"created"`
	// Container is the id of the container used to commit
	Container string `json:"container,omitempty"`
	// ContainerConfig is the configuration of the container that is committed into the image
	ContainerConfig container.Config `json:"container_config,omitempty"`
	// DockerVersion specifies the version of Docker that was used to build the image
	DockerVersion string         `json:"docker_version,omitempty"`
	History       []imageHistory `json:"history,omitempty"`
	// Author is the name of the author that was specified when committing the image
	Author string `json:"author,omitempty"`
	// Config is the configuration of the container received from the client
	Config *container.Config `json:"config,omitempty"`
	// Architecture is the hardware that the image is built and runs on
	Architecture string `json:"architecture,omitempty"`
	// OS is the operating system used to build and run the image
	OS string `json:"os,omitempty"`
	// Size is the total size of the image including all layers it is composed of
	Size   int64        `json:",omitempty"`
	RootFS *imageRootFS `json:"rootfs,omitempty"`

	// rawJSON caches the immutable JSON associated with this image.
	rawJSON []byte
}

type imageRootFS struct {
	Type      string   `json:"type"`
	DiffIDs   []diffID `json:"diff_ids,omitempty"`
	BaseLayer string   `json:"base_layer,omitempty"`
}

type imageHistory struct {
	Created    time.Time `json:"created"`
	Author     string    `json:"author,omitempty"`
	CreatedBy  string    `json:"created_by,omitempty"`
	Comment    string    `json:"comment,omitempty"`
	EmptyLayer bool      `json:"empty_layer,omitempty"`
}

// Manifest is the image manifest struct
type Manifest struct {
	Config   string   `json:"Config,omitempty"`
	Layers   []string `json:"Layers,omitempty"`
	RepoTags []string `json:"RepoTags,omitempty"`
}

// Manifests is an array of Manifest
type Manifests []Manifest

// Repo is a image repo
type Repo struct {
	Tag           string
	DockerVersion string
	Created       string
	Layers        []*Layer
}

// Layer is a image layer
type Layer struct {
	Root    string
	Size    int
	Command string
	Files   []File
}

// File is a layer file object
type File struct {
	Name string // base name of the file
	Path string
	Data os.FileInfo

	Children []File
}

// RawJSON returns the immutable JSON associated with the image.
func (img *Image) RawJSON() []byte {
	return img.rawJSON
}

// NewFromJSON creates an Image configuration from json.
func NewFromJSON(src []byte) (*Image, error) {
	img := &Image{}
	if err := json.Unmarshal(src, img); err != nil {
		return nil, err
	}
	if img.RootFS == nil {
		return nil, errors.New("invalid image JSON, no RootFS key")
	}
	img.rawJSON = src
	return img, nil
}
