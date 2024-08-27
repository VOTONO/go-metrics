package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/VOTONO/go-metrics/internal/constants"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
)

type PostgresMetricStorer struct {
	logger *zap.SugaredLogger
	db     *sql.DB
}

func NewPostgresMetricStorer(logger *zap.SugaredLogger, db *sql.DB) (*PostgresMetricStorer, error) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS metrics (
        id TEXT PRIMARY KEY,
        mtype TEXT,
        delta BIGINT,
        value DOUBLE PRECISION
    );`

	_, execErr := db.Exec(createTableSQL)
	if execErr != nil {
		logger.Errorw("failed to create table", "err", execErr.Error())
		return nil, execErr
	}

	return &PostgresMetricStorer{
		logger: logger,
		db:     db,
	}, nil
}

// StoreSingle inserts or updates a metric in the database
func (p PostgresMetricStorer) StoreSingle(ctx context.Context, newMetric models.Metric) (*models.Metric, error) {

	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return nil, txErr
	}

	var updatedMetric *models.Metric
	var err error
	defer func() {
		var deferErr error
		if err != nil {
			deferErr = tx.Rollback()
		} else {
			deferErr = tx.Commit()
		}
		if deferErr != nil {
			p.logger.Errorw("failed to rollback or commit transaction", "err", deferErr.Error())
		}
	}()

	updatedMetric, err = p.insertOrUpdateMetric(ctx, tx, newMetric)
	if err != nil {
		return nil, err
	}

	return updatedMetric, nil
}

func (p PostgresMetricStorer) StoreSlice(ctx context.Context, newMetrics []models.Metric) error {
	filteredMetrics, filterErr := helpers.ProcessMetricsDuplicates(newMetrics)
	if filterErr != nil {
		return filterErr
	}
	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return txErr
	}
	var err error
	defer func() {
		var deferErr error
		if err != nil {
			deferErr = tx.Rollback()
		} else {
			deferErr = tx.Commit()
		}
		if deferErr != nil {
			p.logger.Errorw("failed to rollback or commit transaction", "err", deferErr.Error())
		}
	}()

	for _, newMetric := range filteredMetrics {
		_, err = p.insertOrUpdateMetric(ctx, tx, newMetric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p PostgresMetricStorer) insertOrUpdateMetric(ctx context.Context, tx *sql.Tx, newMetric models.Metric) (*models.Metric, error) {

	switch newMetric.MType {
	case constants.Gauge:
		metric, err := p.insertMetric(ctx, tx, newMetric)
		if err != nil {
			return nil, err
		}
		return metric, nil
	case constants.Counter:
		updatedMetric, err := p.insertCounterMetric(ctx, newMetric, tx)
		if err != nil {
			return nil, err
		}
		return updatedMetric, nil

	default:
		err := fmt.Errorf("unsupported Metric type: %s", newMetric.MType)
		p.logger.Errorw("error storing Metric", "metric_id", newMetric.ID, "error", err.Error())
		return nil, err
	}
}

func (p PostgresMetricStorer) insertCounterMetric(ctx context.Context, newMetric models.Metric, tx *sql.Tx) (*models.Metric, error) {
	existingMetric, found, getErr := p.getMetricByID(ctx, newMetric.ID)
	if getErr != nil {
		return nil, getErr
	}

	var insertedMetric *models.Metric
	var insertErr error

	if !found {
		insertedMetric, insertErr = p.insertMetric(ctx, tx, newMetric)
		if insertErr != nil {
			return nil, insertErr
		}
		return insertedMetric, nil
	}

	updatedMetric, updateErr := helpers.UpdateCounterMetric(existingMetric, newMetric)
	if updateErr != nil {
		p.logger.Errorw("error storing Metric", "metric_id", newMetric.ID, "error", updateErr.Error())
		return nil, updateErr
	}

	insertedMetric, insertErr = p.insertMetric(ctx, tx, updatedMetric)
	if insertErr != nil {
		return nil, insertErr
	}
	return insertedMetric, nil
}

func (p PostgresMetricStorer) insertMetric(ctx context.Context, tx *sql.Tx, metric models.Metric) (*models.Metric, error) {
	stmt, prepErr := tx.PrepareContext(ctx,
		`
            INSERT INTO metrics (id, mtype, delta, value)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (id) DO UPDATE
            SET mtype = EXCLUDED.mtype, delta = EXCLUDED.delta, value = EXCLUDED.value;`)
	if prepErr != nil {
		p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
		return nil, prepErr
	}

	rows, queryErr := stmt.QueryContext(ctx, metric.ID, metric.MType, metric.Delta, metric.Value)
	if queryErr != nil {
		p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", queryErr.Error())
		return nil, queryErr
	}
	if rowsErr := rows.Err(); rowsErr != nil {
		p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
		return nil, rowsErr
	}
	return &metric, nil
}

// Get retrieves a metric by its ID from the database
func (p PostgresMetricStorer) Get(ctx context.Context, id string) (models.Metric, bool, error) {
	metric, found, getErr := p.getMetricByID(ctx, id)
	return metric, found, getErr
}

func (p PostgresMetricStorer) getMetricByID(ctx context.Context, id string) (models.Metric, bool, error) {
	var metric models.Metric

	stmt, prepErr := p.db.PrepareContext(ctx, `SELECT id, mtype, delta, value FROM metrics WHERE id = $1;`)
	if prepErr != nil {
		p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
		return metric, false, prepErr
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	scanErr := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return metric, false, nil
		}
		p.logger.Errorw("error getting metric", "id", id, "error", scanErr.Error())
		return metric, false, scanErr
	}
	return metric, true, nil
}

func (p PostgresMetricStorer) All(ctx context.Context) (map[string]models.Metric, error) {

	stmt, prepErr := p.db.PrepareContext(ctx, `SELECT id, mtype, delta, value FROM metrics;`)
	if prepErr != nil {
		p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
		return nil, prepErr
	}
	defer stmt.Close()

	rows, queryErr := stmt.QueryContext(ctx)
	if queryErr != nil {
		p.logger.Errorw("error getting all metrics", "error", queryErr.Error())
		return nil, queryErr
	}
	defer rows.Close()

	metrics := make(map[string]models.Metric)
	for rows.Next() {
		var metric models.Metric

		if scanErr := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); scanErr != nil {
			p.logger.Errorw("error scanning metric row", "error", scanErr.Error())
			return nil, scanErr
		}

		metrics[metric.ID] = metric
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
		return nil, rowsErr
	}

	return metrics, nil
}
