BEGIN;

ALTER TABLE transactions ADD COLUMN email TEXT;

UPDATE transactions SET email = 'jj@example.com' WHERE email IS NULL;

ALTER TABLE transactions ALTER COLUMN email SET NOT NULL ;

ALTER TABLE transactions ADD CONSTRAINT empty_email CHECK (email <> '');

COMMIT;

