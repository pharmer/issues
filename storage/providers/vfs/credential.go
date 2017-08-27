package vfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/appscode/pharmer/api"
	"github.com/appscode/pharmer/storage"
	"github.com/graymeta/stow"
)

type CredentialFileStore struct {
	container stow.Container
	prefix    string
}

var _ storage.CredentialStore = &CredentialFileStore{}

func (s *CredentialFileStore) resourceHome() string {
	return filepath.Join(s.prefix, "credentials")
}

func (s *CredentialFileStore) resourceID(name string) string {
	return filepath.Join(s.resourceHome(), name+".json")
}

func (s *CredentialFileStore) List(opts api.ListOptions) ([]*api.Credential, error) {
	result := make([]*api.Credential, 0)
	cursor := stow.CursorStart
	for {
		page, err := s.container.Browse(s.resourceHome(), "/", cursor, pageSize)
		if err != nil {
			return nil, fmt.Errorf("Failed to list credentials. Reason: %v", err)
		}
		for _, item := range page.Items {
			r, err := item.Open()
			if err != nil {
				return nil, fmt.Errorf("Failed to list credentials. Reason: %v", err)
			}
			var obj api.Credential
			err = json.NewDecoder(r).Decode(&obj)
			if err != nil {
				return nil, fmt.Errorf("Failed to list credentials. Reason: %v", err)
			}
			result = append(result, &obj)
			r.Close()
		}
		cursor = page.Cursor
		if stow.IsCursorEnd(cursor) {
			break
		}
	}
	return result, nil
}

func (s *CredentialFileStore) Get(name string) (*api.Credential, error) {
	if name == "" {
		return nil, errors.New("Missing credential name")
	}

	item, err := s.container.Item(s.resourceID(name))
	if err != nil {
		return nil, fmt.Errorf("Credential `%s` does not exist. Reason: %v", name, err)
	}

	r, err := item.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var existing api.Credential
	err = json.NewDecoder(r).Decode(&existing)
	if err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *CredentialFileStore) Create(obj *api.Credential) (*api.Credential, error) {
	if obj == nil {
		return nil, errors.New("Missing credential")
	} else if obj.Name == "" {
		return nil, errors.New("Missing credential name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)

	_, err = s.container.Item(id)
	if err == nil {
		return nil, fmt.Errorf("Credential `%s` already exists", obj.Name)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, &buf, int64(buf.Len()), nil)
	return obj, err
}

func (s *CredentialFileStore) Update(obj *api.Credential) (*api.Credential, error) {
	if obj == nil {
		return nil, errors.New("Missing credential")
	} else if obj.Name == "" {
		return nil, errors.New("Missing credential name")
	}
	err := api.AssignTypeKind(obj)
	if err != nil {
		return nil, err
	}

	id := s.resourceID(obj.Name)

	_, err = s.container.Item(id)
	if err != nil {
		return nil, fmt.Errorf("Credential `%s` does not exist. Reason: %v", obj.Name, err)
	}

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(obj)
	if err != nil {
		return nil, err
	}
	_, err = s.container.Put(id, &buf, int64(buf.Len()), nil)
	return obj, err
}

func (s *CredentialFileStore) Delete(name string) error {
	if name == "" {
		return errors.New("Missing credential name")
	}
	return s.container.RemoveItem(s.resourceID(name))
}
