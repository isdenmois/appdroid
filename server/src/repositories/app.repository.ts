import { eq } from 'drizzle-orm';
import { App, apps, db } from '../db';

export const appRepository = {
  getAll() {
    return db.select().from(apps).execute();
  },
  add(app: App) {
    return db
      .insert(apps)
      .values(app)
      .onConflictDoUpdate({
        target: apps.id,
        set: app,
      })
      .execute();
  },
  get(id: string) {
    return db.select().from(apps).where(eq(apps.id, id)).get();
  },
  delete(id: string) {
    return db.delete(apps).where(eq(apps.id, id)).execute();
  },
};
