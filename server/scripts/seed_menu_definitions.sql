-- Optional local seed if you want to populate menus without NATS.
INSERT INTO menu_definitions (id, domain, label, path, icon, order_index, parent_id, permissions, visible)
VALUES
    ('core:dashboard', 'auth', 'Dashboard', '/dashboard', 'dashboard', 0, NULL, ARRAY[]::text[], true),
    ('auth:system', 'auth', 'System', NULL, 'settings', 90, NULL, ARRAY['auth.role.read'], true),
    ('auth:roles', 'auth', 'Roles', '/system/roles', NULL, 10, 'auth:system', ARRAY['auth.role.read'], true),
    ('cms:root', 'cms', 'CMS', NULL, 'article', 20, NULL, ARRAY['cms.page.read'], true),
    ('cms:pages', 'cms', 'Pages', '/cms/pages', NULL, 10, 'cms:root', ARRAY['cms.page.read'], true),
    ('cms:media', 'cms', 'Media', '/cms/media', NULL, 20, 'cms:root', ARRAY['cms.page.read'], true)
ON CONFLICT (id) DO UPDATE SET
    domain = EXCLUDED.domain,
    label = EXCLUDED.label,
    path = EXCLUDED.path,
    icon = EXCLUDED.icon,
    order_index = EXCLUDED.order_index,
    parent_id = EXCLUDED.parent_id,
    permissions = EXCLUDED.permissions,
    visible = EXCLUDED.visible,
    updated_at = now();
