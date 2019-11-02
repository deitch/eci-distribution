package registry

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes/docker"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

func Push(image string, artifact *Artifact, verbose bool, writer io.Writer) (string, error) {
	ctx := context.Background()
	resolver := docker.NewResolver(docker.ResolverOptions{})

	// Go through each file type in the registry and add the appropriate file type and path, along with annotations
	fileStore := content.NewFileStore("")
	defer fileStore.Close()

	pushContents := []ocispec.Descriptor{}
	var (
		desc            ocispec.Descriptor
		mediaType       string
		customMediaType string
		role            string
		name            string
		path            string
		err             error
		pushOpts        []oras.PushOpt
	)

	if artifact.Kernel != "" {
		role = RoleKernel
		name = "kernel"
		customMediaType = MimeTypeECIKernel
		path = artifact.Kernel
		mediaType = GetLayerMediaType(customMediaType, artifact.Legacy)
		desc, err = fileStore.Add(name, mediaType, path)
		if err != nil {
			return "", fmt.Errorf("error adding %s at %s: %v", name, path, err)
		}
		desc.Annotations[AnnotationMediaType] = customMediaType
		desc.Annotations[AnnotationRole] = role
		desc.Annotations[ocispec.AnnotationTitle] = name
		pushContents = append(pushContents, desc)
	}

	if artifact.Initrd != "" {
		role = RoleInitrd
		name = "initrd"
		customMediaType = MimeTypeECIInitrd
		path = artifact.Initrd
		mediaType = GetLayerMediaType(customMediaType, artifact.Legacy)
		desc, err = fileStore.Add(name, mediaType, path)
		if err != nil {
			return "", fmt.Errorf("error adding %s at %s: %v", name, path, err)
		}
		desc.Annotations[AnnotationMediaType] = customMediaType
		desc.Annotations[AnnotationRole] = role
		desc.Annotations[ocispec.AnnotationTitle] = name
		pushContents = append(pushContents, desc)
	}

	if disk := artifact.Root; disk != nil {
		name := "root"
		role = RoleRootDisk
		customMediaType = TypeToMime[disk.Type]
		path = disk.Path
		mediaType = GetLayerMediaType(customMediaType, artifact.Legacy)
		desc, err = fileStore.Add(name, mediaType, path)
		if err != nil {
			return "", fmt.Errorf("error adding %s disk at %s: %v", name, path, err)
		}
		desc.Annotations[AnnotationMediaType] = customMediaType
		desc.Annotations[AnnotationRole] = role
		desc.Annotations[ocispec.AnnotationTitle] = name
		pushContents = append(pushContents, desc)
	}
	for i, disk := range artifact.Disks {
		if disk != nil {
			name := fmt.Sprintf("disk-%d", i)
			role = RoleAdditionalDisk
			customMediaType = TypeToMime[disk.Type]
			path = disk.Path
			mediaType = GetLayerMediaType(customMediaType, artifact.Legacy)
			desc, err = fileStore.Add(name, mediaType, path)
			if err != nil {
				return "", fmt.Errorf("error adding %s disk at %s: %v", name, path, err)
			}
			desc.Annotations[AnnotationMediaType] = customMediaType
			desc.Annotations[AnnotationRole] = role
			desc.Annotations[ocispec.AnnotationTitle] = name
			pushContents = append(pushContents, desc)
		}
	}

	if verbose {
		pushOpts = append(pushOpts, oras.WithPushBaseHandler(pushStatusTrack(writer)))
	}

	// push the data
	desc, err = oras.Push(ctx, resolver, image, fileStore, pushContents, pushOpts...)
	if err != nil {
		return "", err
	}
	return string(desc.Digest), nil
}

func pushStatusTrack(writer io.Writer) images.Handler {
	var printLock sync.Mutex
	return images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		if name, ok := content.ResolveName(desc); ok {
			printLock.Lock()
			defer printLock.Unlock()
			writer.Write([]byte(fmt.Sprintf("Uploading %s %s\n", desc.Digest.Encoded()[:12], name)))
		}
		return nil, nil
	})
}
