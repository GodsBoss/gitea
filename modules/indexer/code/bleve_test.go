// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package code

import (
	"os"
	"path/filepath"
	"testing"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	models.MainTest(m, filepath.Join("..", "..", ".."))
}

func TestIndexAndSearch(t *testing.T) {
	models.PrepareTestEnv(t)

	dir := "./bleve.index"
	os.RemoveAll(dir)

	setting.Indexer.RepoIndexerEnabled = true
	idx, _, err := NewBleveIndexer(dir)
	if err != nil {
		idx.Close()
		log.Fatal("indexer.Init: %v", err)
	}

	err = idx.Index(1)
	assert.NoError(t, err)

	var (
		keywords = []struct {
			Keyword string
			IDs     []int64
		}{
			{
				Keyword: "Description",
				IDs:     []int64{1},
			},
			{
				Keyword: "repo1",
				IDs:     []int64{1},
			},
			{
				Keyword: "non-exist",
				IDs:     []int64{},
			},
		}
	)

	for _, kw := range keywords {
		total, res, err := idx.Search(nil, kw.Keyword, 1, 10)
		assert.NoError(t, err)
		assert.EqualValues(t, len(kw.IDs), total)

		var ids = make([]int64, 0, len(res))
		for _, hit := range res {
			ids = append(ids, hit.RepoID)
		}
		assert.EqualValues(t, kw.IDs, ids)
	}
}
