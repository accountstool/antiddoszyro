package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *Store) UpsertSSLCertificate(ctx context.Context, domainID uuid.UUID, issuer string, status string, certPath string, keyPath string, expiresAt *time.Time, renewedAt *time.Time) error {
	_, err := s.db.Exec(ctx, `
		insert into ssl_certificates (domain_id, issuer, status, cert_path, key_path, expires_at, last_renew_at)
		values ($1,$2,$3,$4,$5,$6,$7)
		on conflict (domain_id)
		do update set issuer=excluded.issuer,
		              status=excluded.status,
		              cert_path=excluded.cert_path,
		              key_path=excluded.key_path,
		              expires_at=excluded.expires_at,
		              last_renew_at=excluded.last_renew_at,
		              updated_at=now()
	`, domainID, issuer, status, certPath, keyPath, nullableTime(expiresAt), nullableTime(renewedAt))
	return err
}
