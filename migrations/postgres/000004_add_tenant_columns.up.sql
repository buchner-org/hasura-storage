alter table storage.files add column tenant_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
alter table storage.files add column tenant_instance_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
alter table storage.files add column created_user UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
alter table storage.files add column updated_user UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';
alter table storage.files add column object_id UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000';

alter table storage.files add constraint fk_files_tenant_id FOREIGN KEY (tenant_id) REFERENCES public.tenant(id) ON UPDATE RESTRICT ON DELETE CASCADE;
alter table storage.files add constraint fk_files_tenant_instance_id FOREIGN KEY (tenant_id, tenant_instance_id) REFERENCES public.tenant_instance(tenant_id, id) ON UPDATE RESTRICT ON DELETE RESTRICT;

INSERT INTO storage.buckets (id) VALUES ('contact') on conflict do nothing;
INSERT INTO storage.buckets (id) VALUES ('documentation') on conflict do nothing;