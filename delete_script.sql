IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateway_logs]') AND type in (N'U'))
DROP TABLE [dbo].[gateway_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[uhf_status_logs]') AND type in (N'U'))
DROP TABLE [dbo].[uhf_status_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[operation_logs]') AND type in (N'U'))
DROP TABLE [dbo].[operation_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[system_logs]') AND type in (N'U'))
DROP TABLE [dbo].[system_logs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[user_accesses]') AND type in (N'U'))
DROP TABLE [dbo].[user_accesses]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[package_accesses]') AND type in (N'U'))
DROP TABLE [dbo].[package_accesses]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[uhfs]') AND type in (N'U'))
DROP TABLE [dbo].[uhfs]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[gateways]') AND type in (N'U'))
DROP TABLE [dbo].[gateways]

IF  EXISTS (SELECT *
FROM sys.objects
WHERE object_id = OBJECT_ID(N'[dbo].[areas]') AND type in (N'U'))
DROP TABLE [dbo].[areas]
