import { sqliteTable, text } from 'drizzle-orm/sqlite-core';

export const apps = sqliteTable('apps', {
  id: text().notNull().primaryKey(),
  appId: text().notNull(),
  name: text().notNull(),
  version: text().notNull(),
  versionName: text().notNull(),
  type: text().notNull(),
  apk: text().notNull(),
});

export type App = typeof apps.$inferSelect;
