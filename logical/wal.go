package logical

import (
	"encoding/json"
	"strings"
)

// WALPrefix is the prefix within Storage where WAL entries will be written.
const WALPrefix = "wal/"

// PutWAL writes some data to the WAL.
//
// Data within the WAL that is uncommitted (CommitWAL hasn't be called)
// will be given to the rollback callback when an rollback operation is
// received, allowing the backend to clean up some partial states.
//
// The data must be JSON encodable.
//
// This returns a unique ID that can be used to reference this WAL data.
// WAL data cannot be modified. You can only add to the WAL and commit existing
// WAL entries.
func PutWAL(s Storage, data interface{}) (string, error) {
	value, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	id, err := UUID()
	if err != nil {
		return "", err
	}

	return id, s.Put(&StorageEntry{
		Key:   WALPrefix + id,
		Value: value,
	})
}

// GetWAL reads a specific entry from the WAL. If the entry doesn't exist,
// then nil value is returned.
func GetWAL(s Storage, id string) (interface{}, error) {
	entry, err := s.Get(WALPrefix + id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result interface{}
	if err := json.Unmarshal(entry.Value, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteWAL commits the WAL entry with the given ID. Once comitted,
// it is assumed that the operation was a success and doesn't need to
// be rolled back.
func DeleteWAL(s Storage, id string) error {
	return s.Delete(WALPrefix + id)
}

// ListWAL lists all the entries in the WAL.
func ListWAL(s Storage) ([]string, error) {
	keys, err := s.List(WALPrefix)
	if err != nil {
		return nil, err
	}

	for i, k := range keys {
		keys[i] = strings.TrimPrefix(k, WALPrefix)
	}

	return keys, nil
}