package backbone

import (
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
