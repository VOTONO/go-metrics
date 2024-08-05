package repo

import (
	"database/sql"
	"fmt"
	"github.com/VOTONO/go-metrics/internal/helpers"
	"github.com/VOTONO/go-metrics/internal/models"
	"go.uber.org/zap"
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
        delta INTEGER,
        value DOUBLE PRECISION
    );`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		logger.Errorw("failed to create table", "err", err.Error())
		return nil, err
	}

	return &PostgresMetricStorer{
		logger: logger,
		db:     db,
	}, nil
}

// Store inserts or updates a metric in the database
func (p PostgresMetricStorer) Store(metric models.Metric) (*models.Metric, error) {

	switch metric.MType {
	case "gauge":
		query := `
            INSERT INTO metrics (id, mtype, delta, value)
            VALUES ($1, $2, NULL, $3)
            ON CONFLICT (id) DO UPDATE
            SET mtype = EXCLUDED.mtype, value = EXCLUDED.value;`
		_, err := p.db.Exec(query, metric.ID, metric.MType, metric.Value)
		if err != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
			return nil, err
		}
		return &metric, nil
	case "counter":
		selectExistingMetricQuery := `SELECT id, mtype, delta FROM metrics WHERE id = $1;`
		var existingMetric models.Metric

		row := p.db.QueryRow(selectExistingMetricQuery, metric.ID)

		err := row.Scan(&existingMetric.ID, &existingMetric.MType, &existingMetric.Delta)
		if err != nil {
			if err != sql.ErrNoRows {
				p.logger.Errorw("error retrieving metric", "metric_id", metric.ID, "error", err.Error())
				return nil, err
			}
		}

		newMetric, err := helpers.UpdateCounterMetric(existingMetric, metric)
		if err != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
			return nil, err
		}

		storeNewMetricQuery := `
			INSERT INTO metrics (id, mtype, delta, value)
			VALUES ($1, $2, $3, NULL)
			ON CONFLICT (id) DO UPDATE
			SET mtype = EXCLUDED.mtype, delta = EXCLUDED.delta;`

		_, er := p.db.Exec(storeNewMetricQuery, newMetric.ID, newMetric.MType, newMetric.Delta)
		if er != nil {
			p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", er.Error())
			return nil, er
		}
		return &newMetric, nil

	default:
		err := fmt.Errorf("unsupported Metric type: %s", metric.MType)
		p.logger.Errorw("error storing Metric", "metric_id", metric.ID, "error", err.Error())
		return nil, err
	}

}

// Get retrieves a metric by its ID from the database
func (p PostgresMetricStorer) Get(ID string) (models.Metric, bool) {
	query := "SELECT id, mtype, delta, value FROM metrics WHERE id=$1"

	var metric models.Metric

	row := p.db.QueryRow(query, ID)

	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		p.logger.Errorw("error getting metric", "id", ID, "error", err.Error())
		return metric, false
	}

	return metric, true
}

func (p PostgresMetricStorer) All() (map[string]models.Metric, error) {
	query := `SELECT id, mtype, delta, value FROM metrics;`

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Errorw("error getting all metrics", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	metrics := make(map[string]models.Metric)
	for rows.Next() {
		var metric models.Metric

		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value); err != nil {
			p.logger.Errorw("error scanning metric row", "error", err.Error())
			return nil, err
		}

		metrics[metric.ID] = metric
	}
	
	if err = rows.Err(); err != nil {
		p.logger.Errorw("error during rows iteration", "error", err.Error())
		return nil, err
	}

	return metrics, nil
}
