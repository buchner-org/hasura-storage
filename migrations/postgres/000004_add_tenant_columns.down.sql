alter table storage.files drop constraint fk_files_tenant_id;
alter table storage.files drop constraint fk_files_tenant_instance_id;

alter table storage.files drop column tenant_id;
alter table storage.files drop column tenant_instance_id;
alter table storage.files drop column created_user;
alter table storage.files drop column updated_user;
alter table storage.files drop column object_id;