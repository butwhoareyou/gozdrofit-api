package gozdrofitapi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDate_MarshalJSON(t *testing.T) {
	date := Date{time.Date(2021, 1, 1, 10, 10, 15, 0, time.UTC)}

	actual, err := marshal(date)

	require.NoError(t, err)
	assert.Equal(t, actual, "\"2021-01-01\"")
}

func TestDate_UnmarshalJSON(t *testing.T) {
	dateString := []byte("\"2021-01-01\"")

	var actual Date
	err := json.Unmarshal(dateString, &actual)

	require.NoError(t, err)
	assert.Equal(t, actual, Date{time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)})

}

func TestDateTime_MarshalJSON(t *testing.T) {
	date := DateTime{time.Date(1975, 12, 8, 10, 10, 15, 0, time.UTC)}

	actual, err := marshal(date)

	require.NoError(t, err)
	assert.Equal(t, actual, "\"1975-12-08T10:10:15\"")
}

func TestDateTime_UnmarshalJSON(t *testing.T) {
	dateString := []byte("\"2021-01-01T01:15:34\"")

	var actual DateTime
	err := json.Unmarshal(dateString, &actual)

	require.NoError(t, err)
	assert.Equal(t, actual, DateTime{time.Date(2021, 1, 1, 1, 15, 34, 0, time.UTC)})
}

func marshal(payload interface{}) (string, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(jsonPayload), nil
}
