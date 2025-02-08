import { Elysia } from 'elysia';
import { staticPlugin } from '@elysiajs/static';
import { html } from '@elysiajs/html';
import {
  appsController,
  filesController,
  obtainController,
  pingController,
} from './controllers';

const port = +(process.env.PORT || 3000);

new Elysia({
  serve: {
    maxRequestBodySize: 1024 * 1024 * 256, // 256MB
  },
})
  .use(staticPlugin({ prefix: '/', assets: './static' }))
  .use(html())
  .use(filesController)
  .use(obtainController)
  .group('/api', (g) => g.use(appsController).use(pingController))
  .listen(port);
