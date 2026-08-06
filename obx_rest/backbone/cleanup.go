package backbone

import (
	"os"
	"strconv"
	"time"
)

func CleanupExpiredSessions() {
	for {
		result, err := PgSQL.Exec(
			`DELETE FROM "dat_user_session" WHERE expires_at < NOW()`,
		)
		if err != nil {
			Log.Warn().Err(err).Msg("failed to delete expired sessions")
		} else if n, _ := result.RowsAffected(); n > 0 {
			Log.Info().Int64("count", n).Msg("deleted expired sessions")
		}
		time.Sleep(1 * time.Hour)
	}
}

func CleanupContainerStatus() {
	retentionDays := 7
	if v := os.Getenv("IVL_DOCKER"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			retentionDays = n
		}
	}
	for {
		result, err := PgSQL.Exec(
			`DELETE FROM "ict_docker_container_stat" WHERE recorded_at < NOW() - ($1::int * INTERVAL '1 day')`,
			retentionDays,
		)
		if err != nil {
			Log.Warn().Err(err).Int("retention_days", retentionDays).Msg("failed to cleanup container status")
		} else if n, _ := result.RowsAffected(); n > 0 {
			Log.Info().Int64("count", n).Int("retention_days", retentionDays).Msg("deleted expired container status")
		}
		time.Sleep(1 * time.Hour)
	}
}

func CleanupHostStatus() {
	retentionDays := 7
	if v := os.Getenv("IVL_DOCKER"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			retentionDays = n
		}
	}
	for {
		result, err := PgSQL.Exec(
			`DELETE FROM "ict_vm_stat" WHERE recorded_at < NOW() - ($1::int * INTERVAL '1 day')`,
			retentionDays,
		)
		if err != nil {
			Log.Warn().Err(err).Int("retention_days", retentionDays).Msg("failed to cleanup host status")
		} else if n, _ := result.RowsAffected(); n > 0 {
			Log.Info().Int64("count", n).Int("retention_days", retentionDays).Msg("deleted expired host status")
		}
		time.Sleep(1 * time.Hour)
	}
}
