import Elysia from 'elysia';
import { filesService } from '../services';

export const filesController = new Elysia()
  .get('/file/:file', ({ params: { file } }) => filesService.getFile(file))
  .get('/apk/:file/:hash', ({ params: { file } }) =>
    filesService.getFile(file),
  );
