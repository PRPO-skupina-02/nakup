CREATE TYPE reservation_type AS ENUM ('ONLINE', 'POS');

CREATE TABLE IF NOT EXISTS reservations(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    user_id uuid NOT NULL,
    time_slot_id uuid NOT NULL,
    type reservation_type NOT NULL,
    row int NOT NULL,
    col int NOT NULL
);