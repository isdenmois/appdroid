CREATE TABLE `apps` (
	`id` text PRIMARY KEY NOT NULL,
	`appId` text NOT NULL,
	`name` text NOT NULL,
	`version` text NOT NULL,
	`versionName` text NOT NULL,
	`type` text NOT NULL,
	`apk` text NOT NULL
);
