-- 1. Функция
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Триггер
CREATE TRIGGER trg_comments_set_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
