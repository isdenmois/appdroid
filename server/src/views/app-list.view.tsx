import { Html } from '@elysiajs/html';

import { App } from '../db';
import { Page } from './page';

export const appListPageView = (baseUrl: string, apps: App[]) => {
  return (
    <Page>
      <h1>App Droid</h1>

      <div class="upload">
        <div class="all-apks">
          <div class="list-header">All APKs</div>

          <ul class="apps-list">
            {apps.map((app) => (
              <li>
                <a href={createLink(baseUrl, app)}>
                  <div class="name">{app.name}</div>
                  <div class="version">{app.versionName}</div>
                </a>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </Page>
  );
};

const createLink = (baseUrl: string, app: App) => {
  const additionalSettings = JSON.stringify({
    versionExtractionRegEx: '-(.*)\\.apk$',
    matchGroupToUse: '1',
    versionDetection: true,
  });

  const params = JSON.stringify({
    id: app.id,
    url: `${baseUrl}/app/${app.id}`,
    author: 'Appdroid',
    name: app.name,
    additionalSettings,
  });

  return `obtainium://app/${encodeURIComponent(params)}`;
};
