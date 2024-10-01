CREATE OR REPLACE FUNCTION public.random_bigint()
 RETURNS bigint
 LANGUAGE plpgsql
AS $function$
DECLARE
  result BIGINT = 0;
BEGIN
  SELECT FLOOR(random() * (999999999999999 - 5 + 1)) + 5 INTO result;
  RETURN result;
END;
$function$
;
