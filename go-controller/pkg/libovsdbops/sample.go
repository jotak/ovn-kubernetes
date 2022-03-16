package libovsdbops

import (
	"context"
	"fmt"

	libovsdbclient "github.com/ovn-org/libovsdb/client"
	libovsdb "github.com/ovn-org/libovsdb/ovsdb"

	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/nbdb"
	"github.com/ovn-org/ovn-kubernetes/go-controller/pkg/types"
)

// findSamplesByPredicate looks up Samples from the cache based on a given predicate
func findSamplesByPredicate(nbClient libovsdbclient.Client, lookupFunction func(item *nbdb.Sample) bool) ([]nbdb.Sample, error) {
	ctx, cancel := context.WithTimeout(context.Background(), types.OVSDBTimeout)
	defer cancel()
	samples := []nbdb.Sample{}
	err := nbClient.WhereCache(lookupFunction).List(ctx, &samples)
	if err != nil {
		return nil, err
	}

	return samples, nil
}

// IsEquivalentSample if it has same uuid, or if it has same name
// and external ids, or if it has same priority, direction, match
// and action.
func IsEquivalentSample(existing *nbdb.Sample, searched *nbdb.Sample) bool {
	if searched.UUID != "" && existing.UUID == searched.UUID {
		return true
	}
	return (existing.CollectorSetID == searched.CollectorSetID &&
		existing.ObsPointID == searched.ObsPointID &&
		existing.Probability == searched.Probability)
}

// findSample looks up the Sample in the cache and sets the UUID
func findSample(nbClient libovsdbclient.Client, sample *nbdb.Sample) error {
	if sample.UUID != "" && !IsNamedUUID(sample.UUID) {
		return nil
	}

	samples, err := findSamplesByPredicate(nbClient, func(item *nbdb.Sample) bool {
		return IsEquivalentSample(item, sample)
	})

	if err != nil {
		return fmt.Errorf("can't find sample by equivalence %+v: %v", *sample, err)
	}
	if len(samples) > 1 {
		return fmt.Errorf("unexpectedly found multiple equivalent Samples: %+v", samples)
	}
	if len(samples) == 0 {
		return libovsdbclient.ErrNotFound
	}

	sample.UUID = samples[0].UUID
	return nil
}

func BuildSample(namedUUID string, collectorSetID, obsPointID, probability int) *nbdb.Sample {
	return &nbdb.Sample{
		UUID:           namedUUID,
		CollectorSetID: collectorSetID,
		ObsPointID:     obsPointID,
		Probability:    probability,
	}
}
func CreateSamplesOps(nbClient libovsdbclient.Client, ops []libovsdb.Operation, acls []*nbdb.ACL, samples []*nbdb.Sample) ([]libovsdb.Operation, error) {
	if ops == nil {
		ops = []libovsdb.Operation{}
	}
	samplesByNamedUUID := map[string]*nbdb.Sample{}
	for _, sample := range samples {
		samplesByNamedUUID[sample.UUID] = sample
		err := findSample(nbClient, sample)
		if err != nil && err != libovsdbclient.ErrNotFound {
			return nil, err
		}
		// If Sample does not exist, create it
		if err == libovsdbclient.ErrNotFound {
			op, err := nbClient.Create(sample)
			if err != nil {
				return nil, err
			}
			ops = append(ops, op...)
		}
	}
	for _, acl := range acls {
		if acl.Sample != nil {
			sample, ok := samplesByNamedUUID[*acl.Sample]
			if !ok {
				return nil, fmt.Errorf("Unable to find Sample object for ACL: %v", acl)
			}
			// If the Sample already existed it's UUID has now been updated
			// to the one in the cache, update the ACL's
			acl.Sample = &sample.UUID
		}
	}
	return ops, nil
}
