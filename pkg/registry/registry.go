/*
Copyright 2025.

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

package registry

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type RegistryChecker struct {
	Keychain authn.Keychain
}

func NewRegistryChecker() *RegistryChecker {
	return &RegistryChecker{
		Keychain: authn.DefaultKeychain,
	}
}

func (r *RegistryChecker) GetImageDigest(ctx context.Context, imageRef, tag string, authSecret *authn.Secret) (string, error) {
	refStr := fmt.Sprintf("%s:%s", imageRef, tag)
	ref, err := name.ParseReference(refStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse image reference: %w", err)
	}

	var opts []remote.Option
	if authSecret != nil {
		authenticator := &authn.Basic{
			Username: authSecret.Username,
			Password: authSecret.Password,
		}
		opts = append(opts, remote.WithAuth(authenticator))
	} else {
		opts = append(opts, remote.WithAuthFromKeychain(r.Keychain))
	}

	desc, err := remote.Get(ref, opts...)
	if err != nil {
		return "", fmt.Errorf("failed to get image descriptor: %w", err)
	}

	return desc.Digest.String(), nil
}