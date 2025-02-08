import Elysia from 'elysia';
import { appRepository } from '../repositories';
import { appListPageView, appPageView } from '../views';

export const obtainController = new Elysia()
  .get('/app/:id', async ({ params: { id } }) => {
    const app = appRepository.get(id);

    if (app) {
      return appPageView([app]);
    }

    return null;
  })
  .get('/apps', async ({ headers }) => {
    const apps = await appRepository.getAll();
    const { host } = headers;
    const scheme = host?.startsWith('192') ? `http` : `https`;
    const baseUrl = `${scheme}://${host}`;

    return appListPageView(baseUrl, apps);
  });
