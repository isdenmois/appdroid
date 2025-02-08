import Elysia, { t } from 'elysia';
import fs from 'node:fs/promises';
import { appRepository } from '../repositories';
import { apkService, filesService } from '../services';

export const appsController = new Elysia()
  .get('/list', () => appRepository.getAll())
  .post(
    '/upload',
    async ({ body }) => {
      const temp = await filesService.saveTemp(body.file);
      const info = await apkService.getInfo(temp);

      try {
        if (info) {
          const apk = `${info.appId}.apk`;

          await filesService.saveFile(temp, apk);

          await appRepository.add({
            ...info,
            id: info.appId,
            apk,
            type: 'static',
          });
        }
      } finally {
        fs.unlink(temp);
      }
    },
    { body: t.Object({ file: t.File() }) },
  )
  .delete('/:id', async ({ params: { id } }) => {
    const app = await appRepository.get(id);

    if (app?.apk) {
      console.log(app.apk);
      await filesService.removeFile(app.apk);
      await appRepository.delete(id);
    }
  });
