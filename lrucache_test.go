package main

import (
	"testing"
)

func TestLruCache(t *testing.T) {
	type LruOperation int
	const (
		set LruOperation = iota
		get
	)
	tests := []struct {
		op             LruOperation
		key            CacheKey
		value          CacheValue
		expectedResult string
		expectErr      bool
	}{
		{
			op:             set,
			key:            1,
			value:          "1",
			expectedResult: "1 ",
			expectErr:      false,
		},
		{
			op:             set,
			key:            2,
			value:          "2",
			expectedResult: "2 1 ",
			expectErr:      false,
		},
		{
			op:             set,
			key:            3,
			value:          "3",
			expectedResult: "3 2 1 ",
			expectErr:      false,
		},
		{
			op:             set,
			key:            4,
			value:          "4",
			expectedResult: "4 3 2 ",
			expectErr:      false,
		},
		{
			op:             get,
			key:            2,
			expectedResult: "2 4 3 ",
			expectErr:      false,
		},
		{
			op:             get,
			key:            5,
			expectedResult: "key not found",
			expectErr:      true,
		},
		{
			op:             set,
			key:            5,
			value:          "5",
			expectedResult: "5 2 4 ",
			expectErr:      false,
		},
		{
			op:             set,
			key:            4,
			value:          "4",
			expectedResult: "4 5 2 ",
			expectErr:      false,
		},
		{
			op:             get,
			key:            2,
			expectedResult: "2 4 5 ",
			expectErr:      false,
		},
	}
	lru := New(3)
	defer lru.Close()
	for _, tt := range tests {
		switch tt.op {
		case set:
			err := lru.Set(tt.key, tt.value)
			if err != nil {
				if !tt.expectErr {
					t.Errorf("expected no error but got %v", err)
					continue
				}
			}
			if lru.String() != tt.expectedResult {
				t.Errorf("expected %s but got %s", tt.expectedResult, lru.String())
			}
		case get:
			_, err := lru.Get(tt.key)
			if err != nil {
				if !tt.expectErr {
					t.Errorf("expected no error but got %v", err)
				}
				if tt.expectErr && tt.expectedResult != err.Error() {
					t.Errorf("expected error %s but got this error %s", tt.expectedResult, err.Error())
				}
				continue
			}
			if lru.String() != tt.expectedResult {
				t.Errorf("expected %s but got %s", tt.expectedResult, lru.String())
			}
		}
	}
}
