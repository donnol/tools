package db

import (
	"context"
	"fmt"

	"github.com/donnol/tools/worker"
)

// Batch 批量处理查询到的数据
func Batch[S Storer, F Finder[R], R any](db S, finder F, batchNum int, handler func([]R) error) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}
	batch := make([]R, 0, batchNum)
	for rows.Next() {
		t, fields := finder.NewScanObjAndFields(colTypes)
		if err = rows.Scan(fields...); err != nil {
			return
		}

		batch = append(batch, t)
		if batchNum > 0 && len(batch) >= batchNum {
			if err = handler(batch); err != nil {
				err = fmt.Errorf("batch handle failed %w", err)
				return
			}

			batch = make([]R, 0, batchNum)
		}
	}
	if err = rows.Err(); err != nil {
		return
	}

	if len(batch) != 0 {
		if err = handler(batch); err != nil {
			err = fmt.Errorf("batch handle failed %w", err)
			return
		}
	}

	return
}

// BatchConcurrent 并发批量处理查询到的数据
func BatchConcurrent[S Storer, F Finder[R], R any](db S, finder F, batchNum int, handler func([]R) error, concNum int) (err error) {
	query, args := finder.Query()
	rows, err := db.QueryContext(context.TODO(), query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return
	}

	batchWorker := worker.New(concNum)
	batchWorker.Start()
	defer batchWorker.Stop()

	batch := make([]R, 0, batchNum)
	for rows.Next() {
		t, fields := finder.NewScanObjAndFields(colTypes)
		if err = rows.Scan(fields...); err != nil {
			return
		}

		batch = append(batch, t)
		if batchNum > 0 && len(batch) >= batchNum {
			batchCopy := batch
			if err = batchWorker.Push(*worker.NewJob(
				worker.DoWithCtx(func(ctx context.Context) error {
					if err = handler(batchCopy); err != nil {
						err = fmt.Errorf("batch handle failed %w", err)
						return err
					}
					return nil
				}),
				0,
				nil,
			)); err != nil {
				return
			}

			batch = make([]R, 0, batchNum)
		}
	}
	if err = rows.Err(); err != nil {
		return
	}

	if len(batch) != 0 {
		batchCopy := batch
		if err = batchWorker.Push(*worker.NewJob(
			worker.DoWithCtx(func(ctx context.Context) error {
				if err = handler(batchCopy); err != nil {
					err = fmt.Errorf("batch handle failed %w", err)
					return err
				}
				return nil
			}),
			0,
			nil,
		)); err != nil {
			return
		}
	}

	return
}
