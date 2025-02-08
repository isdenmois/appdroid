import { Html } from '@elysiajs/html';

import { App } from '../db';
import { Page } from './page';

export function appPageView(apps: App[]) {
  return (
    <Page>
      <ul>{apps.map((app) => appItem(app))}</ul>
    </Page>
  );
}

function appItem(app: App) {
  const name = `${app.id}-${app.versionName}.apk`;
  const href = `/apk/${app.apk}/${name}`;

  return (
    <li>
      <a href={href}>{name}</a>
    </li>
  );
}
