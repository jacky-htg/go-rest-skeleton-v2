CREATE OR REPLACE FUNCTION public.int64_id(_table text, _col text)
 RETURNS bigint
 LANGUAGE plpgsql
AS $function$
DECLARE
  result BIGINT;
  numrows BIGINT;
BEGIN
  result = random_bigint();
  loop
    execute format('select 1 from %I where %I = %L', _table, _col, result);
    get diagnostics numrows = row_count;
    if numrows = 0 then
      RETURN result; 
    END if;
    result = random_bigint();
  END loop;
END;
$function$
;
