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
func (p PostgresMetricStorer) StoreSingle(ctx context.Context, metric models.Metric) (*models.Metric, error) {

	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return nil, txErr
	}

	switch metric.MType {
	case constants.Gauge:
		stmt, prepErr := tx.PrepareContext(ctx,
			`
            INSERT INTO metrics (id, mtype, delta, value)
            VALUES ($1, $2, NULL, $3)
            ON CONFLICT (id) DO UPDATE
            SET mtype = EXCLUDED.mtype, value = EXCLUDED.value;`)
		if prepErr != nil {
			p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
			tx.Rollback()
			return nil, prepErr
		}
		defer stmt.Close()

		rows, queryErr := stmt.QueryContext(ctx, metric.ID, metric.MType, metric.Value)
		if queryErr != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", queryErr.Error())
			tx.Rollback()
			return nil, queryErr
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
			tx.Rollback()
			return nil, rowsErr
		}
		tx.Commit()
		return &metric, nil
	case constants.Counter:
		stmtSelect, prepSelectErr := tx.PrepareContext(ctx, `SELECT id, mtype, delta FROM metrics WHERE id = $1;`)
		if prepSelectErr != nil {
			p.logger.Errorw("failed to prepare statement", "err", prepSelectErr.Error())
			tx.Rollback()
			return nil, prepSelectErr
		}
		defer stmtSelect.Close()

		var existingMetric models.Metric

		row := stmtSelect.QueryRowContext(ctx, metric.ID)

		scanErr := row.Scan(&existingMetric.ID, &existingMetric.MType, &existingMetric.Delta)
		if scanErr != nil {
			if !errors.Is(scanErr, sql.ErrNoRows) {
				p.logger.Errorw("error retrieving metric", "metric_id", metric.ID, "error", scanErr.Error())
				tx.Rollback()
				return nil, scanErr
			}
		}

		newMetric, updateErr := helpers.UpdateCounterMetric(existingMetric, metric)
		if updateErr != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", updateErr.Error())
			tx.Rollback()
			return nil, updateErr
		}

		stmtInsert, prepInsertErr := tx.PrepareContext(ctx,
			`INSERT INTO metrics (id, mtype, delta, value)
					VALUES ($1, $2, $3, NULL)
					ON CONFLICT (id) DO UPDATE
					SET mtype = EXCLUDED.mtype, delta = EXCLUDED.delta;`)
		if prepInsertErr != nil {
			p.logger.Errorw("failed to prepare statement", "err", prepInsertErr.Error())
			tx.Rollback()
			return nil, prepInsertErr
		}
		defer stmtInsert.Close()

		rows, queryInsertErr := stmtInsert.QueryContext(ctx, newMetric.ID, newMetric.MType, newMetric.Delta)
		if queryInsertErr != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", queryInsertErr.Error())
			tx.Rollback()
			return nil, queryInsertErr
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
			tx.Rollback()
			return nil, rowsErr
		}
		tx.Commit()
		return &newMetric, nil

	default:
		unsupportedErr := fmt.Errorf("unsupported Metric type: %s", metric.MType)
		p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", unsupportedErr.Error())
		tx.Rollback()
		return nil, unsupportedErr
	}
}

func (p PostgresMetricStorer) StoreSlice(ctx context.Context, newMetrics []models.Metric) error {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return txErr
	}
	for _, metric := range newMetrics {
		switch metric.MType {
		case constants.Gauge:
			stmt, prepErr := tx.PrepareContext(ctx,
				`
            INSERT INTO metrics (id, mtype, delta, value)
            VALUES ($1, $2, NULL, $3)
            ON CONFLICT (id) DO UPDATE
            SET mtype = EXCLUDED.mtype, value = EXCLUDED.value;`)
			if prepErr != nil {
				p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
				tx.Rollback()
				return prepErr
			}

			rows, queryErr := stmt.QueryContext(ctx, metric.ID, metric.MType, metric.Value)
			if queryErr != nil {
				p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", queryErr.Error())
				tx.Rollback()
				return queryErr
			}
			if rowsErr := rows.Err(); rowsErr != nil {
				p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
				tx.Rollback()
				return rowsErr
			}
		case constants.Counter:
			stmtSelect, prepSelectErr := tx.PrepareContext(ctx, `SELECT id, mtype, delta FROM metrics WHERE id = $1;`)
			if prepSelectErr != nil {
				p.logger.Errorw("failed to prepare statement", "err", prepSelectErr.Error())
				tx.Rollback()
				return prepSelectErr
			}

			var existingMetric models.Metric

			row := stmtSelect.QueryRowContext(ctx, metric.ID)

			scanErr := row.Scan(&existingMetric.ID, &existingMetric.MType, &existingMetric.Delta)
			if scanErr != nil {
				if !errors.Is(scanErr, sql.ErrNoRows) {
					p.logger.Errorw("error retrieving metric", "metric_id", metric.ID, "error", scanErr.Error())
					tx.Rollback()
					return scanErr
				}
			}

			newMetric, updateErr := helpers.UpdateCounterMetric(existingMetric, metric)
			if updateErr != nil {
				p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", updateErr.Error())
				tx.Rollback()
				return updateErr
			}

			stmtInsert, prepInsertErr := tx.PrepareContext(ctx,
				`INSERT INTO metrics (id, mtype, delta, value)
					VALUES ($1, $2, $3, NULL)
					ON CONFLICT (id) DO UPDATE
					SET mtype = EXCLUDED.mtype, delta = EXCLUDED.delta;`)
			if prepInsertErr != nil {
				p.logger.Errorw("failed to prepare statement", "err", prepInsertErr.Error())
				tx.Rollback()
				return prepInsertErr
			}

			rows, queryInsertErr := stmtInsert.QueryContext(ctx, newMetric.ID, newMetric.MType, newMetric.Delta)
			if queryInsertErr != nil {
				p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", queryInsertErr.Error())
				tx.Rollback()
				return queryInsertErr
			}
			if rowsErr := rows.Err(); rowsErr != nil {
				p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
				tx.Rollback()
				return rowsErr
			}

		default:
			unsupportedErr := fmt.Errorf("unsupported Metric type: %s", metric.MType)
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", unsupportedErr.Error())
			tx.Rollback()
			return unsupportedErr
		}
	}

	tx.Commit()
	return nil
}

// Get retrieves a metric by its ID from the database
func (p PostgresMetricStorer) Get(ctx context.Context, id string) (models.Metric, bool, error) {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return models.Metric{}, false, txErr
	}

	stmt, prepErr := tx.PrepareContext(ctx, `SELECT id, mtype, delta, value FROM metrics WHERE id = $1;`)
	if prepErr != nil {
		p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
		tx.Rollback()
		return models.Metric{}, false, prepErr
	}
	defer stmt.Close()

	var metric models.Metric

	row := stmt.QueryRowContext(ctx, id)

	scanErr := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			tx.Commit()
			return metric, false, nil
		}
		p.logger.Errorw("error getting metric", "id", id, "error", scanErr.Error())
		tx.Rollback()
		return metric, false, scanErr
	}

	tx.Commit()
	return metric, true, nil
}

func (p PostgresMetricStorer) All(ctx context.Context) (map[string]models.Metric, error) {
	tx, txErr := p.db.Begin()
	if txErr != nil {
		p.logger.Errorw("failed to begin transaction", "err", txErr.Error())
		return nil, txErr
	}

	stmt, prepErr := tx.PrepareContext(ctx, `SELECT id, mtype, delta, value FROM metrics;`)
	if prepErr != nil {
		p.logger.Errorw("failed to prepare statement", "err", prepErr.Error())
		tx.Rollback()
		return nil, prepErr
	}
	defer stmt.Close()

	rows, queryErr := stmt.QueryContext(ctx)
	if queryErr != nil {
		p.logger.Errorw("error getting all metrics", "error", queryErr.Error())
		tx.Rollback()
		return nil, queryErr
	}
	defer rows.Close()

	metrics := make(map[string]models.Metric)
	for rows.Next() {
		var metric models.Metric

		if scanErr := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); scanErr != nil {
			p.logger.Errorw("error scanning metric row", "error", scanErr.Error())
			tx.Rollback()
			return nil, scanErr
		}

		metrics[metric.ID] = metric
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		p.logger.Errorw("error during rows iteration", "error", rowsErr.Error())
		tx.Rollback()
		return nil, rowsErr
	}

	tx.Commit()
	return metrics, nil
}
