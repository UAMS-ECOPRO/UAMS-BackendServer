IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateway_logs]') AND type in (N'U'))
DELETE FROM [dbo].[gateway_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[uhf_status_logs]') AND type in (N'U'))
DELETE FROM [dbo].[uhf_status_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[operation_logs]') AND type in (N'U'))
DELETE FROM [dbo].[operation_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[system_logs]') AND type in (N'U'))
DELETE FROM [dbo].[system_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[user_accesses]') AND type in (N'U'))
DELETE FROM [dbo].[user_accesses]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[package_accesses]') AND type in (N'U'))
DELETE FROM [dbo].[package_accesses]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[uhfs]') AND type in (N'U'))
DELETE FROM [dbo].[uhfs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateways]') AND type in (N'U'))
DELETE FROM [dbo].[gateways]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[areas]') AND type in (N'U'))
DELETE FROM [dbo].[areas]
