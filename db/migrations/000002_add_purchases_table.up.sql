CREATE TYPE purchase_type AS ENUM ('FOOD', 'DRINK', 'SNACK');

CREATE TABLE IF NOT EXISTS purchases(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    reservation_id uuid,
    type purchase_type NOT NULL,
    name varchar NOT NULL,
    count int NOT NULL,
    price_per_item_cents int NOT NULL,
    CONSTRAINT "RESERVATION_ID_FKEY" FOREIGN KEY (reservation_id) REFERENCES reservations(id)
);