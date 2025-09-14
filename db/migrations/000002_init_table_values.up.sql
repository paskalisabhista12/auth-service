INSERT INTO public.roles (name, description)
VALUES 
  ('SUPERADMIN', 'Super administrator with full system access'),
  ('ADMIN', 'Administrator with elevated privileges'),
  ('USER', 'Standard user with limited access');

INSERT INTO public.permissions (name, description)
VALUES
    ('ALL', 'Permission to bypass all endpoint');
