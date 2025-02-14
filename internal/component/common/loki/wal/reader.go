package wal

import (
	"fmt"

	"github.com/prometheus/common/model"

	"github.com/blockopsnetwork/telescope/internal/component/common/loki"
	"github.com/grafana/loki/pkg/ingester/wal"
	"github.com/grafana/loki/pkg/util"
	walUtils "github.com/grafana/loki/pkg/util/wal"
)

// ReadWAL will read all entries in the WAL located under dir. Mainly used for testing
func ReadWAL(dir string) ([]loki.Entry, error) {
	reader, closeFn, err := walUtils.NewWalReader(dir, -1)
	if err != nil {
		return nil, err
	}
	defer func() { closeFn.Close() }()

	seenSeries := make(map[uint64]model.LabelSet)
	seenEntries := []loki.Entry{}

	for reader.Next() {
		var walRec = wal.Record{}
		bytes := reader.Record()
		err = wal.DecodeRecord(bytes, &walRec)
		if err != nil {
			return nil, fmt.Errorf("error decoding wal record: %w", err)
		}

		// first read series
		for _, series := range walRec.Series {
			if _, ok := seenSeries[uint64(series.Ref)]; !ok {
				seenSeries[uint64(series.Ref)] = util.MapToModelLabelSet(series.Labels.Map())
			}
		}

		for _, entries := range walRec.RefEntries {
			for _, entry := range entries.Entries {
				labels, ok := seenSeries[uint64(entries.Ref)]
				if !ok {
					return nil, fmt.Errorf("found entry without matching series")
				}
				seenEntries = append(seenEntries, loki.Entry{
					Labels: labels,
					Entry:  entry,
				})
			}
		}

		// reset entry
		walRec.Series = walRec.Series[:]
		walRec.RefEntries = walRec.RefEntries[:]
	}

	return seenEntries, nil
}
